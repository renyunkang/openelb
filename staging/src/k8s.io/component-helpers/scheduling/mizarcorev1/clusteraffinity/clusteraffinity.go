/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package clusteraffinity

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/api/mizarapps/v1alpha1"
	mizarcore "k8s.io/api/mizarcore/v1alpha1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// NodeSelector is a runtime representation of v1.NodeSelector.
type ClusterSelector struct {
	lazy LazyErrorClusterSelector
}

// LazyErrorNodeSelector is a runtime representation of v1.NodeSelector that
// only reports parse errors when no terms match.
type LazyErrorClusterSelector struct {
	terms []clusterSelectorTerm
}

// NewNodeSelector returns a NodeSelector or aggregate parsing errors found.
func NewNodeSelector(ns *v1.ClusterSelector, opts ...field.PathOption) (*ClusterSelector, error) {
	lazy := NewLazyErrorClusterSelector(ns, opts...)
	var errs []error
	for _, term := range lazy.terms {
		if len(term.parseErrs) > 0 {
			errs = append(errs, term.parseErrs...)
		}
	}
	if len(errs) != 0 {
		return nil, errors.Flatten(errors.NewAggregate(errs))
	}
	return &ClusterSelector{lazy: *lazy}, nil
}

// NewLazyErrorNodeSelector creates a NodeSelector that only reports parse
// errors when no terms match.
func NewLazyErrorClusterSelector(ns *v1.ClusterSelector, opts ...field.PathOption) *LazyErrorClusterSelector {
	p := field.ToPath(opts...)
	parsedTerms := make([]clusterSelectorTerm, 0, len(ns.ClusterSelectorTerms))
	path := p.Child("clusterSelectorTerms")
	for i, term := range ns.ClusterSelectorTerms {
		// nil or empty term selects no objects
		if isEmptyClusterSelectorTerm(&term) {
			continue
		}
		p := path.Index(i)
		parsedTerms = append(parsedTerms, newClusterSelectorTerm(&term, p))
	}
	return &LazyErrorClusterSelector{
		terms: parsedTerms,
	}
}

// Match checks whether the node labels and fields match the selector terms, ORed;
// nil or empty term matches no objects.
func (ns *ClusterSelector) Match(node *mizarcore.Cluster) bool {
	// parse errors are reported in NewNodeSelector.
	match, _ := ns.lazy.Match(node)
	return match
}

// Match checks whether the node labels and fields match the selector terms, ORed;
// nil or empty term matches no objects.
// Parse errors are only returned if no terms matched.
func (ns *LazyErrorClusterSelector) Match(node *mizarcore.Cluster) (bool, error) {
	if node == nil {
		return false, nil
	}
	nodeLabels := labels.Set(node.Labels)
	nodeFields := extractClusterFields(node)

	var errs []error
	for _, term := range ns.terms {
		match, tErrs := term.match(nodeLabels, nodeFields)
		if len(tErrs) > 0 {
			errs = append(errs, tErrs...)
			continue
		}
		if match {
			return true, nil
		}
	}
	return false, errors.Flatten(errors.NewAggregate(errs))
}

// PreferredSchedulingTerms is a runtime representation of []v1.PreferredSchedulingTerms.
type PreferredSchedulingTerms struct {
	terms []preferredSchedulingTerm
}

// Score returns a score for a Node: the sum of the weights of the terms that
// match the Node.
func (t *PreferredSchedulingTerms) Score(node *mizarcore.Cluster) int64 {
	var score int64
	nodeLabels := labels.Set(node.Labels)
	nodeFields := extractClusterFields(node)
	for _, term := range t.terms {
		// parse errors are reported in NewPreferredSchedulingTerms.
		if ok, _ := term.match(nodeLabels, nodeFields); ok {
			score += int64(term.weight)
		}
	}
	return score
}

func isEmptyClusterSelectorTerm(term *v1.ClusterSelectorTerm) bool {
	return len(term.MatchExpressions) == 0 && len(term.MatchFields) == 0
}

func extractClusterFields(n *mizarcore.Cluster) fields.Set {
	f := make(fields.Set)
	if len(n.Name) > 0 {
		f[v1alpha1.ClusterLabel] = n.Name
	}
	return f
}

type clusterSelectorTerm struct {
	matchLabels labels.Selector
	matchFields fields.Selector
	parseErrs   []error
}

func newClusterSelectorTerm(term *v1.ClusterSelectorTerm, path *field.Path) clusterSelectorTerm {
	var parsedTerm clusterSelectorTerm
	var errs []error
	if len(term.MatchExpressions) != 0 {
		p := path.Child("matchExpressions")
		parsedTerm.matchLabels, errs = clusterSelectorRequirementsAsSelector(term.MatchExpressions, p)
		if errs != nil {
			parsedTerm.parseErrs = append(parsedTerm.parseErrs, errs...)
		}
	}
	if len(term.MatchFields) != 0 {
		p := path.Child("matchFields")
		parsedTerm.matchFields, errs = clusterSelectorRequirementsAsFieldSelector(term.MatchFields, p)
		if errs != nil {
			parsedTerm.parseErrs = append(parsedTerm.parseErrs, errs...)
		}
	}
	return parsedTerm
}

func (t *clusterSelectorTerm) match(nodeLabels labels.Set, nodeFields fields.Set) (bool, []error) {
	if t.parseErrs != nil {
		return false, t.parseErrs
	}
	if t.matchLabels != nil && !t.matchLabels.Matches(nodeLabels) {
		return false, nil
	}
	if t.matchFields != nil && len(nodeFields) > 0 && !t.matchFields.Matches(nodeFields) {
		return false, nil
	}
	return true, nil
}

// nodeSelectorRequirementsAsSelector converts the []NodeSelectorRequirement api type into a struct that implements
// labels.Selector.
func clusterSelectorRequirementsAsSelector(nsm []v1.ClusterSelectorRequirement, path *field.Path) (labels.Selector, []error) {
	if len(nsm) == 0 {
		return labels.Nothing(), nil
	}
	var errs []error
	selector := labels.NewSelector()
	for i, expr := range nsm {
		p := path.Index(i)
		var op selection.Operator
		switch expr.Operator {
		case v1.ClusterSelectorOpIn:
			op = selection.In
		case v1.ClusterSelectorOpNotIn:
			op = selection.NotIn
		case v1.ClusterSelectorOpExists:
			op = selection.Exists
		case v1.ClusterSelectorOpDoesNotExist:
			op = selection.DoesNotExist
		case v1.ClusterSelectorOpGt:
			op = selection.GreaterThan
		case v1.ClusterSelectorOpLt:
			op = selection.LessThan
		default:
			errs = append(errs, field.NotSupported(p.Child("operator"), expr.Operator, nil))
			continue
		}
		r, err := labels.NewRequirement(expr.Key, op, expr.Values, field.WithPath(p))
		if err != nil {
			errs = append(errs, err)
		} else {
			selector = selector.Add(*r)
		}
	}
	if len(errs) != 0 {
		return nil, errs
	}
	return selector, nil
}

var validFieldSelectorOperators = []string{
	string(v1.ClusterSelectorOpIn),
	string(v1.ClusterSelectorOpNotIn),
}

// nodeSelectorRequirementsAsFieldSelector converts the []NodeSelectorRequirement core type into a struct that implements
// fields.Selector.
func clusterSelectorRequirementsAsFieldSelector(nsr []v1.ClusterSelectorRequirement, path *field.Path) (fields.Selector, []error) {
	if len(nsr) == 0 {
		return fields.Nothing(), nil
	}
	var errs []error

	var selectors []fields.Selector
	for i, expr := range nsr {
		p := path.Index(i)
		switch expr.Operator {
		case v1.ClusterSelectorOpIn:
			if len(expr.Values) != 1 {
				errs = append(errs, field.Invalid(p.Child("values"), expr.Values, "must have one element"))
			} else {
				selectors = append(selectors, fields.OneTermEqualSelector(expr.Key, expr.Values[0]))
			}

		case v1.ClusterSelectorOpNotIn:
			if len(expr.Values) != 1 {
				errs = append(errs, field.Invalid(p.Child("values"), expr.Values, "must have one element"))
			} else {
				selectors = append(selectors, fields.OneTermNotEqualSelector(expr.Key, expr.Values[0]))
			}

		default:
			errs = append(errs, field.NotSupported(p.Child("operator"), expr.Operator, validFieldSelectorOperators))
		}
	}

	if len(errs) != 0 {
		return nil, errs
	}
	return fields.AndSelectors(selectors...), nil
}

type preferredSchedulingTerm struct {
	clusterSelectorTerm
	weight int
}

type RequiredClusterAffinity struct {
	labelSelector   labels.Selector
	clusterSelector *LazyErrorClusterSelector
}

// GetRequiredNodeAffinity returns the parsing result of pod's nodeSelector and nodeAffinity.
func GetRequiredClusterAffinity(pod *v1.Pod) RequiredClusterAffinity {
	var selector labels.Selector
	if len(pod.Spec.ClusterSelector) > 0 {
		selector = labels.SelectorFromSet(pod.Spec.ClusterSelector)
	}
	// Use LazyErrorNodeSelector for backwards compatibility of parsing errors.
	var affinity *LazyErrorClusterSelector
	if pod.Spec.Affinity != nil &&
		pod.Spec.Affinity.ClusterAffinity != nil &&
		pod.Spec.Affinity.ClusterAffinity.RequiredDuringSchedulingIgnoredDuringExecution != nil {
		affinity = NewLazyErrorClusterSelector(pod.Spec.Affinity.ClusterAffinity.RequiredDuringSchedulingIgnoredDuringExecution)
	}
	return RequiredClusterAffinity{labelSelector: selector, clusterSelector: affinity}
}

// Match checks whether the pod is schedulable onto nodes according to
// the requirements in both nodeSelector and nodeAffinity.
func (s RequiredClusterAffinity) Match(cluster *mizarcore.Cluster) (bool, error) {
	if s.labelSelector != nil {
		if !s.labelSelector.Matches(labels.Set(cluster.Labels)) {
			return false, nil
		}
	}
	if s.clusterSelector != nil {
		return s.clusterSelector.Match(cluster)
	}
	return true, nil
}

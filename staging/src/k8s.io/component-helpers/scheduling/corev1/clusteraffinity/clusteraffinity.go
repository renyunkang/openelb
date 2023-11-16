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
	v1alpha1 "k8s.io/api/mizarcore/v1alpha1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ClusterSelector is a runtime representation of v1.ClusterSelector.
type ClusterSelector struct {
	lazy LazyErrorClusterSelector
}

// LazyErrorClusterSelector is a runtime representation of v1.ClusterSelector that
// only reports parse errors when no terms match.
type LazyErrorClusterSelector struct {
	terms []clusterSelectorTerm
}

// NewClusterSelector returns a ClusterSelector or aggregate parsing errors found.
func NewClusterSelector(ns *v1.ClusterSelector, opts ...field.PathOption) (*ClusterSelector, error) {
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

// NewLazyErrorClusterSelector creates a ClusterSelector that only reports parse
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

// Match checks whether the cluster labels and fields match the selector terms, ORed;
// nil or empty term matches no objects.
func (ns *ClusterSelector) Match(cluster *v1alpha1.Cluster) bool {
	// parse errors are reported in NewClusterSelector.
	match, _ := ns.lazy.Match(cluster)
	return match
}

// Match checks whether the cluster labels and fields match the selector terms, ORed;
// nil or empty term matches no objects.
// Parse errors are only returned if no terms matched.
func (ns *LazyErrorClusterSelector) Match(cluster *v1alpha1.Cluster) (bool, error) {
	if cluster == nil {
		return false, nil
	}
	clusterLabels := labels.Set(cluster.Labels)
	clusterFields := extractClusterFields(cluster)

	var errs []error
	for _, term := range ns.terms {
		match, tErrs := term.match(clusterLabels, clusterFields)
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

// NewPreferredSchedulingTerms returns a PreferredSchedulingTerms or all the parsing errors found.
// If a v1.PreferredSchedulingClusterTerm has a 0 weight, its parsing is skipped.
func NewPreferredSchedulingTerms(terms []v1.PreferredSchedulingClusterTerm, opts ...field.PathOption) (*PreferredSchedulingTerms, error) {
	p := field.ToPath(opts...)
	var errs []error
	parsedTerms := make([]preferredSchedulingTerm, 0, len(terms))
	for i, term := range terms {
		path := p.Index(i)
		if term.Weight == 0 || isEmptyClusterSelectorTerm(&term.Preference) {
			continue
		}
		parsedTerm := preferredSchedulingTerm{
			clusterSelectorTerm: newClusterSelectorTerm(&term.Preference, path),
			weight:              int(term.Weight),
		}
		if len(parsedTerm.parseErrs) > 0 {
			errs = append(errs, parsedTerm.parseErrs...)
		} else {
			parsedTerms = append(parsedTerms, parsedTerm)
		}
	}
	if len(errs) != 0 {
		return nil, errors.Flatten(errors.NewAggregate(errs))
	}
	return &PreferredSchedulingTerms{terms: parsedTerms}, nil
}

// Score returns a score for a Cluster: the sum of the weights of the terms that
// match the Cluster.
func (t *PreferredSchedulingTerms) Score(cluster *v1alpha1.Cluster) int64 {
	var score int64
	clusterLabels := labels.Set(cluster.Labels)
	clusterFields := extractClusterFields(cluster)
	for _, term := range t.terms {
		// parse errors are reported in NewPreferredSchedulingTerms.
		if ok, _ := term.match(clusterLabels, clusterFields); ok {
			score += int64(term.weight)
		}
	}
	return score
}

func isEmptyClusterSelectorTerm(term *v1.ClusterSelectorTerm) bool {
	return len(term.MatchExpressions) == 0 && len(term.MatchFields) == 0
}

func extractClusterFields(n *v1alpha1.Cluster) fields.Set {
	f := make(fields.Set)
	if len(n.Name) > 0 {
		f["metadata.name"] = n.Name
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

func (t *clusterSelectorTerm) match(clusterLabels labels.Set, clusterFields fields.Set) (bool, []error) {
	if t.parseErrs != nil {
		return false, t.parseErrs
	}
	if t.matchLabels != nil && !t.matchLabels.Matches(clusterLabels) {
		return false, nil
	}
	if t.matchFields != nil && len(clusterFields) > 0 && !t.matchFields.Matches(clusterFields) {
		return false, nil
	}
	return true, nil
}

// clusterSelectorRequirementsAsSelector converts the []ClusterSelectorRequirement api type into a struct that implements
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

// clusterSelectorRequirementsAsFieldSelector converts the []ClusterSelectorRequirement core type into a struct that implements
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

// GetRequiredClusterAffinity returns the parsing result of pod's clusterSelector and clusterAffinity.
func GetRequiredClusterAffinity(pod *v1.Pod) RequiredClusterAffinity {
	var selector labels.Selector
	if len(pod.Spec.ClusterSelector) > 0 {
		selector = labels.SelectorFromSet(pod.Spec.ClusterSelector)
	}
	// Use LazyErrorClusterSelector for backwards compatibility of parsing errors.
	var affinity *LazyErrorClusterSelector
	if pod.Spec.Affinity != nil &&
		pod.Spec.Affinity.ClusterAffinity != nil &&
		pod.Spec.Affinity.ClusterAffinity.RequiredDuringSchedulingIgnoredDuringExecution != nil {
		affinity = NewLazyErrorClusterSelector(pod.Spec.Affinity.ClusterAffinity.RequiredDuringSchedulingIgnoredDuringExecution)
	}
	return RequiredClusterAffinity{labelSelector: selector, clusterSelector: affinity}
}

// Match checks whether the pod is schedulable onto clusters according to
// the requirements in both clusterSelector and clusterAffinity.
func (s RequiredClusterAffinity) Match(cluster *v1alpha1.Cluster) (bool, error) {
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

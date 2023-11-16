package authorizer

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/apiserver/pkg/authorization/authorizer"
	"k8s.io/apiserver/pkg/endpoints/topology"
	mizarrbaclisters "k8s.io/client-go/listers/mizarrbac/v1alpha1"
	rbaclisters "k8s.io/client-go/listers/rbac/v1"
	"k8s.io/klog/v2"

	// install topology
	_ "k8s.io/IRBAC/topology"
)

func NewIRBACAuthorizer(
	roleGetter rbaclisters.RoleLister,
	roleBindingLister rbaclisters.RoleBindingLister,
	clusterRoleGetter rbaclisters.ClusterRoleLister,
	clusterRoleBindingLister rbaclisters.ClusterRoleBindingLister,
	globalRoleGetter mizarrbaclisters.GlobalRoleLister,
	globalRoleBindingLister mizarrbaclisters.GlobalRoleBindingLister,
) *IRBACAuthorizer {
	return &IRBACAuthorizer{authorizationRuleResolver: NewRequestToRuleMapper(
		roleGetter,
		roleBindingLister,
		clusterRoleGetter,
		clusterRoleBindingLister,
		globalRoleGetter,
		globalRoleBindingLister)}
}

type IRBACAuthorizer struct {
	authorizationRuleResolver RequestToRuleMapper
}

// authorizingVisitor short-circuits once allowed, and collects any resolution errors encountered
type authorizingVisitor struct {
	requestAttributes   authorizer.Attributes
	requestTopologyInfo topology.TopologyInfo

	allowed bool
	reason  string
	errors  []error
}

type TopologyRule struct {
	RegionSelectors  []metav1.LabelSelector
	ClusterSelectors []metav1.LabelSelector
}

func (v *authorizingVisitor) visit(source fmt.Stringer, rule *rbacv1.PolicyRule, topologyRule *TopologyRule, err error) bool {
	if topologyRule != nil && !topologyAllows(v.requestTopologyInfo, topologyRule) {
		return true
	}

	if rule != nil && ruleAllows(v.requestAttributes, rule) {
		v.allowed = true
		v.reason = fmt.Sprintf("IRBAC: allowed by %s", source.String())
		return false
	}

	if err != nil {
		v.errors = append(v.errors, err)
	}
	return true
}

func topologyAllows(requestInfo topology.TopologyInfo, rule *TopologyRule) bool {
	// the requested resource is not "topology resource" and will be skipped directly.
	if !requestInfo.IsTopologyResource() {
		return true
	}

	// If rule is nil , it means that user not have any permission at any region or cluster and will be forbidden directly
	if rule == nil {
		return false
	}

	var regionMatch, clusterMatch bool

	for _, region := range rule.RegionSelectors {
		selector, err := metav1.LabelSelectorAsSelector(&region)
		if err != nil {
			klog.Errorf("failed to get region selector")
			return false
		}
		for _, regionLabel := range requestInfo.GetRegionLabels() {
			if selector.Matches(labels.Set(regionLabel)) {
				regionMatch = true
				break
			}
		}
	}

	for _, clusterSelector := range rule.ClusterSelectors {
		selector, err := metav1.LabelSelectorAsSelector(&clusterSelector)
		if err != nil {
			klog.Errorf("failed to get cluster selector")
			return false
		}

		for _, clusterLabel := range requestInfo.GetClusterLabels() {
			if selector.Matches(labels.Set(clusterLabel)) {
				clusterMatch = true
				break
			}
		}
	}

	return regionMatch || clusterMatch
}

func ruleAllows(requestAttributes authorizer.Attributes, rule *rbacv1.PolicyRule) bool {
	if requestAttributes.IsResourceRequest() {
		combinedResource := requestAttributes.GetResource()
		if len(requestAttributes.GetSubresource()) > 0 {
			combinedResource = requestAttributes.GetResource() + "/" + requestAttributes.GetSubresource()
		}

		return VerbMatches(rule, requestAttributes.GetVerb()) &&
			APIGroupMatches(rule, requestAttributes.GetAPIGroup()) &&
			ResourceMatches(rule, combinedResource, requestAttributes.GetSubresource()) &&
			ResourceNameMatches(rule, requestAttributes.GetName())
	}

	return VerbMatches(rule, requestAttributes.GetVerb()) &&
		NonResourceURLMatches(rule, requestAttributes.GetPath())
}

func (i *IRBACAuthorizer) Authorize(ctx context.Context, requestAttributes authorizer.Attributes) (authorizer.Decision, string, error) {
	requestTopologyInfo, _ := topology.TopologyInfoFrom(ctx)

	if requestTopologyInfo.IsTopologyResource() {
		klog.V(5).Infof("requested TopologyInfo: cluster name: %s, region name: %s, cluster labels: %v, region labels: %v",
			requestTopologyInfo.GetClusterName().List(), requestTopologyInfo.GetRegionName().List(),
			requestTopologyInfo.GetClusterLabels(), requestTopologyInfo.GetRegionLabels())
	}

	ruleCheckingVisitor := &authorizingVisitor{requestAttributes: requestAttributes, requestTopologyInfo: requestTopologyInfo}

	i.authorizationRuleResolver.VisitRulesFor(requestAttributes.GetUser(), requestAttributes.GetNamespace(), ruleCheckingVisitor.visit)
	if ruleCheckingVisitor.allowed {
		return authorizer.DecisionAllow, ruleCheckingVisitor.reason, nil
	}

	// Build a detailed log of the denial.
	// Make the whole block conditional so we don't do a lot of string-building we won't use.
	if klogV := klog.V(5); klogV.Enabled() {
		var operation string
		if requestAttributes.IsResourceRequest() {
			b := &bytes.Buffer{}
			b.WriteString(`"`)
			b.WriteString(requestAttributes.GetVerb())
			b.WriteString(`" resource "`)
			b.WriteString(requestAttributes.GetResource())
			if len(requestAttributes.GetAPIGroup()) > 0 {
				b.WriteString(`.`)
				b.WriteString(requestAttributes.GetAPIGroup())
			}
			if len(requestAttributes.GetSubresource()) > 0 {
				b.WriteString(`/`)
				b.WriteString(requestAttributes.GetSubresource())
			}
			b.WriteString(`"`)
			if len(requestAttributes.GetName()) > 0 {
				b.WriteString(` named "`)
				b.WriteString(requestAttributes.GetName())
				b.WriteString(`"`)
			}
			operation = b.String()
		} else {
			operation = fmt.Sprintf("%q nonResourceURL %q", requestAttributes.GetVerb(), requestAttributes.GetPath())
		}

		var scope string
		if ns := requestAttributes.GetNamespace(); len(ns) > 0 {
			scope = fmt.Sprintf("in namespace %q", ns)
		} else {
			scope = "global-wide"
		}

		forbiddenInfo := fmt.Sprintf("IRBAC: no rules authorize user %q with groups %q to %s %s.",
			requestAttributes.GetUser().GetName(), requestAttributes.GetUser().GetGroups(), operation, scope)

		if requestTopologyInfo.IsTopologyResource() {
			forbiddenInfo = strings.TrimSuffix(forbiddenInfo, ".")
			forbiddenInfo = fmt.Sprintf("%s, with the region %s, the cluster %s", forbiddenInfo, requestTopologyInfo.GetRegionName(), requestTopologyInfo.GetClusterName())
		}
		klogV.Infof(forbiddenInfo)
	}

	reason := ""
	if len(ruleCheckingVisitor.errors) > 0 {
		reason = fmt.Sprintf("RBAC: %v", utilerrors.NewAggregate(ruleCheckingVisitor.errors))
	}
	return authorizer.DecisionNoOpinion, reason, nil
}

func (i *IRBACAuthorizer) RulesFor(user user.Info, namespace string) ([]authorizer.ResourceRuleInfo, []authorizer.NonResourceRuleInfo, bool, error) {
	var (
		resourceRules    []authorizer.ResourceRuleInfo
		nonResourceRules []authorizer.NonResourceRuleInfo
	)

	policyRules, err := i.authorizationRuleResolver.RulesFor(user, namespace)
	if err != nil {
		return nil, nil, true, err
	}
	for _, policyRule := range policyRules {
		if len(policyRule.Resources) > 0 {
			r := authorizer.DefaultResourceRuleInfo{
				Verbs:         policyRule.Verbs,
				APIGroups:     policyRule.APIGroups,
				Resources:     policyRule.Resources,
				ResourceNames: policyRule.ResourceNames,
			}
			var resourceRule authorizer.ResourceRuleInfo = &r
			resourceRules = append(resourceRules, resourceRule)
		}
		if len(policyRule.NonResourceURLs) > 0 {
			r := authorizer.DefaultNonResourceRuleInfo{
				Verbs:           policyRule.Verbs,
				NonResourceURLs: policyRule.NonResourceURLs,
			}
			var nonResourceRule authorizer.NonResourceRuleInfo = &r
			nonResourceRules = append(nonResourceRules, nonResourceRule)
		}
	}
	return resourceRules, nonResourceRules, false, err
}

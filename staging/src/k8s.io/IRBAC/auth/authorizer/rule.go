package authorizer

import (
	"fmt"

	mizarrbacv1alpha1 "k8s.io/api/mizarrbac/v1alpha1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/labels"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apiserver/pkg/authentication/serviceaccount"
	"k8s.io/apiserver/pkg/authentication/user"
	mizarrbaclisters "k8s.io/client-go/listers/mizarrbac/v1alpha1"
	rbaclisters "k8s.io/client-go/listers/rbac/v1"
)

type RequestToRuleMapper interface {
	// RulesFor returns all known PolicyRules and any errors that happened while locating those rules.
	// Any rule returned is still valid, since rules are deny by default.  If you can pass with the rules
	// supplied, you do not have to fail the request.  If you cannot, you should indicate the error along
	// with your denial.
	RulesFor(subject user.Info, namespace string) ([]rbacv1.PolicyRule, error)

	// VisitRulesFor invokes visitor() with each rule that applies to a given user in a given namespace,
	// and each error encountered resolving those rules. Rule may be nil if err is non-nil.
	// If visitor() returns false, visiting is short-circuited.
	VisitRulesFor(user user.Info, namespace string, visitor func(source fmt.Stringer, rule *rbacv1.PolicyRule, topologyRule *TopologyRule, err error) bool)

	GetRoleReferenceRules(roleRef rbacv1.RoleRef, bindingNamespace string) ([]rbacv1.PolicyRule, error)

	GetRoleReferenceTopologyRule(roleRef rbacv1.RoleRef) (*TopologyRule, error)
}

func NewRequestToRuleMapper(roleGetter rbaclisters.RoleLister,
	roleBindingLister rbaclisters.RoleBindingLister,
	clusterRoleGetter rbaclisters.ClusterRoleLister,
	clusterRoleBindingLister rbaclisters.ClusterRoleBindingLister,
	globalRoleGetter mizarrbaclisters.GlobalRoleLister,
	globalRoleBindingLister mizarrbaclisters.GlobalRoleBindingLister) RequestToRuleMapper {
	return RuleResolver{
		RuleOwnerLister: &DefaultRuleOwnerLister{
			RoleGetter:               roleGetter,
			RoleBindingLister:        roleBindingLister,
			ClusterRoleGetter:        clusterRoleGetter,
			ClusterRoleBindingLister: clusterRoleBindingLister,
			GlobalRoleGetter:         globalRoleGetter,
			GlobalRoleBindingLister:  globalRoleBindingLister,
		},
	}
}

type RuleResolver struct {
	RuleOwnerLister
}

func (r RuleResolver) GetRoleReferenceRules(roleRef rbacv1.RoleRef, bindingNamespace string) ([]rbacv1.PolicyRule, error) {
	switch roleRef.Kind {
	case "Role":
		role, err := r.RuleOwnerLister.GetRole(bindingNamespace, roleRef.Name)
		if err != nil {
			return nil, err
		}
		return role.Rules, nil

	case "ClusterRole":
		role, err := r.RuleOwnerLister.GetClusterRole(roleRef.Name)
		if err != nil {
			return nil, err
		}
		return role.Rules, nil

	case "GlobalRole":
		globalRole, err := r.RuleOwnerLister.GetGlobalRole(roleRef.Name)
		if err != nil {
			return nil, err
		}
		return globalRole.Rules, nil

	default:
		return nil, fmt.Errorf("unsupported role reference kind: %q", roleRef.Kind)
	}
}

func (r RuleResolver) GetRoleReferenceTopologyRule(roleRef rbacv1.RoleRef) (*TopologyRule, error) {
	switch roleRef.Kind {
	case "GlobalRole":
		role, err := r.RuleOwnerLister.GetGlobalRole(roleRef.Name)
		if err != nil {
			return nil, err
		}

		return &TopologyRule{
			RegionSelectors: role.RegionSelectors, ClusterSelectors: role.ClusterSelectors,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported role reference kind: %q", roleRef.Kind)
	}
}

func (r RuleResolver) VisitRulesFor(user user.Info, namespace string, visitor func(source fmt.Stringer, rule *rbacv1.PolicyRule, topologyRule *TopologyRule, err error) bool) {
	if globalRoleBindings, err := r.RuleOwnerLister.ListGlobalRoleBindings(); err != nil {
		if !visitor(nil, nil, nil, err) {
			return
		}
	} else {
		sourceDescriber := &globalRoleBindingDescriber{}
		for _, globalRoleBinding := range globalRoleBindings {
			subjectIndex, applies := appliesTo(user, globalRoleBinding.Subjects, "")
			if !applies {
				continue
			}
			rules, err := r.GetRoleReferenceRules(globalRoleBinding.RoleRef, "")
			if err != nil {
				if !visitor(nil, nil, nil, err) {
					return
				}
				continue
			}
			sourceDescriber.binding = globalRoleBinding
			sourceDescriber.subject = &globalRoleBinding.Subjects[subjectIndex]
			topologyRule, err := r.GetRoleReferenceTopologyRule(globalRoleBinding.RoleRef)
			if err != nil {
				if !visitor(nil, nil, nil, err) {
					return
				}
				continue
			}
			for i := range rules {
				if !visitor(sourceDescriber, &rules[i], topologyRule, nil) {
					return
				}
			}
		}
	}

	if clusterRoleBindings, err := r.RuleOwnerLister.ListClusterRoleBinding(); err != nil {
		if !visitor(nil, nil, nil, err) {
			return
		}
	} else {
		sourceDescriber := &clusterRoleBindingDescriber{}
		for _, clusterRoleBinding := range clusterRoleBindings {
			subjectIndex, applies := appliesTo(user, clusterRoleBinding.Subjects, "")
			if !applies {
				continue
			}
			rules, err := r.GetRoleReferenceRules(clusterRoleBinding.RoleRef, "")
			if err != nil {
				if !visitor(nil, nil, nil, err) {
					return
				}
				continue
			}
			sourceDescriber.binding = clusterRoleBinding
			sourceDescriber.subject = &clusterRoleBinding.Subjects[subjectIndex]
			for i := range rules {
				if !visitor(sourceDescriber, &rules[i], nil, nil) {
					return
				}
			}
		}
	}

	if len(namespace) > 0 {
		if roleBindings, err := r.RuleOwnerLister.ListRoleBindings(namespace); err != nil {
			if !visitor(nil, nil, nil, err) {
				return
			}
		} else {
			sourceDescriber := &roleBindingDescriber{}
			for _, roleBinding := range roleBindings {
				subjectIndex, applies := appliesTo(user, roleBinding.Subjects, namespace)
				if !applies {
					continue
				}
				rules, err := r.GetRoleReferenceRules(roleBinding.RoleRef, namespace)
				if err != nil {
					if !visitor(nil, nil, nil, err) {
						return
					}
					continue
				}
				sourceDescriber.binding = roleBinding
				sourceDescriber.subject = &roleBinding.Subjects[subjectIndex]
				for i := range rules {
					if !visitor(sourceDescriber, &rules[i], nil, nil) {
						return
					}
				}
			}
		}
	}
}

// appliesTo returns whether any of the bindingSubjects applies to the specified subject,
// and if true, the index of the first subject that applies
func appliesTo(user user.Info, bindingSubjects []rbacv1.Subject, namespace string) (int, bool) {
	for i, bindingSubject := range bindingSubjects {
		if appliesToUser(user, bindingSubject, namespace) {
			return i, true
		}
	}
	return 0, false
}

func has(set []string, ele string) bool {
	for _, s := range set {
		if s == ele {
			return true
		}
	}
	return false
}

func appliesToUser(user user.Info, subject rbacv1.Subject, namespace string) bool {
	switch subject.Kind {
	case rbacv1.UserKind:
		return user.GetName() == subject.Name

	case rbacv1.GroupKind:
		return has(user.GetGroups(), subject.Name)

	case rbacv1.ServiceAccountKind:
		// default the namespace to namespace we're working in if its available.  This allows rolebindings that reference
		// SAs in th local namespace to avoid having to qualify them.
		saNamespace := namespace
		if len(subject.Namespace) > 0 {
			saNamespace = subject.Namespace
		}
		if len(saNamespace) == 0 {
			return false
		}
		// use a more efficient comparison for RBAC checking
		return serviceaccount.MatchesUsername(saNamespace, subject.Name, user.GetName())
	default:
		return false
	}
}

type globalRoleBindingDescriber struct {
	binding *mizarrbacv1alpha1.GlobalRoleBinding
	subject *rbacv1.Subject
}

func (d *globalRoleBindingDescriber) String() string {
	return fmt.Sprintf("GlobalRoleBinding %q of %s %q to %s",
		d.binding.Name,
		d.binding.RoleRef.Kind,
		d.binding.RoleRef.Name,
		describeSubject(d.subject, ""),
	)
}

type clusterRoleBindingDescriber struct {
	binding *rbacv1.ClusterRoleBinding
	subject *rbacv1.Subject
}

func (d *clusterRoleBindingDescriber) String() string {
	return fmt.Sprintf("ClusterRoleBinding %q of %s %q to %s",
		d.binding.Name,
		d.binding.RoleRef.Kind,
		d.binding.RoleRef.Name,
		describeSubject(d.subject, d.binding.Namespace),
	)
}

type roleBindingDescriber struct {
	binding *rbacv1.RoleBinding
	subject *rbacv1.Subject
}

func (d *roleBindingDescriber) String() string {
	return fmt.Sprintf("RoleBinding %q of %s %q to %s",
		d.binding.Name+"/"+d.binding.Namespace,
		d.binding.RoleRef.Kind,
		d.binding.RoleRef.Name,
		describeSubject(d.subject, d.binding.Namespace),
	)
}

func describeSubject(s *rbacv1.Subject, bindingNamespace string) string {
	switch s.Kind {
	case rbacv1.ServiceAccountKind:
		if len(s.Namespace) > 0 {
			return fmt.Sprintf("%s %q", s.Kind, s.Name+"/"+s.Namespace)
		}
		return fmt.Sprintf("%s %q", s.Kind, s.Name+"/"+bindingNamespace)
	default:
		return fmt.Sprintf("%s %q", s.Kind, s.Name)
	}
}

type RuleOwnerLister interface {
	GetGlobalRole(name string) (*mizarrbacv1alpha1.GlobalRole, error)
	ListGlobalRoleBindings() ([]*mizarrbacv1alpha1.GlobalRoleBinding, error)
	GetClusterRole(name string) (*rbacv1.ClusterRole, error)
	ListClusterRoleBinding() ([]*rbacv1.ClusterRoleBinding, error)
	GetRole(namespace, name string) (*rbacv1.Role, error)
	ListRoleBindings(namespace string) ([]*rbacv1.RoleBinding, error)
}

type DefaultRuleOwnerLister struct {
	GlobalRoleGetter         mizarrbaclisters.GlobalRoleLister
	GlobalRoleBindingLister  mizarrbaclisters.GlobalRoleBindingLister
	ClusterRoleGetter        rbaclisters.ClusterRoleLister
	ClusterRoleBindingLister rbaclisters.ClusterRoleBindingLister
	RoleGetter               rbaclisters.RoleLister
	RoleBindingLister        rbaclisters.RoleBindingLister
}

func (i *DefaultRuleOwnerLister) GetGlobalRole(name string) (*mizarrbacv1alpha1.GlobalRole, error) {
	return i.GlobalRoleGetter.Get(name)
}

func (i *DefaultRuleOwnerLister) ListGlobalRoleBindings() ([]*mizarrbacv1alpha1.GlobalRoleBinding, error) {
	return i.GlobalRoleBindingLister.List(labels.Everything())
}

func (i *DefaultRuleOwnerLister) GetClusterRole(name string) (*rbacv1.ClusterRole, error) {
	return i.ClusterRoleGetter.Get(name)
}

func (i *DefaultRuleOwnerLister) ListClusterRoleBinding() ([]*rbacv1.ClusterRoleBinding, error) {
	return i.ClusterRoleBindingLister.List(labels.Everything())
}

func (i *DefaultRuleOwnerLister) GetRole(namespace, name string) (*rbacv1.Role, error) {
	return i.RoleGetter.Roles(namespace).Get(name)
}

func (i *DefaultRuleOwnerLister) ListRoleBindings(namespace string) ([]*rbacv1.RoleBinding, error) {
	return i.RoleBindingLister.RoleBindings(namespace).List(labels.Everything())
}

type ruleAccumulator struct {
	rules  []rbacv1.PolicyRule
	errors []error
}

func (r *ruleAccumulator) visit(source fmt.Stringer, rule *rbacv1.PolicyRule, topologyRule *TopologyRule, err error) bool {
	if rule != nil {
		r.rules = append(r.rules, *rule)
	}
	if err != nil {
		r.errors = append(r.errors, err)
	}
	return true
}

func (r RuleResolver) RulesFor(subject user.Info, namespace string) ([]rbacv1.PolicyRule, error) {
	visitor := &ruleAccumulator{}
	r.VisitRulesFor(subject, namespace, visitor.visit)
	return visitor.rules, utilerrors.NewAggregate(visitor.errors)
}

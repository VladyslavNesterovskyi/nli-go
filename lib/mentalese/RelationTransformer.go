package mentalese

import (
	"nli-go/lib/common"
)

type RelationTransformer struct {
	matcher *RelationMatcher
	log     *common.SystemLog
}

// using transformations
func NewRelationTransformer(matcher *RelationMatcher, log *common.SystemLog) *RelationTransformer {
	return &RelationTransformer{
		matcher: matcher,
		log:     log,
	}
}

// return the original relations, but replace the ones that have matched with their replacements
func (transformer *RelationTransformer) Replace(rules []Rule, relationSet RelationSet) RelationSet {

	// replace the relations embedded in quants
	replacedSet := transformer.replaceRelations(rules, relationSet, Binding{})

	return replacedSet
}

func (transformer *RelationTransformer) replaceRelations(transformations []Rule, relationSet RelationSet, binding Binding) RelationSet {

	replacedSet := RelationSet{}
	for _, relation := range relationSet {

		// replace inside hierarchical relations
		deepRelation := NewRelation(relation.Predicate, relation.Arguments)

		for i, argument := range deepRelation.Arguments {
			if argument.IsRelationSet() {
				deepRelation.Arguments[i] = NewRelationSet(transformer.replaceRelations(transformations, argument.TermValueRelationSet, binding))
			}
		}

		// replace according to rules
		newRelations := RelationSet{ }

		found := false
		for _, rule := range transformations {
			aBinding, ok := transformer.matcher.MatchTwoRelations(rule.Goal, deepRelation, Binding{})
			if  ok {
				boundRelations := rule.Pattern.BindSingle(aBinding)
				newRelations = append(newRelations, boundRelations...)
				found = true
			}
		}

		if !found {
			newRelations = append(newRelations, deepRelation)
		}

		replacedSet = append(replacedSet, newRelations...)
	}

	return replacedSet
}

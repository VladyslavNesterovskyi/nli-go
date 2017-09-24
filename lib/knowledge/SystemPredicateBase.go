package knowledge

import (
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"strconv"
)

type SystemPredicateBase struct {
	rules []mentalese.Rule
	log   *common.SystemLog
}

func NewSystemPredicateBase(log *common.SystemLog) *SystemPredicateBase {
	return &SystemPredicateBase{log: log}
}

func (ruleBase *SystemPredicateBase) GetMatchingGroups(set mentalese.RelationSet, knowledgeBaseIndex int) RelationGroups {

	matchingGroups := RelationGroups{}
	predicates := []string{"numberOf", "exists"}

	for _, setRelation := range set {
		for _, predicate:= range predicates {
			if predicate == setRelation.Predicate {
				matchingGroups = append(matchingGroups, RelationGroup{mentalese.RelationSet{setRelation}, knowledgeBaseIndex, 0.0})
				break
			}
		}
	}

	return matchingGroups
}

func (ruleBase *SystemPredicateBase) GetKnownRelations(set mentalese.RelationSet) mentalese.RelationSet {

	knownRelations := mentalese.RelationSet{}
	predicates := []string{"numberOf", "exists"}

	for _, setRelation := range set {
		for _, predicate:= range predicates {
			if predicate == setRelation.Predicate {
				knownRelations = append(knownRelations, setRelation)
				break
			}
		}
	}

	return knownRelations
}

func (ruleBase *SystemPredicateBase) Knows(relation mentalese.Relation) bool {

	predicates := []string{"numberOf", "exists"}

	found := false
	for _, predicate := range predicates{
		if predicate == relation.Predicate {
			found = true
			break
		}
	}
	return found
}

func (base *SystemPredicateBase) Bind(goal mentalese.Relation, bindings []mentalese.Binding) ([]mentalese.Binding, bool) {

	base.log.StartDebug("SystemPredicateBase Bind", goal, bindings)

	newBindings := []mentalese.Binding{}
	found := true
	aggregate := mentalese.Term{}

	resultArgument := goal.Arguments[0]
	resultVariable := resultArgument.TermValue

	if goal.Predicate == "numberOf" {

		subjectVariable := goal.Arguments[1].TermValue

		differentValues := base.getDifferentValues(bindings, subjectVariable)
		aggregate = mentalese.Term{TermType: mentalese.Term_number, TermValue: strconv.Itoa(len(differentValues))}

	} else if goal.Predicate == "exists" {

		subjectVariable := goal.Arguments[1].TermValue

		differentValues := base.getDifferentValues(bindings, subjectVariable)
		val := "false"
		if len(differentValues) > 0 {
			val = "true"
		}
		aggregate = mentalese.Term{TermType: mentalese.Term_predicateAtom, TermValue: val}

	} else {
		found = false
	}

	if found {
		newBindings = []mentalese.Binding{}

		// numberOf(4, E1)
		if resultArgument.IsNumber() {

			if resultArgument.TermValue == aggregate.TermValue {
				newBindings = bindings
			}

			// numberOf(N, E1)
		} else {

			if len(bindings) > 0 {

				for _, binding := range bindings {
					newBinding := binding.Copy()
					newBinding[resultVariable] = aggregate
					newBindings = append(newBindings, newBinding)
				}
			} else {

				newBinding := mentalese.Binding{}
				newBinding[resultVariable] = aggregate
				newBindings = append(newBindings, newBinding)

			}
		}
	}

	base.log.EndDebug("SystemPredicateBase Bind", newBindings, found)

	return newBindings, found
}

func (base *SystemPredicateBase) getDifferentValues(bindings []mentalese.Binding, subjectVariable string) []mentalese.Term {

	differentValues := []mentalese.Term{}

	for _, binding := range bindings {
		value, found := binding[subjectVariable]
		if found {
			preExists := false
			for _, differentValue := range differentValues {
				if differentValue.Equals(value) {
					preExists = true
				}
			}
			if !preExists {
				differentValues = append(differentValues, value)
			}
		}
	}

	return differentValues
}

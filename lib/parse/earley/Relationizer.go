package earley

import (
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse"
	"strconv"
)

// The relationizer turns a parse tree into a relation set
// It also subsumes the range and quantifier relation sets inside its quantification relation
type Relationizer struct {
	senseBuilder parse.SenseBuilder
	log          *common.SystemLog
}

func NewRelationizer(log *common.SystemLog) *Relationizer {
	return &Relationizer{
		senseBuilder: parse.NewSenseBuilder(),
		log:          log,
	}
}

func (relationizer Relationizer) Relationize(rootNode ParseTreeNode, nameResolver *central.NameResolver) (mentalese.RelationSet, mentalese.Binding) {
	rootEntityVariable := relationizer.senseBuilder.GetNewVariable("Sentence")
	sense, nameBinding, constantBinding := relationizer.extractSenseFromNode(rootNode, nameResolver, []string{ rootEntityVariable } )
	sense = sense.BindSingle(constantBinding)
	return sense, nameBinding
}

// Returns the sense of a node and its children
// node contains a rule with NP -> Det NBar
// antecedentVariable the actual variable used for the antecedent (for example: E5)
func (relationizer Relationizer) extractSenseFromNode(node ParseTreeNode, nameResolver *central.NameResolver, antecedentVariables []string) (mentalese.RelationSet, mentalese.Binding, mentalese.Binding) {

	relationizer.log.StartDebug("extractSenseFromNode", antecedentVariables, node.rule, node.rule.Sense)

	nameBinding := mentalese.Binding{}
	constantBinding := mentalese.Binding{}
	relationSet := mentalese.RelationSet{}

	if len(node.nameInformations) > 0 {
		firstAntecedentVariable := antecedentVariables[0]
		resolvedNameInformations := nameResolver.Resolve(node.nameInformations)
		for _, nameInformation := range resolvedNameInformations {
			nameBinding[firstAntecedentVariable] = mentalese.NewId(nameInformation.SharedId, nameInformation.EntityType)
		}
	}

	variableMap := relationizer.senseBuilder.CreateVariableMap(antecedentVariables, node.rule.GetAntecedentVariables(), node.rule.GetAllConsequentVariables())

	// create relations for each of the children
	boundChildSets := []mentalese.RelationSet{}
	for i, childNode := range node.constituents {

		consequentVariables := node.rule.GetConsequentVariables(i)

		mappedConsequentVariables := []string{}
		for _, consequentVariable := range consequentVariables {
			mappedConsequentVariables = append(mappedConsequentVariables, variableMap[consequentVariable].TermValue)
		}

		childRelations, childNameBinding, childConstantBinding := relationizer.extractSenseFromNode(childNode, nameResolver, mappedConsequentVariables)
		nameBinding = nameBinding.Merge(childNameBinding)
		boundChildSets = append(boundChildSets, childRelations)
		constantBinding = constantBinding.Merge(childConstantBinding)

		if node.rule.GetConsequentPositionType(i) == parse.PosTypeRegExp {
			constantBinding[antecedentVariables[0]] = mentalese.NewString(childNode.form)
		}
	}

	variableMap = relationizer.senseBuilder.ExtendVariableMap(node.rule.Sense, variableMap)

	boundParentSet := relationizer.senseBuilder.CreateGrammarRuleRelations(node.rule.Sense, variableMap)
	relationSet = relationizer.combineParentsAndChildren(boundParentSet, boundChildSets, node.rule)

	relationizer.log.EndDebug("extractSenseFromNode", relationSet)

	return relationSet, nameBinding, constantBinding
}

// Adds all childSets to parentSet
// Special case: if parentSet contains relation set placeholders [], like `quantification(X, [], Y, [])`, then these placeholders
// will be filled with the child set of the preceding variable
func (relationizer Relationizer) combineParentsAndChildren(parentSet mentalese.RelationSet, childSets []mentalese.RelationSet, rule parse.GrammarRule) mentalese.RelationSet {

	relationizer.log.StartDebug("processChildRelations", parentSet, childSets, rule)

	referencedChildrenIndexes := []int{}
	compoundRelation := mentalese.Relation{}

	// process sem(1) sem(2)
	newSet1 := mentalese.RelationSet{}
	for i, parentRelation := range parentSet {
		compoundRelation, referencedChildrenIndexes = relationizer.includeChildSenses(parentRelation, i, childSets, rule, referencedChildrenIndexes)
		newSet1 = append(newSet1, compoundRelation)
	}

	// collect children with sem(parent)
	childrenWithParentReferenceIndexes := relationizer.collectChildrenWithParentReferences(childSets)

	combination := parentSet

	// add simple children
	for i, childSet := range childSets {
		if !common.IntArrayContains(referencedChildrenIndexes, i) && !common.IntArrayContains(childrenWithParentReferenceIndexes, i) {
			combination = append(combination, childSet...)
		}
	}

	// raise children with sem(parent), eg. quants
	for _, i := range childrenWithParentReferenceIndexes {
		combination = relationizer.bindParent(combination, childSets[i])
	}

	relationizer.log.EndDebug("processChildRelations", combination)

	return combination
}

// Replaces a childRelation's sem(parent) with parentRelations and returns the new childRelation
func (relationizer Relationizer) bindParent(parentRelations mentalese.RelationSet, childSet mentalese.RelationSet) mentalese.RelationSet {

	newParentRelations := parentRelations

	for r, childRelation := range childSet {
		for a, argument := range childRelation.Arguments {
			if argument.IsRelationSet() {
				for _, argumentRelation := range argument.TermValueRelationSet {
					if argumentRelation.Predicate == mentalese.PredicateSem && argumentRelation.Arguments[0].TermValue == mentalese.AtomParent {

						if len(parentRelations) == 0 {
							// the quant should be raised above its "parent"; but in between may be empty-semantics nodes
							// note! there may in theory also be non-empty semantic nodes in between; don't know how to deal with this
							newParentRelations = childSet
						} else {
							prevParent := newParentRelations
							// the the sem of P is replaced by this quant
							newParentRelations = childSet.Copy()
							// the argument 'scope' in the quant of C is replaced by the current sem of P
							newParentRelations[r].Arguments[a] = mentalese.NewRelationSet(prevParent)
						}
					}
				}
			}
		}
	}

	return newParentRelations
}

func (relationizer Relationizer) collectChildrenWithParentReferences(childSets []mentalese.RelationSet) []int {

	childIndexes := []int{}

	for s, childSet := range childSets {
		for _, childRelation := range childSet {
			for _, argument := range childRelation.Arguments {
				if argument.IsRelationSet() {
					for _, argumentRelation := range argument.TermValueRelationSet {
						if argumentRelation.Predicate == mentalese.PredicateSem && argumentRelation.Arguments[0].TermValue == mentalese.AtomParent {
							childIndexes = append(childIndexes, s)
						}
					}
				}
			}
		}
	}

	return childIndexes
}

// Example:
// relation = quantification(E1, [], D1, [])
// extractedSetIndexes = []
// childSets = [ [], [isa(E1, dog)], [], [isa(D1, every)] ]
// rule = np(E1) -> dp(D1) nbar(E1);
func (relationizer Relationizer) includeChildSenses(parentRelation mentalese.Relation, childIndex int, childSets []mentalese.RelationSet, rule parse.GrammarRule, childIndexes []int) (mentalese.Relation, []int) {

	relationizer.log.StartDebug("includeChildSenses", parentRelation, childSets, rule)

	ruleRelation := rule.Sense[childIndex]

	for i, formalArgument := range ruleRelation.Arguments {
		if formalArgument.IsRelationSet() {
			if len(formalArgument.TermValueRelationSet) > 0 {
				firstRelation := formalArgument.TermValueRelationSet[0]
				if firstRelation.Predicate == mentalese.PredicateSem {
					index, err := strconv.Atoi(firstRelation.Arguments[0].TermValue)
					if err == nil {
						index = index - 1
						childIndexes = append(childIndexes, index)
						subSet := childSets[index]
						relationSetArgument := mentalese.Term{TermType: mentalese.TermRelationSet, TermValueRelationSet: subSet}
						parentRelation.Arguments[i] = relationSetArgument
					}
				}
			}
		}
	}

	relationizer.log.EndDebug("includeChildSenses", parentRelation, childIndexes)

	return parentRelation, childIndexes
}
package mentalese

import (
	"nli-go/lib/common"
	"strings"
)

type Relation struct {
	Predicate string
	Arguments []Term
}

const PredicateQuant = "quant"
const PredicateSem = "sem"
const PredicateSequence = "sequence"
const PredicateNot = "not"
const PredicateCall = "call"
const PredicateName = "name"
const PredicateText = "text"

const PredicateAssert = "assert"
const PredicateRetract = "retract"

const PredicateNumber = "number"
const PredicateThe = "the"
const PredicateSome = "some"
const PredicateReference = "reference"
const PredicateAll = "all"
const AtomAll = "all"
const AtomParent = "parent"

const QuantQuantifierVariableIndex = 0
const QuantQuantifierIndex = 1
const QuantRangeVariableIndex = 2
const QuantRangeIndex = 3
const QuantScopeIndex = 4

const SeqFirstOperandIndex = 1
const SeqSecondOperandIndex = 2

const NotScopeIndex = 0

func NewRelation(predicate string, arguments []Term) Relation {
	return Relation{
		Predicate: predicate,
		Arguments: arguments,
	}
}

func (relation Relation) GetVariableNames() []string {

	var names []string

	for _, argument := range relation.Arguments {
		if argument.IsVariable() {
			names = append(names, argument.TermValue)
		}
	}

	return common.StringArrayDeduplicate(names)
}

func (relation Relation) Equals(otherRelation Relation) bool {

	equals := relation.Predicate == otherRelation.Predicate

	for i, argument := range relation.Arguments {
		equals = equals && argument.Equals(otherRelation.Arguments[i])
	}

	return equals
}

func (relation Relation) Copy() Relation {

	newRelation := Relation{}
	newRelation.Predicate = relation.Predicate
	newRelation.Arguments = []Term{}
	for _, argument := range relation.Arguments {
		newRelation.Arguments = append(newRelation.Arguments, argument.Copy())
	}
	return newRelation
}

// Returns a new relation, that has all variables bound to bindings
func (relation Relation) BindSingleRelationSingleBinding(binding Binding) Relation {

	boundRelation := Relation{}
	boundRelation.Predicate = relation.Predicate

	for _, argument := range relation.Arguments {
		arg := argument.Bind(binding)
		boundRelation.Arguments = append(boundRelation.Arguments, arg)
	}

	return boundRelation
}

// Returns multiple relations, that has all variables bound to bindings
func (relation Relation) BindSingleRelationMultipleBindings(bindings Bindings) []Relation {

	boundRelations := []Relation{}

	for _, binding := range bindings {
		boundRelations = append(boundRelations, relation.BindSingleRelationSingleBinding(binding))
	}

	return boundRelations
}

// check if relation uses a variable (perhaps in one of its nested arguments)
func (relation Relation) RelationUsesVariable(variable string) bool {

	var found = false

	for _, argument := range relation.Arguments {
		if argument.IsVariable() {
			found = found || argument.TermValue == variable
		} else if argument.IsRelationSet() {
			for _, rel := range argument.TermValueRelationSet {
				found = found || rel.RelationUsesVariable(variable)
			}
		}
	}

	return found
}

func (relation Relation) ConvertVariablesToConstants() Relation {

	newRelation := Relation{ Predicate: relation.Predicate }

	for _, argument := range relation.Arguments {

		newArgument := argument.Copy()

		if argument.IsVariable() {
			newArgument = NewPredicateAtom(strings.ToLower(argument.TermValue))
		} else if argument.IsRelationSet() {
			newArgument = NewRelationSet(argument.TermValueRelationSet.ConvertVariablesToConstants())
		}

		newRelation.Arguments = append(newRelation.Arguments, newArgument)
	}

	return newRelation
}

func (relation Relation) String() string {

	args, sep := "", ""

	for _, Argument := range relation.Arguments {

		args += sep + Argument.String()
		sep = ", "
	}

	return relation.Predicate + "(" + args + ")"
}

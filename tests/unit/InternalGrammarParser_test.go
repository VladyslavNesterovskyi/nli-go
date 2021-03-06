package tests

import (
	"fmt"
	"nli-go/lib/importer"
	"testing"
)

func TestInternalGrammarParser(t *testing.T) {

	parser := importer.NewInternalGrammarParser()

	tests := []string{
		"determiner(E, [], D, [])",
	}

	for _, test := range tests {
		result := parser.CreateRelation(test)
		if result.String() != test {
			t.Errorf("got %s, want %s", result.String(), test)
		}
	}

	// =====================================================

	grammar := parser.CreateGrammar("[" +
		"{ rule: s(P) -> np(E) vp(P),         sense: subject(P, E) }" +
		"]")

	rules := grammar.FindRules("s", 1)
	if len(rules) == 0 {
		t.Error("No rules found")
	}

	if rules[0].GetAntecedent() != "s" {
		t.Error(fmt.Printf("Error in rule: %s", rules[0].GetAntecedent()))
	}
	if rules[0].GetConsequent(0) != "np" {
		t.Error(fmt.Printf("Error in rule: %s", rules[0].GetConsequent(0)))
	}
	if rules[0].GetAntecedentVariables()[0] != "P" {
		t.Error(fmt.Printf("Error in rule: %s", rules[0].GetAntecedentVariables()))
	}
	if rules[0].GetConsequentVariables(0)[0] != "E" {
		t.Error(fmt.Printf("Error in rule: %s", rules[0].GetConsequentVariables(0)))
	}
	if len(rules[0].Sense) != 1 {
		t.Error(fmt.Printf("Error in number of sense relations: %d", len(rules[0].Sense)))
	}

	grammar = parser.CreateGrammar("[" +
		"{ rule: s(P) -> np(E) vp(P),    sense: subject(P, E) }" +
		"{ rule: np(P) -> nbar(E) }" +
		"]")

	rules = grammar.FindRules("s", 1)
	if len(rules) != 1 {
		t.Error("No rules found")
	}
	rules = grammar.FindRules("np", 1)
	if len(rules) != 1 {
		t.Error("No rules found")
	}

	grammar = parser.CreateGrammar("[]")

	parser.CreateRelationSet("[assert(at(5, 3))]")
}

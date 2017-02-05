package tests

import (
	"testing"
	"nli-go/lib/parse"
	"nli-go/lib/importer"
	"nli-go/lib/mentalese"
	"nli-go/lib/central"
	"nli-go/lib/knowledge"
	"nli-go/lib/common"
	"nli-go/lib/parse/earley"
)

func TestRelease1(t *testing.T) {

	internalGrammarParser := importer.NewInternalGrammarParser()
	internalGrammarParser.SetPanicOnParseFail(true)

	// Data

	grammar := internalGrammarParser.LoadGrammar(common.GetCurrentDir() + "/../../resources/english-1.grammar")

	lexicon := internalGrammarParser.CreateLexicon(`[
		form: 'who',        pos: whWord;
		form: 'married',    pos: verb, 	        sense: isa(E, marry);
		form: 'did',		pos: auxDo;
		form: 'marry',		pos: verb,		    sense: isa(E, marry);
		form: 'de',		    pos: insertion,     sense: name(E, 'de', insertion);
		form: 'van',		pos: insertion,     sense: name(E, 'van', insertion);
		form: /^[A-Z]/,	    pos: lastName,      sense: name(E, Form, lastName);
		form: /^[A-Z]/,	    pos: firstName,     sense: name(E, Form, firstName);
		form: 'are',		pos: auxBe,		    sense: isa(E, be);
		form: 'and',		pos: conjunction;
		form: 'siblings',	pos: noun,		    sense: isa(E, sibling);
		form: '?',          pos: questionMark;
	]`)

	generic2ds := internalGrammarParser.CreateTransformations(`[
		married_to(A, B) :- isa(P1, marry) subject(P1, A) object(P1, B);
		siblings(A1, A2) :- isa(P1, be) subject(P1, A) conjunction(A, A1, A2) object(P1, B) isa(B, sibling);
		name(A, N) :- name(A, N, _);
		act(question, who) married_to(A, B) :- question(Q) isa(Q, marry) subject(Q, A) object(Q, B);
		act(question, yesno) :- question(Q) isa(Q, marry) subject(Q, A) object(Q, B);
		act(question, howmany) child(A, B) :- question(Q) isa(Q, marry) subject(Q, A) object(Q, B);
	]`)

	dsSolutions := internalGrammarParser.CreateSolutions(`[
		condition: married_to(A, B) question(A),
		answer: grammatical_subject(B) married_to(A, B) gender(B, G) name(A, N);
	]`)

	dsInferenceRules := internalGrammarParser.CreateRules(`[
		sibling(A, B) :- parent(C, A) parent(C, B);
	]`)

	ds2db := internalGrammarParser.CreateRules(`[
		married_to(A, B) :- marriages(A, B, _);
		name(A, N) :- person(A, N, _, _);
		gender(A, male) :- person(A, _, 'M', _);
		gender(A, female) :- person(A, _, 'F', _);
	]`)

	dbFacts := internalGrammarParser.CreateRelationSet(`[
		marriages(1, 2, '1992')
		parent(4, 2)
		parent(4, 3)
		person(1, 'Jacqueline de Boer', 'F', '1964')
		person(2, 'Mark van Dongen', 'M', '1967')
		person(3, 'Suzanne van Dongen', 'F', '1967')
		person(4, 'John van Dongen', 'M', '1938')
	]`)

	// Services

	tokenizer := parse.NewTokenizer()
	parser := earley.NewParser(grammar, lexicon)
	transformer := mentalese.NewRelationTransformer()
	factBase1 := knowledge.NewFactBase(dbFacts, ds2db)
	ruleBase1 := knowledge.NewRuleBase(dsInferenceRules)
	problemSolver := central.NewAnswerer()
	problemSolver.AddSolutions(dsSolutions)
	problemSolver.AddKnowledgeBase(factBase1)
	problemSolver.AddKnowledgeBase(ruleBase1)

	// Tests

	var tests = []struct {
		question string
		want string
	} {
		{"Who married Jacqueline de Boer?", "Marty"},
		//{"Did Bob marry Sally?", "Yes"},
		//{"Are Jane and Janelle siblings?", "No"},
		//{"Which children has John van Dongen?", "Mark van Dongen and Suzanne van Dongen"},
		//{"How many children has John van Dongen?", "He has 2 children"},
	}

	for _, test := range tests {

		tokens := tokenizer.Process(test.question)

common.LoggerActive=false
		genericSense, _, _ := parser.Parse(tokens)
common.LoggerActive=false
//fmt.Print(genericSense)
		domainSpecificSense := transformer.Extract(generic2ds, genericSense)

		//dsAnswersSenses :=
			problemSolver.Answer(domainSpecificSense)

		answer := ""

		if answer != test.want {
//			t.Errorf("release1: got %v, want %v", answer, test.want)
		}
	}
}

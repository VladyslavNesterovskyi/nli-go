package tests

import (
	"nli-go/lib/central"
	"nli-go/lib/common"
	"nli-go/lib/importer"
	"nli-go/lib/mentalese"
	"nli-go/lib/parse/earley"
	"testing"
)

func TestFillerStack(t *testing.T) {

	internalGrammarParser := importer.NewInternalGrammarParser()

	grammar := internalGrammarParser.CreateGrammar(`[

		{ rule: s(P1) -> which() np(E1) dep_vp(P1, E1),			sense: which(E1) }
		{ rule: np(E1) -> nbar(E1), 							sense: quant(_, some(_), E1, sem(1), sem(parent)) }
		{ rule: nbar(E) -> noun(E) }
		{ rule: dep_vp(P1, E1) -> be() np(E2) advp(P1) vp(P1, E1, E2) }
		{ rule: np(E1) -> qp(Q1) nbar(E1), 						sense: quant(Q1, sem(1), E1, sem(2), sem(parent)) }
		{ rule: advp(P1) -> adverb(P1) }
		{ rule: vp(P1, E1, E2) -> to() take() from(), 			sense: take_from(P1, E2, E1)  }
	]`)

	lexicon := internalGrammarParser.CreateLexicon(`[
		{ form: 'which', pos: which }
		{ form: 'were', pos: be }
		{ form: 'babies', pos: noun, 		sense: baby(E) }
		{ form: 'toys', pos: noun, 			sense: toy(E) }
		{ form: 'the', pos: qp, 				sense: the(E) }
		{ form: 'easiest', pos: adverb, 		sense: easiest(E) }
		{ form: 'to', pos: to }
		{ form: 'take', pos: take }
		{ form: 'from', pos: from }
	]`)
	log := common.NewSystemLog(true)

	matcher := mentalese.NewRelationMatcher(log)
	dialogContext := central.NewDialogContext()
	predicates := mentalese.Predicates{}
	solver := central.NewProblemSolver(matcher, predicates, dialogContext, log)
	nameResolver := central.NewNameResolver(solver, matcher, predicates, log, dialogContext)

	parser := earley.NewParser(grammar, lexicon, nameResolver, predicates, log)

	relationizer := earley.NewRelationizer(lexicon, log)

	parseTrees := parser.Parse([]string{"Which", "babies", "were", "the", "toys", "easiest", "to", "take", "from"})

	if len(parseTrees) == 0 {
		t.Error(log.String())
		return
	}

	result, _ := relationizer.Relationize(parseTrees[0], nameResolver)

	want := "[quant(_, [some(_)], E5, [baby(E5)], [which(E5) quant(Q5, [the(Q5)], E6, [toy(E6)], [easiest(S5) take_from(S5, E6, E5)])])]"
	if result.String() != want {
		t.Errorf("got %s, want %s", result.String(), want)
	}
}

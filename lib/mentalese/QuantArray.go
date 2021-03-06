package mentalese

type QuantArray RelationSet

func (s QuantArray) Len() int {
	return len(s)
}

func (s QuantArray) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s QuantArray) Less(i, j int) bool {

	less := false

	first := s[i]
	firstQuantifier := first.Arguments[QuantQuantifierIndex]

	second := s[j]
	secondQuantifier := second.Arguments[QuantQuantifierIndex]

	firstQuantifierSimple := ""
	secondQuantifierSimple := ""

	// for now, we're just doing `all`, and `some`
	if len(firstQuantifier.TermValueRelationSet) == 1 {
		// isa(X, all)
		firstQuantifierSimple = firstQuantifier.TermValueRelationSet[0].Arguments[1].TermValue
	}
	if len(secondQuantifier.TermValueRelationSet) == 1 {
		secondQuantifierSimple = secondQuantifier.TermValueRelationSet[0].Arguments[1].TermValue
	}

	if firstQuantifierSimple == AtomAll && secondQuantifierSimple == AtomAll {
		less = i < j
	} else if firstQuantifierSimple == AtomAll {
		less = true
	} else if secondQuantifierSimple == AtomAll {
		less = false
	} else if false { // interrogative determiner
		less = true
	} else {
		// by default, reading order is order of preference
		less = i < j
	}

	return less
}

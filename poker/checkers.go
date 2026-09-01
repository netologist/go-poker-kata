package poker

type royalFlushChecker struct{}

func (royalFlushChecker) Check(cards []Card) (HandResult, bool) {
	if isFlush(cards) && straightHighValueOrZero(cards) == int(Ace) {
		return HandResult{Rank: RoyalFlush, Tiebreakers: []int{int(Ace)}}, true
	}
	return HandResult{}, false
}

type straightFlushChecker struct{}

func (straightFlushChecker) Check(cards []Card) (HandResult, bool) {
	high := straightHighValueOrZero(cards)
	if isFlush(cards) && high != 0 {
		return HandResult{Rank: StraightFlush, Tiebreakers: []int{high}}, true
	}
	return HandResult{}, false
}

type fourOfAKindChecker struct{}

func (fourOfAKindChecker) Check(cards []Card) (HandResult, bool) {
	counts := valueCountsByFrequency(cards)
	if counts[0].count == 4 {
		return HandResult{Rank: FourOfAKind, Tiebreakers: valuesOf(counts)}, true
	}
	return HandResult{}, false
}

type fullHouseChecker struct{}

func (fullHouseChecker) Check(cards []Card) (HandResult, bool) {
	counts := valueCountsByFrequency(cards)
	if len(counts) == 2 && counts[0].count == 3 {
		return HandResult{Rank: FullHouse, Tiebreakers: valuesOf(counts)}, true
	}
	return HandResult{}, false
}

type flushChecker struct{}

func (flushChecker) Check(cards []Card) (HandResult, bool) {
	if isFlush(cards) {
		return HandResult{Rank: Flush, Tiebreakers: sortedValuesDesc(cards)}, true
	}
	return HandResult{}, false
}

type straightChecker struct{}

func (straightChecker) Check(cards []Card) (HandResult, bool) {
	high := straightHighValueOrZero(cards)
	if high != 0 {
		return HandResult{Rank: Straight, Tiebreakers: []int{high}}, true
	}
	return HandResult{}, false
}

type threeOfAKindChecker struct{}

func (threeOfAKindChecker) Check(cards []Card) (HandResult, bool) {
	counts := valueCountsByFrequency(cards)
	if len(counts) == 3 && counts[0].count == 3 {
		return HandResult{Rank: ThreeOfAKind, Tiebreakers: valuesOf(counts)}, true
	}
	return HandResult{}, false
}

type twoPairChecker struct{}

func (twoPairChecker) Check(cards []Card) (HandResult, bool) {
	counts := valueCountsByFrequency(cards)
	if len(counts) == 3 && counts[0].count == 2 && counts[1].count == 2 {
		return HandResult{Rank: TwoPair, Tiebreakers: valuesOf(counts)}, true
	}
	return HandResult{}, false
}

type onePairChecker struct{}

func (onePairChecker) Check(cards []Card) (HandResult, bool) {
	counts := valueCountsByFrequency(cards)
	if len(counts) == 4 && counts[0].count == 2 {
		return HandResult{Rank: OnePair, Tiebreakers: valuesOf(counts)}, true
	}
	return HandResult{}, false
}

// highCardChecker is the fallback: any five cards form at least a high card hand.
type highCardChecker struct{}

func (highCardChecker) Check(cards []Card) (HandResult, bool) {
	return HandResult{Rank: HighCard, Tiebreakers: sortedValuesDesc(cards)}, true
}

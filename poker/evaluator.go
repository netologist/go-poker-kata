package poker

import "fmt"

// HandEvaluator evaluates an exact 5-card hand into its HandResult by trying
// a list of HandChecker implementations from strongest to weakest category;
// the first match wins. It depends only on the HandChecker interface
// (dependency inversion) and the checker list can be extended without
// modifying this type (open/closed).
type HandEvaluator struct {
	checkersStrongestFirst []HandChecker
}

// NewHandEvaluator builds an evaluator with the standard set of checkers.
func NewHandEvaluator() *HandEvaluator {
	return NewHandEvaluatorWithCheckers([]HandChecker{
		royalFlushChecker{},
		straightFlushChecker{},
		fourOfAKindChecker{},
		fullHouseChecker{},
		flushChecker{},
		straightChecker{},
		threeOfAKindChecker{},
		twoPairChecker{},
		onePairChecker{},
		highCardChecker{},
	})
}

// NewHandEvaluatorWithCheckers allows injecting a custom ordered checker list
// (e.g. for testing, or a house-rule variant with extra hand categories).
func NewHandEvaluatorWithCheckers(checkers []HandChecker) *HandEvaluator {
	copied := make([]HandChecker, len(checkers))
	copy(copied, checkers)
	return &HandEvaluator{checkersStrongestFirst: copied}
}

// Evaluate scores exactly five cards. It returns an error if fewer or more
// than five cards are given.
func (e *HandEvaluator) Evaluate(fiveCards []Card) (HandResult, error) {
	if len(fiveCards) != 5 {
		return HandResult{}, fmt.Errorf("a poker hand must contain exactly 5 cards, got %d", len(fiveCards))
	}
	for _, checker := range e.checkersStrongestFirst {
		if result, ok := checker.Check(fiveCards); ok {
			return result, nil
		}
	}
	// Unreachable: highCardChecker always matches.
	return HandResult{}, fmt.Errorf("no checker matched hand: %v", fiveCards)
}

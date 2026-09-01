package poker

import "fmt"

// BestHandFinder finds the highest-ranked 5-card poker hand obtainable from
// a set of 5 or more cards (e.g. picking the best 5 out of 7 hole+community
// cards). It composes a HandEvaluator rather than embedding its logic,
// keeping combination-generation and hand-ranking as separate responsibilities.
type BestHandFinder struct {
	evaluator *HandEvaluator
}

func NewBestHandFinder() *BestHandFinder {
	return &BestHandFinder{evaluator: NewHandEvaluator()}
}

func NewBestHandFinderWithEvaluator(evaluator *HandEvaluator) *BestHandFinder {
	return &BestHandFinder{evaluator: evaluator}
}

func (f *BestHandFinder) FindBestHand(cards []Card) (HandResult, error) {
	if len(cards) < 5 {
		return HandResult{}, fmt.Errorf("at least 5 cards are required, got %d", len(cards))
	}
	if len(cards) == 5 {
		return f.evaluator.Evaluate(cards)
	}

	var best HandResult
	haveBest := false
	for _, combo := range fiveCardCombinations(cards) {
		candidate, err := f.evaluator.Evaluate(combo)
		if err != nil {
			return HandResult{}, err
		}
		if !haveBest || candidate.Compare(best) > 0 {
			best = candidate
			haveBest = true
		}
	}
	return best, nil
}

func fiveCardCombinations(cards []Card) [][]Card {
	var results [][]Card
	current := make([]Card, 0, 5)
	var combine func(start int)
	combine = func(start int) {
		if len(current) == 5 {
			combo := make([]Card, 5)
			copy(combo, current)
			results = append(results, combo)
			return
		}
		for i := start; i < len(cards); i++ {
			current = append(current, cards[i])
			combine(i + 1)
			current = current[:len(current)-1]
		}
	}
	combine(0)
	return results
}

package poker

import "fmt"

// HandResult is the outcome of evaluating a 5-card hand: its category plus
// an ordered list of tiebreaker values (most significant first), used to
// compare two hands that share the same category.
type HandResult struct {
	Rank        HandRank
	Tiebreakers []int
}

// Compare returns a negative number if h is weaker than other, zero if equal
// in strength, and a positive number if h is stronger.
func (h HandResult) Compare(other HandResult) int {
	if h.Rank != other.Rank {
		return int(h.Rank) - int(other.Rank)
	}
	length := len(h.Tiebreakers)
	if len(other.Tiebreakers) < length {
		length = len(other.Tiebreakers)
	}
	for i := 0; i < length; i++ {
		if diff := h.Tiebreakers[i] - other.Tiebreakers[i]; diff != 0 {
			return diff
		}
	}
	return len(h.Tiebreakers) - len(other.Tiebreakers)
}

func (h HandResult) String() string {
	return fmt.Sprintf("%s %v", h.Rank, h.Tiebreakers)
}

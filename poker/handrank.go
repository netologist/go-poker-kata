package poker

// HandRank is one of the ten standard poker hand categories, ordered by strength.
type HandRank int

const (
	HighCard HandRank = iota + 1
	OnePair
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
	RoyalFlush
)

func (h HandRank) String() string {
	names := map[HandRank]string{
		HighCard:      "High Card",
		OnePair:       "One Pair",
		TwoPair:       "Two Pair",
		ThreeOfAKind:  "Three of a Kind",
		Straight:      "Straight",
		Flush:         "Flush",
		FullHouse:     "Full House",
		FourOfAKind:   "Four of a Kind",
		StraightFlush: "Straight Flush",
		RoyalFlush:    "Royal Flush",
	}
	if name, ok := names[h]; ok {
		return name
	}
	return "Unknown"
}

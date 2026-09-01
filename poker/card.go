package poker

import "fmt"

// Rank is a card's face value. The underlying int doubles as its numeric
// strength (2..14, Ace high).
type Rank int

const (
	Two   Rank = 2
	Three Rank = 3
	Four  Rank = 4
	Five  Rank = 5
	Six   Rank = 6
	Seven Rank = 7
	Eight Rank = 8
	Nine  Rank = 9
	Ten   Rank = 10
	Jack  Rank = 11
	Queen Rank = 12
	King  Rank = 13
	Ace   Rank = 14
)

func (r Rank) String() string {
	names := map[Rank]string{
		Two: "2", Three: "3", Four: "4", Five: "5", Six: "6", Seven: "7",
		Eight: "8", Nine: "9", Ten: "10", Jack: "J", Queen: "Q", King: "K", Ace: "A",
	}
	if name, ok := names[r]; ok {
		return name
	}
	return fmt.Sprintf("Rank(%d)", int(r))
}

// Suit is a card's suit.
type Suit int

const (
	Hearts Suit = iota
	Diamonds
	Clubs
	Spades
)

func (s Suit) String() string {
	switch s {
	case Hearts:
		return "H"
	case Diamonds:
		return "D"
	case Clubs:
		return "C"
	case Spades:
		return "S"
	default:
		return "?"
	}
}

// Card is a single playing card.
type Card struct {
	Rank Rank
	Suit Suit
}

func NewCard(rank Rank, suit Suit) Card {
	return Card{Rank: rank, Suit: suit}
}

func (c Card) String() string {
	return c.Rank.String() + c.Suit.String()
}

// Command poker parses cards given as command line arguments in the form
// "RANKSUIT" (e.g. "AH KH QH JH TH") and prints the best 5-card hand found.
// Rank letters: 2-9,T,J,Q,K,A. Suit letters: H,D,C,S.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/netologist/go-poker-kata/poker"
)

var rankLetters = map[byte]poker.Rank{
	'2': poker.Two, '3': poker.Three, '4': poker.Four, '5': poker.Five,
	'6': poker.Six, '7': poker.Seven, '8': poker.Eight, '9': poker.Nine,
	'T': poker.Ten, 'J': poker.Jack, 'Q': poker.Queen, 'K': poker.King, 'A': poker.Ace,
}

var suitLetters = map[byte]poker.Suit{
	'H': poker.Hearts, 'D': poker.Diamonds, 'C': poker.Clubs, 'S': poker.Spades,
}

func parseCard(token string) (poker.Card, error) {
	if len(token) != 2 {
		return poker.Card{}, fmt.Errorf("invalid card token: %s", token)
	}
	upper := strings.ToUpper(token)
	rank, rankOK := rankLetters[upper[0]]
	suit, suitOK := suitLetters[upper[1]]
	if !rankOK || !suitOK {
		return poker.Card{}, fmt.Errorf("invalid card token: %s", token)
	}
	return poker.NewCard(rank, suit), nil
}

func main() {
	if len(os.Args) < 6 {
		fmt.Println("Usage: poker AH KH QH JH TH [more cards...]")
		fmt.Println("Ranks: 2-9,T,J,Q,K,A  Suits: H,D,C,S")
		return
	}

	cards := make([]poker.Card, 0, len(os.Args)-1)
	for _, arg := range os.Args[1:] {
		card, err := parseCard(arg)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		cards = append(cards, card)
	}

	result, err := poker.NewBestHandFinder().FindBestHand(cards)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("Cards:", cards)
	fmt.Println("Best hand:", result.Rank)
	fmt.Println("Tiebreakers:", result.Tiebreakers)
}

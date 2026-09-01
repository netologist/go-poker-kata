package poker

import "testing"

func TestFindBestHandWithExactlyFiveCards(t *testing.T) {
	hand := []Card{
		NewCard(Ace, Hearts), NewCard(King, Hearts), NewCard(Queen, Hearts),
		NewCard(Jack, Hearts), NewCard(Ten, Hearts),
	}
	result, err := NewBestHandFinder().FindBestHand(hand)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Rank != RoyalFlush {
		t.Errorf("want RoyalFlush, got %v", result.Rank)
	}
}

func TestFindBestHandWithSevenCardsPicksBestFive(t *testing.T) {
	// Two irrelevant "hole" cards plus five cards forming a flush.
	cards := []Card{
		NewCard(Two, Clubs), NewCard(Three, Spades),
		NewCard(King, Hearts), NewCard(Eight, Hearts), NewCard(Six, Hearts),
		NewCard(Four, Hearts), NewCard(Two, Hearts),
	}
	result, err := NewBestHandFinder().FindBestHand(cards)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Rank != Flush {
		t.Errorf("want Flush, got %v", result.Rank)
	}
}

func TestFindBestHandWithSixCardsFindsFourOfAKind(t *testing.T) {
	cards := []Card{
		NewCard(Nine, Clubs), NewCard(Nine, Hearts), NewCard(Nine, Spades),
		NewCard(Nine, Diamonds), NewCard(Two, Clubs), NewCard(Two, Hearts),
	}
	result, err := NewBestHandFinder().FindBestHand(cards)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Rank != FourOfAKind {
		t.Errorf("want FourOfAKind, got %v", result.Rank)
	}
}

func TestFindBestHandRejectsTooFewCards(t *testing.T) {
	cards := []Card{NewCard(Ace, Hearts), NewCard(King, Hearts)}
	_, err := NewBestHandFinder().FindBestHand(cards)
	if err == nil {
		t.Error("expected an error for fewer than 5 cards, got none")
	}
}

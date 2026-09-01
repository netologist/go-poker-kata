package poker

import (
	"reflect"
	"testing"
)

func TestRoyalFlush(t *testing.T) {
	hand := []Card{
		NewCard(Ace, Hearts), NewCard(King, Hearts), NewCard(Queen, Hearts),
		NewCard(Jack, Hearts), NewCard(Ten, Hearts),
	}
	result, err := NewHandEvaluator().Evaluate(hand)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Rank != RoyalFlush {
		t.Errorf("want RoyalFlush, got %v", result.Rank)
	}
}

func TestStraightFlush(t *testing.T) {
	hand := []Card{
		NewCard(Ten, Clubs), NewCard(Nine, Clubs), NewCard(Eight, Clubs),
		NewCard(Seven, Clubs), NewCard(Six, Clubs),
	}
	result, _ := NewHandEvaluator().Evaluate(hand)
	if result.Rank != StraightFlush {
		t.Errorf("want StraightFlush, got %v", result.Rank)
	}
	if !reflect.DeepEqual(result.Tiebreakers, []int{int(Ten)}) {
		t.Errorf("want tiebreaker [10], got %v", result.Tiebreakers)
	}
}

func TestWheelStraightFlushIsNotRoyal(t *testing.T) {
	hand := []Card{
		NewCard(Ace, Spades), NewCard(Five, Spades), NewCard(Four, Spades),
		NewCard(Three, Spades), NewCard(Two, Spades),
	}
	result, _ := NewHandEvaluator().Evaluate(hand)
	if result.Rank != StraightFlush {
		t.Errorf("want StraightFlush (wheel), got %v", result.Rank)
	}
	if !reflect.DeepEqual(result.Tiebreakers, []int{int(Five)}) {
		t.Errorf("want tiebreaker [5], got %v", result.Tiebreakers)
	}
}

func TestFourOfAKind(t *testing.T) {
	hand := []Card{
		NewCard(Queen, Clubs), NewCard(Queen, Hearts), NewCard(Queen, Spades),
		NewCard(Queen, Diamonds), NewCard(Five, Clubs),
	}
	result, _ := NewHandEvaluator().Evaluate(hand)
	if result.Rank != FourOfAKind {
		t.Errorf("want FourOfAKind, got %v", result.Rank)
	}
}

func TestFullHouse(t *testing.T) {
	hand := []Card{
		NewCard(Ace, Hearts), NewCard(Ace, Spades), NewCard(Ace, Diamonds),
		NewCard(King, Spades), NewCard(King, Hearts),
	}
	result, _ := NewHandEvaluator().Evaluate(hand)
	if result.Rank != FullHouse {
		t.Errorf("want FullHouse, got %v", result.Rank)
	}
}

func TestFlush(t *testing.T) {
	hand := []Card{
		NewCard(King, Hearts), NewCard(Eight, Hearts), NewCard(Six, Hearts),
		NewCard(Four, Hearts), NewCard(Two, Hearts),
	}
	result, _ := NewHandEvaluator().Evaluate(hand)
	if result.Rank != Flush {
		t.Errorf("want Flush, got %v", result.Rank)
	}
}

func TestStraight(t *testing.T) {
	hand := []Card{
		NewCard(Eight, Hearts), NewCard(Seven, Clubs), NewCard(Six, Diamonds),
		NewCard(Five, Spades), NewCard(Four, Hearts),
	}
	result, _ := NewHandEvaluator().Evaluate(hand)
	if result.Rank != Straight {
		t.Errorf("want Straight, got %v", result.Rank)
	}
}

func TestThreeOfAKind(t *testing.T) {
	hand := []Card{
		NewCard(Queen, Clubs), NewCard(Queen, Hearts), NewCard(Queen, Spades),
		NewCard(Seven, Hearts), NewCard(Two, Clubs),
	}
	result, _ := NewHandEvaluator().Evaluate(hand)
	if result.Rank != ThreeOfAKind {
		t.Errorf("want ThreeOfAKind, got %v", result.Rank)
	}
}

func TestTwoPair(t *testing.T) {
	hand := []Card{
		NewCard(Jack, Diamonds), NewCard(Jack, Spades), NewCard(Nine, Spades),
		NewCard(Nine, Diamonds), NewCard(Five, Clubs),
	}
	result, _ := NewHandEvaluator().Evaluate(hand)
	if result.Rank != TwoPair {
		t.Errorf("want TwoPair, got %v", result.Rank)
	}
}

func TestOnePair(t *testing.T) {
	hand := []Card{
		NewCard(King, Hearts), NewCard(King, Spades), NewCard(Nine, Diamonds),
		NewCard(Eight, Spades), NewCard(Four, Hearts),
	}
	result, _ := NewHandEvaluator().Evaluate(hand)
	if result.Rank != OnePair {
		t.Errorf("want OnePair, got %v", result.Rank)
	}
}

func TestHighCard(t *testing.T) {
	hand := []Card{
		NewCard(Ace, Hearts), NewCard(Queen, Clubs), NewCard(Six, Hearts),
		NewCard(Four, Spades), NewCard(Two, Diamonds),
	}
	result, _ := NewHandEvaluator().Evaluate(hand)
	if result.Rank != HighCard {
		t.Errorf("want HighCard, got %v", result.Rank)
	}
	if result.Tiebreakers[0] != int(Ace) {
		t.Errorf("want top tiebreaker Ace, got %v", result.Tiebreakers[0])
	}
}

func TestCategoryOrderingRespectedByCompare(t *testing.T) {
	evaluator := NewHandEvaluator()
	pair, _ := evaluator.Evaluate([]Card{
		NewCard(King, Hearts), NewCard(King, Spades), NewCard(Nine, Diamonds),
		NewCard(Eight, Spades), NewCard(Four, Hearts),
	})
	flush, _ := evaluator.Evaluate([]Card{
		NewCard(Two, Clubs), NewCard(Four, Clubs), NewCard(Six, Clubs),
		NewCard(Eight, Clubs), NewCard(Nine, Clubs),
	})
	if flush.Compare(pair) <= 0 {
		t.Errorf("expected flush to outrank pair")
	}
}

func TestEvaluateRejectsWrongHandSize(t *testing.T) {
	fourCards := []Card{
		NewCard(Ace, Hearts), NewCard(King, Hearts),
		NewCard(Queen, Hearts), NewCard(Jack, Hearts),
	}
	_, err := NewHandEvaluator().Evaluate(fourCards)
	if err == nil {
		t.Error("expected an error for a 4-card hand, got none")
	}
}

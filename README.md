# go-poker-kata

A **production-quality poker hand evaluator** built in idiomatic Go using object-oriented design principles — built step by step as a coding kata.

This project is a **design-patterns tutorial disguised as a poker engine**. Every file introduces exactly one idea — value objects, the Strategy pattern, the Open/Closed Principle, Dependency Inversion, composition — and this README walks through each decision in the exact order the code was built.

---

## Table of Contents

1. [What the Kata Asks For](#what-the-kata-asks-for)
2. [Architecture at a Glance](#architecture-at-a-glance)
3. [Libraries and Dependencies](#libraries-and-dependencies)
4. [Step 1 — Value Objects: `Card`, `Rank`, `Suit`](#step-1--value-objects-card-rank-suit)
5. [Step 2 — The Hand Taxonomy: `HandRank`](#step-2--the-hand-taxonomy-handrank)
6. [Step 3 — `HandResult`: Rank + Tiebreakers](#step-3--handresult-rank--tiebreakers)
7. [Step 4 — Pure Analysis Functions](#step-4--pure-analysis-functions)
8. [Step 5 — Strategy Pattern: `HandChecker` interface](#step-5--strategy-pattern-handchecker-interface)
9. [Step 6 — One Checker Per Hand Category](#step-6--one-checker-per-hand-category)
10. [Step 7 — `HandEvaluator`: Open/Closed + Dependency Inversion](#step-7--handevaluator-openclosed--dependency-inversion)
11. [Step 8 — `BestHandFinder`: Composition over Inheritance](#step-8--besthandfinder-composition-over-inheritance)
12. [Step 9 — CLI Entry Point](#step-9--cli-entry-point)
13. [Tiebreaker Logic Explained](#tiebreaker-logic-explained)
14. [The Wheel Straight Edge Case](#the-wheel-straight-edge-case)
15. [Extending the Kata](#extending-the-kata)
16. [Project Layout](#project-layout)
17. [Running the Project](#running-the-project)

---

## What the Kata Asks For

Given 5 (or more) playing cards, determine the **best 5-card poker hand** and correctly **rank two hands against each other** including tiebreakers.

The ten standard hand categories, from weakest to strongest:

| Rank | Name | Example |
|---|---|---|
| 1 | High Card | A♥ Q♣ 6♥ 4♠ 2♦ |
| 2 | One Pair | K♥ K♠ 9♦ 8♠ 4♥ |
| 3 | Two Pair | J♦ J♠ 9♠ 9♦ 5♣ |
| 4 | Three of a Kind | Q♣ Q♥ Q♠ 7♥ 2♣ |
| 5 | Straight | 8♥ 7♣ 6♦ 5♠ 4♥ |
| 6 | Flush | K♥ 8♥ 6♥ 4♥ 2♥ |
| 7 | Full House | A♥ A♠ A♦ K♠ K♥ |
| 8 | Four of a Kind | Q♣ Q♥ Q♠ Q♦ 5♣ |
| 9 | Straight Flush | T♣ 9♣ 8♣ 7♣ 6♣ |
| 10 | Royal Flush | A♥ K♥ Q♥ J♥ T♥ |

A complete solution must also handle:
- The **wheel straight** (A-2-3-4-5, Ace plays low, ranks as 5-high)
- **Tiebreakers** within the same category (e.g. pair of Kings beats pair of Jacks)
- **Best hand from N cards** (Texas Hold'em: pick best 5 from 7 hole + community cards)

---

## Architecture at a Glance

```
┌─────────────────────────────────────────────────────────┐
│                    main.go (CLI)                        │
│   parse args → []Card → BestHandFinder → print result  │
└─────────────────────────────┬───────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────┐
│                  BestHandFinder                         │
│   fiveCardCombinations([]Card) → try every 5-combo      │
│   ──────────────────────────────────────────────────    │
│   delegates each combo to → HandEvaluator               │
└─────────────────────────────┬───────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────┐
│                   HandEvaluator                         │
│   walks checkersStrongestFirst []HandChecker            │
│   first checker.Check() that returns ok=true wins       │
└────────────┬────────────────────────────────────────────┘
             │
             ▼ (each implements HandChecker interface)
┌────────────────────────────────────────────────────────────────┐
│  royalFlushChecker  straightFlushChecker  fourOfAKindChecker   │
│  fullHouseChecker   flushChecker          straightChecker       │
│  threeOfAKindChecker  twoPairChecker  onePairChecker           │
│  highCardChecker (fallback — always matches)                   │
└────────────────────────────────────────────────────────────────┘
             │
             ▼ (shared pure functions)
┌────────────────────────────────────────────────────────────────┐
│  analysis.go                                                   │
│  isFlush()  straightHighValueOrZero()  valueCountsByFrequency()│
│  sortedValuesDesc()  valuesOf()                                │
└────────────────────────────────────────────────────────────────┘
```

---

## Libraries and Dependencies

### Production

**Zero external dependencies.** `go.mod` contains only:

```
module github.com/netologist/go-poker-kata

go 1.22
```

Everything is the Go standard library: `sort`, `fmt`, `strings`, `os`.

> **Why no third-party libraries?**
> A kata is a pure algorithmic exercise — external libraries would obscure the design decisions.
> The only data structure is a slice; the only algorithm is a linear scan.
> `sort.Search` (binary search) and `sort.Slice` from the stdlib are the only non-trivial tools used.

### Testing

```go
import (
    "reflect"
    "testing"
)
```

`testing` and `reflect` from the standard library. Table-driven tests, no assertion frameworks. The test files live inside `package poker` (white-box testing) so they can access unexported helpers if needed.

---

## Step 1 — Value Objects: `Card`, `Rank`, `Suit`

> **file:** `poker/card.go`

The first decision: how to represent a playing card. A card has exactly two attributes — its face value (rank) and its suit. Both are **value objects**: immutable, comparable by value, with no identity beyond their attributes.

### Why typed integers, not strings?

```go
// Rank is a card's face value. The underlying int doubles as its numeric
// strength (2..14, Ace high).
type Rank int

const (
    Two   Rank = 2
    Three Rank = 3
    // ...
    Ace   Rank = 14
)
```

Using `Rank int` (not `type Rank string`) gives us:
- **Numeric comparison for free** — `King > Jack` is just `13 > 11`
- **Zero-cost tiebreaker encoding** — `int(card.Rank)` directly becomes the tiebreaker value
- **Type safety** — the compiler rejects `card.Rank + card.Suit` (different types)

The `Two = 2` assignment is deliberate: lower cards start at 2 (not 0 or 1), so `int(Rank)` is always the correct poker value with no translation.

```go
type Suit int

const (
    Hearts Suit = iota
    Diamonds
    Clubs
    Spades
)
```

Suits use `iota` because their numeric value is irrelevant — suits are only compared for equality (flush detection), never ordered.

### The `Card` struct

```go
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
```

`Card` is a value type (struct, not pointer). Slices of cards are contiguous in memory — no heap allocation per card. The constructor `NewCard` is provided for readability in tests, not because it adds logic.

### `String()` methods — implementing `fmt.Stringer`

```go
func (r Rank) String() string {
    names := map[Rank]string{
        Two: "2", Three: "3", ..., Jack: "J", Queen: "Q", King: "K", Ace: "A",
    }
    if name, ok := names[r]; ok {
        return name
    }
    return fmt.Sprintf("Rank(%d)", int(r))
}
```

Implementing `fmt.Stringer` on `Rank`, `Suit`, and `Card` means every `fmt.Println(card)` produces a readable output like `AH` (Ace of Hearts) — essential for test diagnostics and the CLI output.

---

## Step 2 — The Hand Taxonomy: `HandRank`

> **file:** `poker/handrank.go`

```go
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
```

`iota + 1` starts at 1, not 0. This matters: when two hands have the same `HandRank`, their comparison returns 0. If `HighCard = 0`, the zero value of `HandResult` would silently look like a valid High Card, masking bugs.

The numeric ordering directly encodes hand strength — `int(RoyalFlush) > int(StraightFlush) > ... > int(HighCard)` — so comparing two `HandRank` values with subtraction gives the correct sign:

```go
// In HandResult.Compare:
if h.Rank != other.Rank {
    return int(h.Rank) - int(other.Rank)  // positive = h is stronger
}
```

---

## Step 3 — `HandResult`: Rank + Tiebreakers

> **file:** `poker/handresult.go`

The result of evaluating a hand is not just its category — two players can both have a Flush and one still wins. The tiebreaker sequence encodes *which specific flush wins*.

```go
type HandResult struct {
    Rank        HandRank
    Tiebreakers []int  // most significant first
}
```

### Tiebreaker ordering by hand type

| Hand | Tiebreakers |
|---|---|
| High Card | `[14, 12, 6, 4, 2]` — all five ranks, highest first |
| One Pair | `[13, 9, 8, 4]` — pair rank, then kickers highest-first |
| Two Pair | `[11, 9, 5]` — high pair, low pair, kicker |
| Three of a Kind | `[12, 7, 2]` — trips rank, kickers |
| Straight | `[8]` — high card of the straight |
| Flush | `[13, 8, 6, 4, 2]` — all five ranks, highest first |
| Full House | `[14, 13]` — trips rank, pair rank |
| Four of a Kind | `[12, 5]` — quads rank, kicker |
| Straight Flush | `[10]` — high card of the straight |
| Royal Flush | `[14]` — always Ace, always ties |

### `Compare` — lexicographic ordering of tiebreakers

```go
func (h HandResult) Compare(other HandResult) int {
    // 1. Different categories — the higher category wins outright
    if h.Rank != other.Rank {
        return int(h.Rank) - int(other.Rank)
    }
    // 2. Same category — compare tiebreakers left to right
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
```

The return sign convention is: **negative = receiver is weaker, zero = equal, positive = receiver is stronger**. This matches the `sort.Interface.Less` pattern and enables: `results[i].Compare(results[j]) > 0` means `i` beats `j`.

---

## Step 4 — Pure Analysis Functions

> **file:** `poker/analysis.go`

Before any hand-recognition logic, we extract the pure mathematical primitives that every checker needs. These are package-level functions with no state — easy to test in isolation, impossible to misuse.

### `sortedValuesDesc` — rank values, highest first

```go
func sortedValuesDesc(cards []Card) []int {
    values := make([]int, len(cards))
    for i, c := range cards {
        values[i] = int(c.Rank)
    }
    sort.Sort(sort.Reverse(sort.IntSlice(values)))
    return values
}
```

Used by `flushChecker` and `highCardChecker` — their tiebreakers are simply all five rank values in descending order.

### `valueCountsByFrequency` — the pairs/trips/quads engine

```go
type rankValueCount struct {
    value int
    count int
}

func valueCountsByFrequency(cards []Card) []rankValueCount {
    counts := map[int]int{}
    for _, c := range cards {
        counts[int(c.Rank)]++
    }
    result := make([]rankValueCount, 0, len(counts))
    for value, count := range counts {
        result = append(result, rankValueCount{value: value, count: count})
    }
    // Primary sort: frequency DESC (quads before trips before pairs before kickers)
    // Secondary sort: rank value DESC (higher pair beats lower pair)
    sort.Slice(result, func(i, j int) bool {
        if result[i].count != result[j].count {
            return result[i].count > result[j].count
        }
        return result[i].value > result[j].value
    })
    return result
}
```

This dual-sort gives the tiebreakers in the right order automatically for any grouped hand. For a full house (A-A-A-K-K):
- Aces: count=3, value=14 → first
- Kings: count=2, value=13 → second
- Tiebreakers: `[14, 13]` — correct.

For two pair (J-J-9-9-5):
- Jacks: count=2, value=11 → first
- Nines: count=2, value=9 → second
- Five: count=1, value=5 → third
- Tiebreakers: `[11, 9, 5]` — correct.

### `isFlush` — same suit check

```go
func isFlush(cards []Card) bool {
    suit := cards[0].Suit
    for _, c := range cards[1:] {
        if c.Suit != suit {
            return false
        }
    }
    return true
}
```

Short-circuits on the first mismatch. O(n) with early exit.

### `straightHighValueOrZero` — straight detection with wheel support

```go
func straightHighValueOrZero(cards []Card) int {
    // Collect distinct values (a duplicate immediately disqualifies a straight)
    seen := map[int]bool{}
    distinct := make([]int, 0, 5)
    for _, c := range cards {
        v := int(c.Rank)
        if !seen[v] {
            seen[v] = true
            distinct = append(distinct, v)
        }
    }
    if len(distinct) != 5 {
        return 0  // has a pair — cannot be a straight
    }
    sort.Sort(sort.Reverse(sort.IntSlice(distinct)))

    // Normal straight: max - min == 4 (e.g. T-9-8-7-6: 10-6=4)
    if distinct[0]-distinct[4] == 4 {
        return distinct[0]
    }
    // Wheel: A-5-4-3-2 — Ace plays as 1
    wheel := []int{int(Ace), int(Five), int(Four), int(Three), int(Two)}
    for i, v := range wheel {
        if distinct[i] != v {
            return 0
        }
    }
    return int(Five)  // Wheel is a 5-high straight, not 14-high
}
```

See [The Wheel Straight Edge Case](#the-wheel-straight-edge-case) for the full explanation.

---

## Step 5 — Strategy Pattern: `HandChecker` interface

> **file:** `poker/checker.go`

```go
// HandChecker recognises a single hand category. Implementations are
// stateless and only know how to detect their own category, so adding a new
// hand type never requires touching an existing checker (Open/Closed
// Principle). ok is false when the five cards do not match this category.
type HandChecker interface {
    Check(fiveCards []Card) (result HandResult, ok bool)
}
```

This is the **Strategy pattern**: a family of algorithms (hand recognition strategies) unified behind a single interface.

**Why a two-value return `(HandResult, bool)` instead of `(HandResult, error)`?**

"This hand is not a flush" is not an error — it's the normal, expected outcome for most checkers. Using `bool` (the idiomatic Go "ok" pattern) avoids allocating an error value on the hot path (every checker is tried for every hand).

**Why is `HandChecker` defined separately from its implementations?**

The interface is the contract. Placing it in its own file (`checker.go`) makes the contract immediately visible without wading through implementations, and it lives at the level of the concept, not any particular implementation.

---

## Step 6 — One Checker Per Hand Category

> **file:** `poker/checkers.go`

Each checker is a zero-size struct implementing `HandChecker`. Zero-size structs cost no memory — 10 checker instances allocate nothing on the heap:

```go
type royalFlushChecker struct{}

func (royalFlushChecker) Check(cards []Card) (HandResult, bool) {
    if isFlush(cards) && straightHighValueOrZero(cards) == int(Ace) {
        return HandResult{Rank: RoyalFlush, Tiebreakers: []int{int(Ace)}}, true
    }
    return HandResult{}, false
}
```

**Royal Flush** = flush AND straight with Ace as the high card. The Ace tiebreaker is constant — all Royal Flushes are equal.

```go
type straightFlushChecker struct{}

func (straightFlushChecker) Check(cards []Card) (HandResult, bool) {
    high := straightHighValueOrZero(cards)
    if isFlush(cards) && high != 0 {
        return HandResult{Rank: StraightFlush, Tiebreakers: []int{high}}, true
    }
    return HandResult{}, false
}
```

**Straight Flush** = flush AND any straight. Note that `royalFlushChecker` runs *before* `straightFlushChecker` in the evaluator — the Royal Flush would also satisfy Straight Flush's condition. Order in the checker list is what separates them.

```go
type fourOfAKindChecker struct{}

func (fourOfAKindChecker) Check(cards []Card) (HandResult, bool) {
    counts := valueCountsByFrequency(cards)
    if counts[0].count == 4 {
        return HandResult{Rank: FourOfAKind, Tiebreakers: valuesOf(counts)}, true
    }
    return HandResult{}, false
}
```

`valueCountsByFrequency` sorts by frequency then value, so `counts[0]` is always the most frequent group. If that group has count 4, it's four-of-a-kind. The tiebreakers are `[quads_rank, kicker_rank]`.

```go
type fullHouseChecker struct{}

func (fullHouseChecker) Check(cards []Card) (HandResult, bool) {
    counts := valueCountsByFrequency(cards)
    // Exactly 2 distinct ranks (trips + pair), and the top group is trips
    if len(counts) == 2 && counts[0].count == 3 {
        return HandResult{Rank: FullHouse, Tiebreakers: valuesOf(counts)}, true
    }
    return HandResult{}, false
}
```

`len(counts) == 2` means only two distinct ranks in the five cards. That could be four-of-a-kind + kicker (4+1) or full house (3+2). `counts[0].count == 3` distinguishes the full house.

```go
type flushChecker struct{}

func (flushChecker) Check(cards []Card) (HandResult, bool) {
    if isFlush(cards) {
        return HandResult{Rank: Flush, Tiebreakers: sortedValuesDesc(cards)}, true
    }
    return HandResult{}, false
}
```

A plain flush (not a straight flush, because `straightFlushChecker` runs first). All five rank values in descending order form the tiebreakers — K♥8♥6♥4♥2♥ beats K♥8♥6♥4♥3♥ because the fourth tiebreaker is 4 vs 4 (equal), then 2 vs 3 — wait, `sortedValuesDesc` puts them highest-first: `[13,8,6,4,3]` beats `[13,8,6,4,2]` because the fifth element `3 > 2`. ✓

```go
type straightChecker struct{}

func (straightChecker) Check(cards []Card) (HandResult, bool) {
    high := straightHighValueOrZero(cards)
    if high != 0 {
        return HandResult{Rank: Straight, Tiebreakers: []int{high}}, true
    }
    return HandResult{}, false
}
```

A straight is fully described by its high card. Two straights with the same high card are always equal (they have the same five consecutive ranks).

```go
type threeOfAKindChecker struct{}

func (threeOfAKindChecker) Check(cards []Card) (HandResult, bool) {
    counts := valueCountsByFrequency(cards)
    // 3 distinct ranks (trips + 2 different kickers), top group is trips
    if len(counts) == 3 && counts[0].count == 3 {
        return HandResult{Rank: ThreeOfAKind, Tiebreakers: valuesOf(counts)}, true
    }
    return HandResult{}, false
}
```

`len(counts) == 3` means three distinct ranks. That could be three-of-a-kind (3+1+1) or two pair (2+2+1). `counts[0].count == 3` distinguishes trips.

```go
type twoPairChecker struct{}

func (twoPairChecker) Check(cards []Card) (HandResult, bool) {
    counts := valueCountsByFrequency(cards)
    // 3 distinct ranks, top two groups are both pairs
    if len(counts) == 3 && counts[0].count == 2 && counts[1].count == 2 {
        return HandResult{Rank: TwoPair, Tiebreakers: valuesOf(counts)}, true
    }
    return HandResult{}, false
}
```

```go
type onePairChecker struct{}

func (onePairChecker) Check(cards []Card) (HandResult, bool) {
    counts := valueCountsByFrequency(cards)
    // 4 distinct ranks (pair + 3 kickers)
    if len(counts) == 4 && counts[0].count == 2 {
        return HandResult{Rank: OnePair, Tiebreakers: valuesOf(counts)}, true
    }
    return HandResult{}, false
}
```

```go
// highCardChecker is the fallback: any five cards form at least a high card hand.
type highCardChecker struct{}

func (highCardChecker) Check(cards []Card) (HandResult, bool) {
    return HandResult{Rank: HighCard, Tiebreakers: sortedValuesDesc(cards)}, true
}
```

`highCardChecker` **always** returns `ok=true`. It must be the last checker in the list — it's the catch-all that guarantees every five-card hand gets a classification.

### The checker dispatch table

All ten hand categories mapped to their distinguishing condition:

| Category | `len(counts)` | `counts[0].count` | Extra condition |
|---|---|---|---|
| Royal Flush | — | — | `isFlush` AND `high==Ace` |
| Straight Flush | — | — | `isFlush` AND `high!=0` |
| Four of a Kind | 2 | 4 | — |
| Full House | 2 | 3 | — |
| Flush | — | — | `isFlush` |
| Straight | — | — | `high!=0` |
| Three of a Kind | 3 | 3 | — |
| Two Pair | 3 | 2 | `counts[1].count==2` |
| One Pair | 4 | 2 | — |
| High Card | 5 | 1 | always |

---

## Step 7 — `HandEvaluator`: Open/Closed + Dependency Inversion

> **file:** `poker/evaluator.go`

```go
// HandEvaluator evaluates an exact 5-card hand into its HandResult by trying
// a list of HandChecker implementations from strongest to weakest category;
// the first match wins. It depends only on the HandChecker interface
// (Dependency Inversion Principle) and the checker list can be extended
// without modifying this type (Open/Closed Principle).
type HandEvaluator struct {
    checkersStrongestFirst []HandChecker
}
```

### Open/Closed Principle

The evaluator is **open for extension** (add a new hand category by appending a new `HandChecker` implementation) and **closed for modification** (adding it requires zero changes to `HandEvaluator` itself).

Want to add a "Five of a Kind" category for a wild-card variant? Write `fiveOfAKindChecker{}`, prepend it to the list — done.

### Dependency Inversion Principle

`HandEvaluator` depends on the *abstraction* (`HandChecker` interface), not any concrete checker. This enables:

```go
// Production — all ten standard checkers
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

// Test / house rules — inject only the checkers you care about
func NewHandEvaluatorWithCheckers(checkers []HandChecker) *HandEvaluator {
    copied := make([]HandChecker, len(checkers))
    copy(copied, checkers)  // defensive copy — caller cannot mutate the evaluator's list
    return &HandEvaluator{checkersStrongestFirst: copied}
}
```

### The evaluation loop

```go
func (e *HandEvaluator) Evaluate(fiveCards []Card) (HandResult, error) {
    if len(fiveCards) != 5 {
        return HandResult{}, fmt.Errorf(
            "a poker hand must contain exactly 5 cards, got %d", len(fiveCards))
    }
    for _, checker := range e.checkersStrongestFirst {
        if result, ok := checker.Check(fiveCards); ok {
            return result, nil  // first match wins — strongest category first
        }
    }
    // Unreachable: highCardChecker always matches.
    return HandResult{}, fmt.Errorf("no checker matched hand: %v", fiveCards)
}
```

The loop is **O(n)** where n = number of checkers (10). Each checker call is O(5) = O(1). Total: O(10 × 5) = O(1). The evaluator is effectively constant time.

---

## Step 8 — `BestHandFinder`: Composition over Inheritance

> **file:** `poker/besthand.go`

`BestHandFinder` solves a different problem from `HandEvaluator`: given *more* than 5 cards, find the best possible 5-card hand. In Texas Hold'em, each player sees 7 cards (2 hole + 5 community) and plays the best 5.

```go
// BestHandFinder finds the highest-ranked 5-card poker hand obtainable from
// a set of 5 or more cards (e.g. picking the best 5 out of 7 hole+community
// cards). It composes a HandEvaluator rather than embedding its logic,
// keeping combination-generation and hand-ranking as separate responsibilities.
type BestHandFinder struct {
    evaluator *HandEvaluator
}
```

### Composition over inheritance

`BestHandFinder` *has* a `HandEvaluator` — it does not *extend* it. If Go had class inheritance, you might be tempted to subclass `HandEvaluator`. Composition is cleaner:
- `BestHandFinder` adds one responsibility (combination enumeration) without touching evaluation logic
- `HandEvaluator` can be swapped (different house rules) via the constructor

```go
func NewBestHandFinder() *BestHandFinder {
    return &BestHandFinder{evaluator: NewHandEvaluator()}
}

func NewBestHandFinderWithEvaluator(evaluator *HandEvaluator) *BestHandFinder {
    return &BestHandFinder{evaluator: evaluator}
}
```

### Combination generation

For 7 cards, C(7,5) = 21 combinations. For 6 cards, C(6,5) = 6. The recursive generator uses a reusable slice to avoid allocating a new slice per recursion level:

```go
func fiveCardCombinations(cards []Card) [][]Card {
    var results [][]Card
    current := make([]Card, 0, 5)
    var combine func(start int)
    combine = func(start int) {
        if len(current) == 5 {
            combo := make([]Card, 5)
            copy(combo, current)  // snapshot — current will be modified
            results = append(results, combo)
            return
        }
        for i := start; i < len(cards); i++ {
            current = append(current, cards[i])
            combine(i + 1)
            current = current[:len(current)-1]  // backtrack
        }
    }
    combine(0)
    return results
}
```

The `copy` on each complete combination is essential — without it, all combinations would point to the same backing array, and only the last combination's state would survive.

### `FindBestHand` — linear scan for maximum

```go
func (f *BestHandFinder) FindBestHand(cards []Card) (HandResult, error) {
    if len(cards) < 5 {
        return HandResult{}, fmt.Errorf("at least 5 cards are required, got %d", len(cards))
    }
    if len(cards) == 5 {
        return f.evaluator.Evaluate(cards)  // fast path — no combinations needed
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
```

`haveBest` avoids comparing against the zero value of `HandResult`, which would be `HandRank(0)` — less than any real hand, but semantically undefined.

---

## Step 9 — CLI Entry Point

> **file:** `main.go`

```go
// Command poker parses cards given as command line arguments in the form
// "RANKSUIT" (e.g. "AH KH QH JH TH") and prints the best 5-card hand found.
// Rank letters: 2-9,T,J,Q,K,A. Suit letters: H,D,C,S.
package main
```

The card parser is a simple lookup table:

```go
var rankLetters = map[byte]poker.Rank{
    '2': poker.Two, '3': poker.Three, '4': poker.Four, '5': poker.Five,
    '6': poker.Six, '7': poker.Seven, '8': poker.Eight, '9': poker.Nine,
    'T': poker.Ten, 'J': poker.Jack, 'Q': poker.Queen, 'K': poker.King, 'A': poker.Ace,
}

var suitLetters = map[byte]poker.Suit{
    'H': poker.Hearts, 'D': poker.Diamonds, 'C': poker.Clubs, 'S': poker.Spades,
}
```

`'T'` for Ten (not `'1'` or `"10"`) is the standard poker shorthand that fits the 2-character token format.

```bash
# Royal Flush
go run . AH KH QH JH TH
# Cards: [AH KH QH JH TH]
# Best hand: Royal Flush
# Tiebreakers: [14]

# Texas Hold'em: 7 cards, picks the best 5
go run . 2C 3S KH 8H 6H 4H 2H
# Cards: [2C 3S KH 8H 6H 4H 2H]
# Best hand: Flush
# Tiebreakers: [13 8 6 4 2]
```

---

## Tiebreaker Logic Explained

Every `HandResult` carries an ordered `Tiebreakers []int` slice. `Compare` walks this slice lexicographically — the first position where they differ determines the winner.

### Example: Two Pair vs Two Pair

```
Hand A: J♦ J♠ 9♠ 9♦ 5♣  → Tiebreakers: [11, 9, 5]
Hand B: J♣ J♥ 8♦ 8♣ K♥  → Tiebreakers: [11, 8, 13]

Compare position 0: 11 == 11 → tie, continue
Compare position 1: 9  >  8  → Hand A wins
```

Hand A wins because its second pair (Nines) beats Hand B's second pair (Eights). The King kicker in Hand B is irrelevant once position 1 differs.

### Example: Full House vs Full House

```
Hand A: A♥ A♠ A♦ K♠ K♥  → Tiebreakers: [14, 13]
Hand B: K♣ K♥ K♦ A♠ A♦  → Tiebreakers: [13, 14]

Compare position 0: 14 > 13 → Hand A wins
```

Aces-full-of-Kings beats Kings-full-of-Aces. The trips rank is always the first tiebreaker, which is why `valueCountsByFrequency` sorts by frequency (descending) before value.

---

## The Wheel Straight Edge Case

The **wheel** (`A-2-3-4-5`) is the lowest possible straight. The Ace plays as a 1, not a 14. This creates a special case because in all other contexts `Ace = 14`.

```go
func straightHighValueOrZero(cards []Card) int {
    // ... build `distinct` (5 unique values sorted DESC) ...

    // Normal check: consecutive range of 5
    if distinct[0]-distinct[4] == 4 {
        return distinct[0]
    }

    // Wheel check: exactly [14, 5, 4, 3, 2]
    wheel := []int{int(Ace), int(Five), int(Four), int(Three), int(Two)}
    for i, v := range wheel {
        if distinct[i] != v {
            return 0
        }
    }
    return int(Five)  // ← return 5, not 14
}
```

Returning `int(Five) = 5` means:
- The wheel compares as a **5-high straight**, below any 6-high or higher straight
- `A-2-3-4-5♠ (flush)` correctly becomes a **5-high Straight Flush**, not an Ace-high one

Without this special case, `A♠-2♠-3♠-4♠-5♠` would incorrectly be classified as a 14-high Straight Flush — stronger than `K-high` — which is wrong by poker rules.

The test `TestWheelStraightFlushIsNotRoyal` explicitly verifies this:

```go
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
```

---

## Extending the Kata

### Adding a new hand category (e.g. Five of a Kind for wild-card games)

1. Add the constant to `handrank.go`:
```go
const (
    HighCard HandRank = iota + 1
    // ... existing values ...
    RoyalFlush
    FiveOfAKind  // ← add here, stronger than Royal Flush
)
```

2. Write the checker in `checkers.go`:
```go
type fiveOfAKindChecker struct{}

func (fiveOfAKindChecker) Check(cards []Card) (HandResult, bool) {
    counts := valueCountsByFrequency(cards)
    if counts[0].count == 5 {
        return HandResult{Rank: FiveOfAKind, Tiebreakers: valuesOf(counts)}, true
    }
    return HandResult{}, false
}
```

3. Prepend it when building the evaluator:
```go
return NewHandEvaluatorWithCheckers([]HandChecker{
    fiveOfAKindChecker{},  // ← new, must come before royalFlushChecker
    royalFlushChecker{},
    // ... rest unchanged ...
})
```

Zero changes to `HandEvaluator`, `BestHandFinder`, or any existing checker. This is the Open/Closed Principle in practice.

### Testing a single checker in isolation

```go
func TestMyChecker(t *testing.T) {
    evaluator := NewHandEvaluatorWithCheckers([]HandChecker{
        myChecker{},
        highCardChecker{}, // fallback
    })
    result, _ := evaluator.Evaluate(testHand)
    // ...
}
```

### Hand comparison for a showdown

```go
func showdown(hands [][]poker.Card) (int, error) {
    finder := poker.NewBestHandFinder()
    results := make([]poker.HandResult, len(hands))
    for i, hand := range hands {
        r, err := finder.FindBestHand(hand)
        if err != nil {
            return -1, err
        }
        results[i] = r
    }
    winner := 0
    for i := 1; i < len(results); i++ {
        if results[i].Compare(results[winner]) > 0 {
            winner = i
        }
    }
    return winner, nil
}
```

---

## Project Layout

```
go-poker-kata/
├── go.mod                   Module declaration (zero external deps)
├── main.go                  CLI entry point — parse "AH KH QH JH TH" → best hand
└── poker/
    ├── card.go              Value objects: Card, Rank (int 2-14), Suit
    ├── handrank.go          HandRank enum (HighCard=1 … RoyalFlush=10)
    ├── handresult.go        HandResult{Rank, Tiebreakers} + Compare()
    ├── analysis.go          Pure functions: isFlush, straightHighValueOrZero, valueCountsByFrequency
    ├── checker.go           HandChecker interface (Strategy pattern)
    ├── checkers.go          10 concrete checker implementations (one per hand category)
    ├── evaluator.go         HandEvaluator — walks checker chain strongest-first
    ├── besthand.go          BestHandFinder — C(n,5) combinations + best result
    ├── evaluator_test.go    17 tests: one per category + edge cases + Compare
    └── besthand_test.go     4 tests: exact 5, 7-card best, 6-card best, too few cards
```

---

## Running the Project

```bash
git clone https://github.com/netologist/go-poker-kata
cd go-poker-kata

# Run all tests
go test ./...

# Run with verbose output
go test ./poker/... -v

# Run with race detector
go test ./... -race

# Run the CLI
go run . AH KH QH JH TH           # Royal Flush
go run . TC 9C 8C 7C 6C           # Straight Flush (Ten-high)
go run . AC 5S 4S 3S 2S           # Wheel Straight Flush (Five-high)
go run . QC QH QS QD 5C           # Four of a Kind
go run . AH AS AD KS KH           # Full House (Aces full of Kings)
go run . KH 8H 6H 4H 2H           # Flush
go run . 8H 7C 6D 5S 4H           # Straight (Eight-high)
go run . QC QH QS 7H 2C           # Three of a Kind
go run . JD JS 9S 9D 5C           # Two Pair
go run . KH KS 9D 8S 4H           # One Pair
go run . AH QC 6H 4S 2D           # High Card

# Texas Hold'em: 7 cards, picks best 5
go run . 2C 3S KH 8H 6H 4H 2H
```

### Test output

```
=== RUN   TestRoyalFlush
--- PASS: TestRoyalFlush (0.00s)
=== RUN   TestStraightFlush
--- PASS: TestStraightFlush (0.00s)
=== RUN   TestWheelStraightFlushIsNotRoyal
--- PASS: TestWheelStraightFlushIsNotRoyal (0.00s)
=== RUN   TestFourOfAKind
--- PASS: TestFourOfAKind (0.00s)
=== RUN   TestFullHouse
--- PASS: TestFullHouse (0.00s)
=== RUN   TestFlush
--- PASS: TestFlush (0.00s)
=== RUN   TestStraight
--- PASS: TestStraight (0.00s)
=== RUN   TestThreeOfAKind
--- PASS: TestThreeOfAKind (0.00s)
=== RUN   TestTwoPair
--- PASS: TestTwoPair (0.00s)
=== RUN   TestOnePair
--- PASS: TestOnePair (0.00s)
=== RUN   TestHighCard
--- PASS: TestHighCard (0.00s)
=== RUN   TestCategoryOrderingRespectedByCompare
--- PASS: TestCategoryOrderingRespectedByCompare (0.00s)
=== RUN   TestEvaluateRejectsWrongHandSize
--- PASS: TestEvaluateRejectsWrongHandSize (0.00s)
=== RUN   TestFindBestHandWithExactlyFiveCards
--- PASS: TestFindBestHandWithExactlyFiveCards (0.00s)
=== RUN   TestFindBestHandWithSevenCardsPicksBestFive
--- PASS: TestFindBestHandWithSevenCardsPicksBestFive (0.00s)
=== RUN   TestFindBestHandWithSixCardsFindsFourOfAKind
--- PASS: TestFindBestHandWithSixCardsFindsFourOfAKind (0.00s)
=== RUN   TestFindBestHandRejectsTooFewCards
--- PASS: TestFindBestHandRejectsTooFewCards (0.00s)
PASS
ok      github.com/netologist/go-poker-kata/poker  0.002s
```

---

## Key Design Decisions

### 1. Typed integer ranks, numeric strength baked in

`Rank int` with `Two=2 … Ace=14` means `int(card.Rank)` is the correct poker value everywhere. No translation table, no off-by-one errors.

### 2. `HandChecker` interface — Strategy pattern

Ten hand categories → ten independent strategies. Each is a zero-size struct. Adding a category means adding one struct and one entry in the constructor — nothing else changes.

### 3. Checker order is the disambiguation mechanism

Royal Flush and Straight Flush overlap. Four of a Kind and Full House both have `len(counts)==2`. Running checkers strongest-first and returning the first match means the evaluation logic itself needs no special disambiguation logic.

### 4. `valueCountsByFrequency` — one sort, many uses

One dual-keyed sort (frequency DESC, value DESC) produces the correct tiebreaker ordering for all grouped hands (pairs, trips, quads, full house, two pair) without any hand-specific post-processing.

### 5. Zero external dependencies

The entire implementation is standard library. `sort.Slice`, `sort.Search`, maps, slices — the building blocks are readable and auditable by any Go developer. No hidden framework magic.

### 6. Value types for `Card` and `HandResult`

Both are structs passed by value. Slices of cards are contiguous in memory. Comparison via `HandResult.Compare` is pure arithmetic on `int` fields — no heap allocations, no pointer chasing.

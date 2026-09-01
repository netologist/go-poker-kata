package poker

import "sort"

// rankValueCount pairs a rank value with how many times it occurs.
type rankValueCount struct {
	value int
	count int
}

// sortedValuesDesc returns the numeric rank values of the given cards,
// sorted highest first.
func sortedValuesDesc(cards []Card) []int {
	values := make([]int, len(cards))
	for i, c := range cards {
		values[i] = int(c.Rank)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(values)))
	return values
}

// valueCountsByFrequency groups cards by rank value and orders the groups by
// (count desc, rank value desc) - directly usable as tiebreaker ordering for
// pair/trips/quads style hands.
func valueCountsByFrequency(cards []Card) []rankValueCount {
	counts := map[int]int{}
	for _, c := range cards {
		counts[int(c.Rank)]++
	}
	result := make([]rankValueCount, 0, len(counts))
	for value, count := range counts {
		result = append(result, rankValueCount{value: value, count: count})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].count != result[j].count {
			return result[i].count > result[j].count
		}
		return result[i].value > result[j].value
	})
	return result
}

func isFlush(cards []Card) bool {
	suit := cards[0].Suit
	for _, c := range cards[1:] {
		if c.Suit != suit {
			return false
		}
	}
	return true
}

// straightHighValueOrZero returns the straight's high value (treating the
// Ace-low "wheel" as high=5), or 0 if the five cards do not form a straight.
func straightHighValueOrZero(cards []Card) int {
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
		return 0
	}
	sort.Sort(sort.Reverse(sort.IntSlice(distinct)))
	if distinct[0]-distinct[4] == 4 {
		return distinct[0]
	}
	wheel := []int{int(Ace), int(Five), int(Four), int(Three), int(Two)}
	for i, v := range wheel {
		if distinct[i] != v {
			return 0
		}
	}
	return int(Five)
}

func valuesOf(counts []rankValueCount) []int {
	values := make([]int, len(counts))
	for i, rc := range counts {
		values[i] = rc.value
	}
	return values
}

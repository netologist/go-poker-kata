package poker

// HandChecker recognises a single hand category. Implementations are
// stateless and only know how to detect their own category, so adding a new
// hand type never requires touching an existing checker (Open/Closed
// principle). ok is false when the five cards do not match this category.
type HandChecker interface {
	Check(fiveCards []Card) (result HandResult, ok bool)
}

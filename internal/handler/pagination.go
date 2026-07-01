package handler

// Pagination constants for consistent pagination across all handlers.
//
// Default limits:
//   - Courses & Blogs: 10 (content-heavy, typically fewer items)
//   - Leaderboard:     10 (top-N display)
//   - Everything else: 20 (general-purpose listing)
//
// MaxLimit applies to all endpoints to prevent abuse (e.g. ?limit=999999).

const (
	DefaultPage       = 1
	DefaultLimit      = 20
	DefaultCourseLimit = 10
	DefaultBlogLimit  = 10
	DefaultLeaderLimit = 10
	MaxLimit          = 100
)

// clampLimit ensures limit is within [1, MaxLimit].
// If limit is <= 0, it falls back to the provided default.
func clampLimit(limit, defaultVal int) int {
	if limit < 1 {
		return defaultVal
	}
	if limit > MaxLimit {
		return MaxLimit
	}
	return limit
}

// clampPage ensures page is at least 1.
func clampPage(page int) int {
	if page < 1 {
		return DefaultPage
	}
	return page
}

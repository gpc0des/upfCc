// package types defines the types that are globally used in the application.
package types

// ValidDimensions defines the valid dimensions for social media interactions.
// A dimension is a type of social media interaction, such as likes, comments, favorites, or retweets.
type Dimension string

const (
	Likes Dimension = "likes"
	Comments Dimension = "comments"
	Favorites Dimension = "favorites"
	Retweets Dimension = "retweets"
)

// IsValidDimension verifies if the given dimension is valid.
func IsValidDimension(dimension Dimension) bool {
	switch dimension {
	case Likes, Comments, Favorites, Retweets:
		return true
	default:
		return false
	}
}
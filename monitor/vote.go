package monitor

import "math/rand"

const (
	// SignLike for voting
	SignLike = "👍"
	// SignDislike for voting
	SignDislike = "👎"
)

// GetVote for voting
func GetVote() string {
	signs := []string{
		SignLike,
		SignDislike,
	}

	return signs[rand.Intn(len(signs))]
}

package monitor

import (
	"math/rand"
	"time"
)

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

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	return signs[r.Intn(len(signs))]
}

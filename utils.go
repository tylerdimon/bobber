package bobber

import (
	"github.com/google/uuid"
	"time"
)

const MonotonicClock = 1

// Generator this interface allows for easy mocking of ID and timestamp generation in tests
type Generator struct {
	Now  func() time.Time
	UUID func() uuid.UUID
}

func NewGenerator() Generator {
	return Generator{
		Now: func() time.Time {
			// the monotonic clock causes issues when converting to and from strings
			return time.Now().Truncate(MonotonicClock)
		},
		UUID: uuid.New,
	}
}

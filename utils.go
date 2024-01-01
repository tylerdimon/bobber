package bobber

import (
	"github.com/google/uuid"
	"time"
)

const MonotonicClock = 1

// this Generator interface allows for easy mocking in tests
type Generator struct {
	Now  func() time.Time
	UUID func() uuid.UUID
}

func GetGenerator() Generator {
	return Generator{
		Now: func() time.Time {
			return time.Now().Truncate(MonotonicClock)
		},
		UUID: uuid.New,
	}
}

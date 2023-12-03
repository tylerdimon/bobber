package bobber

import (
	"github.com/google/uuid"
	"time"
)

// this Generator interface allows for easy mocking in tests

type Generator struct {
	Now  func() time.Time
	UUID func() uuid.UUID
}

func GetGenerator() Generator {
	return Generator{
		Now:  time.Now,
		UUID: uuid.New,
	}
}

package mock

import (
	"github.com/google/uuid"
	"github.com/tylerdimon/bobber"
	"time"
)

const UUIDString = "6e300e63-3b0a-470e-b169-f4460e1ccd81"
const TimestampString = "2009-11-10 23:00:00 +0000 UTC"

func Generator() bobber.Generator {
	return bobber.Generator{
		Now:  TimeGenerator(TimestampString),
		UUID: UUIDGenerator(UUIDString),
	}
}

func UUIDGenerator(uuidStr string) func() uuid.UUID {
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		panic(err)
	}
	return func() uuid.UUID {
		return id
	}
}

func TimeGenerator(timestamp string) func() time.Time {
	return func() time.Time {
		return ParseTime(timestamp)
	}
}

func ParseTime(timestamp string) time.Time {
	layout := "2006-01-02 15:04:05 +0000 UTC"
	now, err := time.Parse(layout, timestamp)
	if err != nil {
		panic(err)
	}
	return now
}

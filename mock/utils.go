package mock

import (
	"github.com/google/uuid"
	"github.com/tylerdimon/bobber"
	"time"
)

const StaticUUIDValue = "6e300e63-3b0a-470e-b169-f4460e1ccd81"
const StaticTimeValue = "2009-11-10 23:00:00 +0000 UTC"

func Generator() bobber.Generator {
	return bobber.Generator{
		Now:  StaticTime,
		UUID: StaticUUID,
	}
}

func StaticUUID() uuid.UUID {
	id, err := uuid.Parse(StaticUUIDValue)
	if err != nil {
		panic(err)
	}
	return id
}

func StaticTime() time.Time {
	layout := "2006-01-02 15:04:05 +0000 UTC"

	now, err := time.Parse(layout, StaticTimeValue)
	if err != nil {
		panic(err)
	}
	return now
}

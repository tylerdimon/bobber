package mock

import (
	"github.com/google/uuid"
	"github.com/tylerdimon/bobber/sqlite"
	"time"
)

const StaticUUID = "6e300e63-3b0a-470e-b169-f4460e1ccd81"
const StaticTime = "2009-11-10 23:00:00 +0000 UTC m=+0.000000001"

func InitTestDB() *sqlite.DB {
	mockDB := sqlite.DB{
		UUID: staticUUID,
		Now:  staticTime,
	}

	return &mockDB
}

func staticUUID() uuid.UUID {
	id, err := uuid.Parse(StaticUUID)
	if err != nil {
		panic(err)
	}
	return id
}

func staticTime() time.Time {
	layout := "2006-01-02 15:04:05 +0000 UTC m=+0.000000001"

	now, err := time.Parse(layout, StaticTime)
	if err != nil {
		panic(err)
	}
	return now
}

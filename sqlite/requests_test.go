package sqlite

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tylerdimon/bobber"
	"github.com/tylerdimon/bobber/mock"
	"log"
	"testing"
)

func populateRequests(db *DB) {
	tx, err := db.conn.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Commit()

	stmt, err := tx.Prepare(`
		INSERT INTO requests (id, method, url, host, path, timestamp, body, headers) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	request1 := bobber.Request{
		ID:        mock.UUIDString,
		Timestamp: mock.TimestampString,
		Method:    "GET",
		URL:       "/path/one",
		Host:      "google.com",
		Path:      "",
		Body:      "",
	}

	request2 := bobber.Request{
		ID:        "6e300e63-3b0a-470e-b169-f4460e1ccd82",
		Timestamp: "2009-11-10 23:00:01 +0000 UTC",
		Method:    "POST",
		URL:       "/path/two",
		Host:      "example.com",
		Path:      "",
		Body:      "some body text",
	}

	_, err = stmt.Exec(request1.ID, request1.Method, request1.URL, request1.Host, request1.Path, request1.Timestamp, request1.Body, "")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(request2.ID, request2.Method, request2.URL, request2.Host, request2.Path, request2.Timestamp, request2.Body, "")
	if err != nil {
		log.Fatal(err)
	}
}

func TestRequestGetById(t *testing.T) {
	db := initDB()
	defer db.Close()

	populateRequests(db)

	service := &RequestService{
		DB: db,
	}

	tests := []struct {
		name     string
		id       string
		expected bobber.Request
		wantErr  bool
	}{
		{
			name: "Get Request By ID",
			id:   mock.UUIDString,
			expected: bobber.Request{
				ID:        mock.UUIDString,
				Method:    "GET",
				URL:       "/path/one",
				Host:      "google.com",
				Path:      "",
				Timestamp: mock.TimestampString,
				Body:      "",
				Headers:   nil,
			},
			wantErr: false,
		},
		{
			name:    "Request Not Found",
			id:      "non-existent-id",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetById(tt.id)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.expected, *got)
			}
		})
	}
}

func TestRequestGetAll(t *testing.T) {
	db := initDB()
	defer db.Close()

	populateRequests(db)

	service := &RequestService{
		DB: db,
	}

	expected := []*bobber.Request{
		{
			ID:        "6e300e63-3b0a-470e-b169-f4460e1ccd82",
			Method:    "POST",
			URL:       "/path/two",
			Host:      "example.com",
			Path:      "",
			Timestamp: "2009-11-10 23:00:01 +0000 UTC",
			Body:      "some body text",
			Headers:   nil,
		},
		{
			ID:        mock.UUIDString,
			Method:    "GET",
			URL:       "/path/one",
			Host:      "google.com",
			Path:      "",
			Timestamp: mock.TimestampString,
			Body:      "",
			Headers:   nil,
		},
	}

	actual, err := service.GetAll()
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func TestRequestDeleteAll(t *testing.T) {
	db := initDB()
	defer db.Close()

	populateRequests(db)

	var count int
	err := db.conn.Get(&count, "SELECT COUNT(*) FROM requests")
	require.Nil(t, err)
	assert.Equal(t, 2, count)

	service := &RequestService{
		DB: db,
	}

	err = service.DeleteAll()
	assert.Nil(t, err)

	err = db.conn.Get(&count, "SELECT COUNT(*) FROM requests")
	require.Nil(t, err)
	assert.Equal(t, 0, count)
}

//Add

//DeleteById

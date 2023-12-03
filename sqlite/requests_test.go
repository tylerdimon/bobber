package sqlite

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tylerdimon/bobber"
	"github.com/tylerdimon/bobber/mock"
	"testing"
)

func initDB() (*DB, error) {
	db := &DB{
		DSN: ":memory:",
	}
	db.ctx, db.cancel = context.WithCancel(context.Background())
	err := db.Open()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func populateDB(db *DB) error {
	request := bobber.Request{
		ID:        mock.UUIDString,
		Timestamp: mock.TimestampString,
		Method:    "GET",
		URL:       "/path/one",
		Host:      "google.com",
		Path:      "",
		Body:      "",
		Headers:   "",
	}
	_, err := db.conn.NamedExec(`INSERT INTO requests (id, method, url, host, path, timestamp, body, headers)
	                               VALUES (:id, :method, :url, :host, :path, :timestamp, :body, :headers)`, &request)
	if err != nil {
		return err
	}

	request2 := bobber.Request{
		ID:        "6e300e63-3b0a-470e-b169-f4460e1ccd82",
		Timestamp: "2009-11-10 23:00:01 +0000 UTC",
		Method:    "POST",
		URL:       "/path/two",
		Host:      "example.com",
		Path:      "",
		Body:      "some body text",
		Headers:   "",
	}
	_, err = db.conn.NamedExec(`INSERT INTO requests (id, method, url, host, path, timestamp, body, headers)
	                               VALUES (:id, :method, :url, :host, :path, :timestamp, :body, :headers)`, &request2)
	return err
}

func TestGetByID(t *testing.T) {
	db, err := initDB()
	require.Nil(t, err)
	defer db.Close()

	err = populateDB(db)
	require.Nil(t, err)

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
				Headers:   "",
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
			got, err := service.GetByID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && *got != tt.expected {
				t.Errorf("GetByID() got = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	db, err := initDB()
	require.Nil(t, err)
	defer db.Close()

	err = populateDB(db)
	require.Nil(t, err)

	service := &RequestService{
		DB: db,
	}

	expected := []bobber.Request{
		{
			ID:        "6e300e63-3b0a-470e-b169-f4460e1ccd82",
			Method:    "POST",
			URL:       "/path/two",
			Host:      "example.com",
			Path:      "",
			Timestamp: "2009-11-10 23:00:01 +0000 UTC",
			Body:      "some body text",
			Headers:   "",
		},
		{
			ID:        mock.UUIDString,
			Method:    "GET",
			URL:       "/path/one",
			Host:      "google.com",
			Path:      "",
			Timestamp: mock.TimestampString,
			Body:      "",
			Headers:   "",
		},
	}

	actual, err := service.GetAll()
	if err != nil {
		t.Errorf("GetAll() got unexpected error %v", err)
	}
	assert.Equal(t, expected, actual)
}

func TestDeleteAll(t *testing.T) {
	db, err := initDB()
	require.Nil(t, err)
	defer db.Close()

	err = populateDB(db)
	require.Nil(t, err)

	service := &RequestService{
		DB: db,
	}

	err = service.DeleteAll()
	assert.Nil(t, err)

	var count int
	err = db.conn.Get(&count, "SELECT COUNT(*) FROM requests")
	require.Nil(t, err)
	assert.Equal(t, 0, count)
}

//Add

//DeleteByID

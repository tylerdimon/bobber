package sqlite

import (
	"context"
	"github.com/tylerdimon/bobber"
	"github.com/tylerdimon/bobber/mock"
	"reflect"
	"testing"
)

func GetService(t *testing.T) *RequestService {
	db := &DB{
		DSN: ":memory:",
	}
	db.ctx, db.cancel = context.WithCancel(context.Background())

	if err := db.Open(); err != nil {
		t.Fatal(err)
	}

	requestService := &RequestService{
		DB:  db,
		Gen: mock.Generator(),
	}

	return requestService
}

func TestGetByID(t *testing.T) {
	service := GetService(t)
	defer service.DB.Close()

	_, err := service.Add(bobber.Request{
		ID:        mock.StaticUUIDValue,
		Method:    "GET",
		URL:       "/path/one",
		Host:      "google.com",
		Path:      "",
		Timestamp: mock.StaticTimeValue,
		Body:      "",
		Headers:   "",
	})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		id          string
		wantRequest *bobber.Request
		wantErr     bool
	}{
		{
			name: "Get Request By ID",
			id:   mock.StaticUUIDValue,
			wantRequest: &bobber.Request{
				ID:        mock.StaticUUIDValue,
				Method:    "GET",
				URL:       "/path/one",
				Host:      "google.com",
				Path:      "",
				Timestamp: mock.StaticTimeValue,
				Body:      " ",
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
			if !tt.wantErr && reflect.DeepEqual(got, tt.wantRequest) {
				t.Errorf("GetByID() got = %v, want %v", got, tt.wantRequest)
			}
		})
	}
}

//GetByID:
//Test retrieving an existing request.
//Test retrieving a non-existing request.
//Test database errors (using mocking).
//
//GetAll:
//Test retrieving when there are multiple requests in the database.
//Test retrieving when the database is empty.
//Test database errors.
//
//Add:
//Test adding a new request.
//Test adding a request with incomplete or invalid data.
//Test database errors.
//
//DeleteByID:
//Test deleting an existing request.
//Test deleting a non-existing request.
//Test database errors.
//
//DeleteAll:
//Test deleting when there are multiple requests.
//Test deleting when the database is empty.
//Test database errors.

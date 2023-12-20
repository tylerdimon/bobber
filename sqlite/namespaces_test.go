package sqlite

import (
	"github.com/stretchr/testify/assert"
	"github.com/tylerdimon/bobber"
	"github.com/tylerdimon/bobber/mock"
	"log"
	"testing"
)

func populateNamespaces(db *DB) {
	tx, err := db.conn.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Commit()

	stmt, err := tx.Prepare(`
		INSERT INTO namespaces (id, slug, name, created_at) 
		VALUES (?, ?, ?, ?)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	namespace1 := bobber.Namespace{
		ID:        mock.UUIDString,
		CreatedAt: mock.TimestampString,
		Slug:      "first-space",
		Name:      "First Space",
	}

	namespace2 := bobber.Namespace{
		ID:        "6e300e63-3b0a-470e-b169-f4460e1ccd82",
		CreatedAt: "2009-11-10 23:00:01 +0000 UTC",
		Slug:      "second-space",
		Name:      "Second Space",
	}

	_, err = stmt.Exec(namespace1.ID, namespace1.Slug, namespace1.Name, namespace1.CreatedAt)
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(namespace2.ID, namespace2.Slug, namespace2.Name, namespace2.CreatedAt)
	if err != nil {
		log.Fatal(err)
	}
}

func TestNamespaceGetById(t *testing.T) {
	db := initDB()
	defer db.Close()

	populateNamespaces(db)

	service := &NamespaceService{
		DB: db,
	}

	tests := []struct {
		name     string
		id       string
		expected bobber.Namespace
		wantErr  bool
	}{
		{
			name: "Get Namespace By ID",
			id:   mock.UUIDString,
			expected: bobber.Namespace{
				ID:        mock.UUIDString,
				Slug:      "first-space",
				Name:      "First Space",
				CreatedAt: mock.TimestampString,
				UpdatedAt: "",
			},
			wantErr: false,
		},
		{
			name:    "Namespace Not Found",
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

func TestNamespaceGetAll(t *testing.T) {
	db := initDB()
	defer db.Close()

	populateNamespaces(db)

	service := &NamespaceService{
		DB: db,
	}

	expected := []*bobber.Namespace{
		{
			ID:        mock.UUIDString,
			Slug:      "first-space",
			Name:      "First Space",
			CreatedAt: mock.TimestampString,
			UpdatedAt: "",
		},
		{
			ID:        "6e300e63-3b0a-470e-b169-f4460e1ccd82",
			Slug:      "second-space",
			Name:      "Second Space",
			CreatedAt: "2009-11-10 23:00:01 +0000 UTC",
			UpdatedAt: "",
		},
	}

	actual, err := service.GetAll()
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

// TODO

//Add

//DeleteById

//Update

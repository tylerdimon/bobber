package sqlite

import (
	"github.com/stretchr/testify/suite"
	"github.com/tylerdimon/bobber"
	"github.com/tylerdimon/bobber/mock"
	"log"
	"testing"
)

const UUID1 = "6e300e63-3b0a-470e-b169-f4460e1ccd82"
const UUID2 = "6e300e63-3b0a-470e-b169-f4460e1ccd83"
const TS1 = "2009-11-10 23:00:01 +0000 UTC"
const TS2 = "2009-11-10 23:00:02 +0000 UTC"

type NamespaceDbSuite struct {
	db *DB
	suite.Suite
}

func (s *NamespaceDbSuite) BeforeTest(suiteName, testName string) {
	s.db = initDB()
	s.populateNamespaces()
}

func (s *NamespaceDbSuite) AfterTest(suiteName, testName string) {
	s.db.Close()
}

func (s *NamespaceDbSuite) populateNamespaces() {
	tx, err := s.db.conn.Begin()
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
		ID:        UUID1,
		CreatedAt: TS1,
		Slug:      "first-space",
		Name:      "First Space",
	}

	namespace2 := bobber.Namespace{
		ID:        UUID2,
		CreatedAt: TS2,
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

func (s *NamespaceDbSuite) TestGetById() {
	service := &NamespaceService{
		DB: s.db,
	}

	tests := []struct {
		name     string
		id       string
		expected bobber.Namespace
		wantErr  bool
	}{
		{
			name: "Get Namespace By ID",
			id:   UUID1,
			expected: bobber.Namespace{
				ID:        UUID1,
				Slug:      "first-space",
				Name:      "First Space",
				CreatedAt: TS1,
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
		s.Suite.Run(tt.name, func() {
			got, err := service.GetById(tt.id)
			if tt.wantErr {
				s.NotNil(err)
			} else {
				s.Nil(err)
				s.Equal(tt.expected, *got)
			}
		})
	}
}

func (s *NamespaceDbSuite) TestGetAll() {
	service := &NamespaceService{
		DB: s.db,
	}

	expected := []*bobber.Namespace{
		{
			ID:        UUID1,
			Slug:      "first-space",
			Name:      "First Space",
			CreatedAt: TS1,
			UpdatedAt: "",
		},
		{
			ID:        UUID2,
			Slug:      "second-space",
			Name:      "Second Space",
			CreatedAt: TS2,
			UpdatedAt: "",
		},
	}

	actual, err := service.GetAll()
	s.Nil(err)
	s.Equal(expected, actual)
}

func (s *NamespaceDbSuite) TestAdd() {
	var count int
	err := s.db.conn.Get(&count, "SELECT COUNT(*) FROM namespaces")
	s.Require().Nil(err)
	s.Equal(2, count)

	service := &NamespaceService{
		DB:  s.db,
		Gen: mock.Generator(),
	}

	tests := []struct {
		name     string
		toAdd    bobber.Namespace
		expected bobber.Namespace
		wantErr  bool
	}{
		{
			name: "Get Namespace By ID",
			toAdd: bobber.Namespace{
				Slug: "third-space",
				Name: "Third Space",
			},
			expected: bobber.Namespace{
				ID:        mock.UUIDString,
				Slug:      "third-space",
				Name:      "Third Space",
				CreatedAt: mock.TimestampString,
				UpdatedAt: "",
				Endpoints: nil,
			},
			wantErr: false,
		},
		{
			name: "Namespace Not Found",
			toAdd: bobber.Namespace{
				Slug:      "first-space",
				Name:      "First Space",
				UpdatedAt: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Suite.Run(tt.name, func() {
			got, err := service.Add(tt.toAdd)
			if tt.wantErr {
				s.NotNil(err)
			} else {
				s.Nil(err)
				s.Equal(tt.expected, *got)
			}
		})
	}

	err = s.db.conn.Get(&count, "SELECT COUNT(*) FROM namespaces")
	s.Require().Nil(err)
	s.Equal(3, count)
}

func (s *NamespaceDbSuite) TestDeleteById() {

}

func (s *NamespaceDbSuite) TestUpdate() {

}

func TestNamespaceDbSuite(t *testing.T) {
	suite.Run(t, new(NamespaceDbSuite))
}

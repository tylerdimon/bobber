package sqlite

import (
	"github.com/google/uuid"
	"github.com/tylerdimon/bobber"
	"time"
)

type NamespaceService struct {
	DB *DB
}

func (s *NamespaceService) GetByID(id string) (*bobber.Namespace, error) {
	var namespace bobber.Namespace
	err := s.DB.conn.Get(&namespace, "SELECT * FROM namespaces WHERE id = ?", id)
	return &namespace, err
}

func (s *NamespaceService) GetAll() ([]*bobber.Namespace, error) {
	// TODO alphabetical by name ordering
	var spaces []*bobber.Namespace
	err := s.DB.conn.Select(&spaces, "SELECT * FROM namespaces")
	return spaces, err
}

func (s *NamespaceService) Add(namespace bobber.Namespace) (*bobber.Namespace, error) {
	namespace.ID = uuid.New().String()
	namespace.CreatedAt = time.Now().String()

	result, err := s.DB.conn.NamedExec(
		`INSERT INTO namespaces (id, slug, name, created_at) VALUES (:id, :slug, :name, :created_at)`, &namespace)
	if err != nil {
		return nil, err
	}

	_, err = result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &namespace, nil
}

func (s *NamespaceService) Update(namespace bobber.Namespace) (bobber.Namespace, error) {
	namespace.UpdatedAt = time.Now().String()
	_, err := s.DB.conn.NamedExec(`UPDATE namespaces SET slug = :slug, name = :name, updated_at = :updated_at WHERE id = :id`, &namespace)
	return namespace, err
}

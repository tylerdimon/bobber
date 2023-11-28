package sqlite

import (
	"github.com/tylerdimon/bobber"
	"strconv"
)

type NamespaceService struct {
	DB *DB
}

func (s *NamespaceService) GetAll() ([]bobber.Namespace, error) {
	// TODO alphabetical by name ordering
	var spaces []bobber.Namespace
	err := s.DB.conn.Select(&spaces, "SELECT * FROM namespaces")
	return spaces, err
}

func (s *NamespaceService) Add(namespace bobber.Namespace) (*bobber.Namespace, error) {
	result, err := s.DB.conn.NamedExec(`INSERT INTO namespaces (slug, name, timestamp) 
                                    VALUES (:slug, :name, :timestamp)`, &namespace)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	namespace.ID = strconv.Itoa(int(id))
	return &namespace, nil
}

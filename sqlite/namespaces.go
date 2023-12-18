package sqlite

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/tylerdimon/bobber"
	"log"
	"time"
)

type NamespaceService struct {
	DB *DB
}

func (s *NamespaceService) GetById(id string) (*bobber.Namespace, error) {
	var ns bobber.Namespace
	ns.ID = id
	var updatedAt sql.NullString

	query := `
SELECT slug, name, created_at, updated_at FROM namespaces ORDER BY name
`
	err := s.DB.conn.QueryRow(query, id).Scan(&ns.Slug, &ns.Name, &ns.CreatedAt, &updatedAt)
	ns.UpdatedAt = Unwrap(updatedAt)

	endpoints, err := s.getEndpointsByNamespaceId(id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	ns.Endpoints = endpoints

	return &ns, err
}

func (s *NamespaceService) getEndpointsByNamespaceId(id string) ([]bobber.Endpoint, error) {
	query := `SELECT id, method, path, response, created_at, updated_at FROM endpoints WHERE namespace_id = ? ORDER BY path`
	rows, err := s.DB.conn.Query(query, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var endpoints []bobber.Endpoint
	for rows.Next() {
		var e bobber.Endpoint
		var updatedAt sql.NullString

		err := rows.Scan(&e.ID, &e.Method, &e.Path, &e.Response, &e.CreatedAt, &updatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		e.UpdatedAt = Unwrap(updatedAt)

		endpoints = append(endpoints, e)
	}
	return endpoints, nil
}

func (s *NamespaceService) GetAll() ([]*bobber.Namespace, error) {
	query := `
SELECT id, slug, name, created_at, updated_at FROM namespaces ORDER BY name
`
	rows, err := s.DB.conn.Query(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var spaces []*bobber.Namespace
	for rows.Next() {
		var ns bobber.Namespace
		var updatedAt sql.NullString

		err := rows.Scan(&ns.ID, &ns.Slug, &ns.Name, &ns.CreatedAt, &updatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		ns.UpdatedAt = Unwrap(updatedAt)

		spaces = append(spaces, &ns)
	}
	return spaces, nil
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

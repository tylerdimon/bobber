package sqlite

import (
	"github.com/google/uuid"
	"github.com/tylerdimon/bobber"
	"strconv"
	"time"
)

type EndpointService struct {
	DB *DB
}

func (s *EndpointService) GetAll() ([]bobber.Endpoint, error) {
	// TODO alphabetical by name ordering
	var spaces []bobber.Endpoint
	err := s.DB.conn.Select(&spaces, "SELECT * FROM endpoints")
	return spaces, err
}

func (s *EndpointService) Add(endpoint bobber.Endpoint) (*bobber.Endpoint, error) {
	endpoint.ID = uuid.New().String()
	endpoint.CreatedAt = time.Now().String()
	result, err := s.DB.conn.NamedExec(`INSERT INTO endpoints (id, path, response, namespace_id, created_at) 
                                    VALUES (:id, :path, :response, :namespace_id, :created_at)`, &endpoint)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	endpoint.ID = strconv.Itoa(int(id))
	return &endpoint, nil
}

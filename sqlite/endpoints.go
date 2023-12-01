package sqlite

import (
	"github.com/tylerdimon/bobber"
	"strconv"
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
	result, err := s.DB.conn.NamedExec(`INSERT INTO endpoints (slug, name, timestamp) 
                                    VALUES (:slug, :name, :timestamp)`, &endpoint)
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

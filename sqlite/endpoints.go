package sqlite

import (
	"github.com/google/uuid"
	"github.com/tylerdimon/bobber"
	"log"
	"strconv"
	"time"
)

type EndpointService struct {
	DB *DB
}

func (s *EndpointService) GetAll() ([]bobber.Endpoint, error) {
	// TODO alphabetical by name ordering
	var spaces []bobber.Endpoint
	err := s.DB.conn.Select(&spaces, "SELECT * FROM endpoints ORDER BY path asc")
	return spaces, err
}

func (s *EndpointService) Add(endpoint bobber.Endpoint) (*bobber.Endpoint, error) {
	endpoint.ID = uuid.New().String()
	endpoint.CreatedAt = time.Now().String()
	result, err := s.DB.conn.NamedExec(`INSERT INTO endpoints (id, name, method, path, response, namespace_id, created_at) 
                                    VALUES (:id, :name, :method, :path, :response, :namespace_id, :created_at)`, &endpoint)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	endpoint.ID = strconv.Itoa(int(id))
	return &endpoint, nil
}

func (s *EndpointService) DeleteById(id string) error {
	log.Printf("Deleting namespace %s", id)
	_, err := s.DB.conn.NamedExec("DELETE FROM endpoints WHERE id = $1", id)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/tylerdimon/bobber"
	"log"
	"strconv"
	"time"
)

type EndpointService struct {
	DB  *DB
	Gen *bobber.Generator
}

func NewEndpointService(db *DB, gen *bobber.Generator) *EndpointService {
	if db == nil {
		log.Fatal("Endpoint service requires database")

	}
	if gen == nil {
		log.Fatal("Endpoint service requires generator")
	}
	return &EndpointService{
		DB:  db,
		Gen: gen,
	}
}

func (s *EndpointService) GetById(id string) (*bobber.Endpoint, error) {
	var endpoint bobber.Endpoint

	query := `SELECT id, name, method, path, response FROM endpoints WHERE id = ?`

	if err := s.DB.conn.QueryRow(query, id).Scan(&endpoint.ID, &endpoint.Name, &endpoint.Method, &endpoint.Path, &endpoint.Response); err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("NamespaceService GetById %s: Not Found", id)
		}
		return nil, fmt.Errorf("NamespaceService GetById %s: %v", id, err)
	}

	return &endpoint, nil
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

func (s *EndpointService) Update(endpoint bobber.Endpoint) (*bobber.Endpoint, error) {
	endpoint.UpdatedAt = s.Gen.Now().String()
	result, err := s.DB.conn.NamedExec(`UPDATE endpoints SET name = :name, method = :method, path = :path, response = :response, updated_at = :updated_at WHERE id = :id`, &endpoint)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	updates, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if updates == 0 {
		msg := fmt.Sprintf("Endpoint with ID %s does not exist", endpoint.ID)
		log.Print(msg)
		return nil, fmt.Errorf(msg)
	}

	return &endpoint, err
}

func (s *EndpointService) DeleteById(id string) error {
	log.Printf("Deleting endpoint %s", id)
	_, err := s.DB.conn.Exec("DELETE FROM endpoints WHERE id = ?", id)

	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

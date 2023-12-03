package sqlite

import (
	"github.com/tylerdimon/bobber"
	"log"
)

type RequestService struct {
	DB *DB
}

func (s *RequestService) GetByID(id string) (*bobber.Request, error) {
	var req bobber.Request
	err := s.DB.conn.Get(&req, "SELECT * FROM requests WHERE id = ?", id)
	return &req, err
}

func (s *RequestService) GetAll() ([]bobber.Request, error) {
	// TODO add ordering by timestamp descending
	var reqs []bobber.Request
	err := s.DB.conn.Select(&reqs, "SELECT * FROM requests")
	return reqs, err
}

func (s *RequestService) Add(request bobber.Request) (*bobber.Request, error) {
	request.ID = s.DB.UUID().String()
	request.Timestamp = s.DB.Now().String()
	// TODO convert timestamp
	result, err := s.DB.conn.NamedExec(`INSERT INTO requests (id, method, url, host, path, timestamp, body, headers)
	                               VALUES (:id, :method, :url, :host, :path, :timestamp, :body, :headers)`, &request)
	if err != nil {
		log.Printf("Error saving request to database - Request %v : %v", request, err)
		return nil, err
	}

	_, err = result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (s *RequestService) DeleteByID(id string) (*bobber.Request, error) {
	req, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	_, err = s.DB.conn.Exec("DELETE FROM requests WHERE id = ?", id)
	return req, err
}

func (s *RequestService) DeleteAll() error {
	_, err := s.DB.conn.Exec("DELETE FROM requests")
	return err
}

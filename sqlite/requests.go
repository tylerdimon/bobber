package sqlite

import (
	"github.com/tylerdimon/bobber"
	"strconv"
)

type RequestService struct {
	DB *DB
}

func (s *RequestService) GetByID(id string) (bobber.Request, error) {
	var req bobber.Request
	err := s.DB.conn.Get(&req, "SELECT * FROM requests WHERE id = ?", id)
	return req, err
}

func (s *RequestService) GetAll() ([]bobber.Request, error) {
	// TODO add ordering by timestamp descending
	var reqs []bobber.Request
	err := s.DB.conn.Select(&reqs, "SELECT * FROM requests")
	return reqs, err
}

func (s *RequestService) Add(request bobber.Request) (*bobber.Request, error) {
	result, err := s.DB.conn.NamedExec(`INSERT INTO requests (method, url, host, path, timestamp, body, headers) 
                                    VALUES (:method, :url, :host, :path, :timestamp, :body, :headers)`, &request)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	request.ID = strconv.Itoa(int(id))
	return &request, nil
}

func (s *RequestService) Update(request bobber.Request) (bobber.Request, error) {
	_, err := s.DB.conn.NamedExec(`UPDATE requests SET method = :method, url = :url, host = :host, path = :path, 
                              timestamp = :timestamp, body = :body, headers = :headers WHERE id = :id`, &request)
	return request, err
}

func (s *RequestService) DeleteByID(id string) (bobber.Request, error) {
	req, err := s.GetByID(id)
	if err != nil {
		return bobber.Request{}, err
	}

	_, err = s.DB.conn.Exec("DELETE FROM requests WHERE id = ?", id)
	return req, err
}

func (s *RequestService) DeleteAll() error {
	_, err := s.DB.conn.Exec("DELETE FROM requests")
	return err
}

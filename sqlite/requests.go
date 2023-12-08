package sqlite

import (
	"github.com/tylerdimon/bobber"
	"log"
	"strings"
)

type RequestService struct {
	DB  *DB
	Gen bobber.Generator
}

func (s *RequestService) GetByID(id string) (*bobber.Request, error) {
	var req bobber.Request
	err := s.DB.conn.Get(&req, "SELECT * FROM requests WHERE id = ?", id)
	return &req, err
}

func (s *RequestService) GetAll() ([]bobber.RequestDetail, error) {
	var requestDetails []bobber.RequestDetail

	query := `
	SELECT requests.id,
       requests.method,
       requests.url,
       requests.host,
       requests.path,
       requests.timestamp,
       requests.body,
       requests.headers,
       requests.namespace_id,
       namespaces.name "namespace_name",
       requests.endpoint_id,
       endpoints.method_path  "endpoint_path"
FROM requests
         LEFT JOIN namespaces on requests.namespace_id = namespaces.id
         LEFT JOIN endpoints on requests.endpoint_id = endpoints.id
ORDER BY timestamp DESC;`

	err := s.DB.conn.Select(&requestDetails, query)

	return requestDetails, err
}

func (s *RequestService) Add(request bobber.RequestDetail) (*bobber.RequestDetail, error) {
	request.ID = s.Gen.UUID().String()
	request.Timestamp = s.Gen.Now().String()
	result, err := s.DB.conn.NamedExec(`INSERT INTO requests (id, method, url, host, path, timestamp, body, headers, namespace_id, endpoint_id)
	                               VALUES (:id, :method, :url, :host, :path, :timestamp, :body, :headers, :namespace_id, :endpoint_id)`, &request)
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

func (s *RequestService) Match(method string, path string) (namespaceID, endpointID, response string) {
	log.Printf("Matching request with method %s and path %s", method, path)
	parts := strings.SplitN(path, "/", 4)

	// parts[0] is throwaway value, will always be empty string with leading slash
	// parts[1] is always going to be requests
	// parts[2] might be a namespace, but only if it matches
	possibleNamespace := parts[2]
	namespaceID = s.matchNamespace(possibleNamespace)
	if namespaceID == "" {
		log.Printf("Request with method %s and path %s did not match a namespace", method, path)
		// request is unmatched
		return
	}

	// parts[2] is the rest of the path without leading slash
	endpointID, response = s.matchEndpoint(namespaceID, method, "/"+parts[3])
	log.Printf("Request with method %s and path %s matched the following namespace %s and endpoint %s", method, path, namespaceID, endpointID)
	return
}

func (s *RequestService) matchNamespace(namespace string) (namespaceID string) {
	log.Printf("Looking for match for namespace %s", namespace)

	var id string
	err := s.DB.conn.Get(&id, " SELECT id FROM namespaces WHERE slug = ?", namespace)
	if err != nil {
		log.Print(err)
	}
	return id
}

func (s *RequestService) matchEndpoint(namespaceID, method, path string) (endpointID string, response string) {
	var endpoint bobber.Endpoint
	err := s.DB.conn.Get(&endpoint, " SELECT id, response FROM endpoints WHERE namespace_id = ? AND method_path = ?", namespaceID, method+" "+path)
	if err != nil {
		log.Print(err)
	}
	return endpoint.ID, endpoint.Response
}

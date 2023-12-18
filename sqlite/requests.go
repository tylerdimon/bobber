package sqlite

import (
	"encoding/json"
	"github.com/tylerdimon/bobber"
	"log"
	"strings"
)

type RequestService struct {
	DB  *DB
	Gen bobber.Generator
}

func (s *RequestService) GetById(id string) (*bobber.Request, error) {
	var req bobber.Request
	err := s.DB.conn.Get(&req, "SELECT * FROM requests WHERE id = ?", id)
	return &req, err
}

func (s *RequestService) GetAll() ([]bobber.Request, error) {
	query := `
SELECT r.id, r.method, r.url, r.host, r.path, r.timestamp, r.body,
	   r.headers, r.namespace_id, r.endpoint_id, n.name
FROM requests r
LEFT JOIN namespaces n on r.namespace_id = n.id
ORDER BY timestamp DESC;`

	rows, err := s.DB.conn.Query(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var requests []bobber.Request

	for rows.Next() {
		var req bobber.Request
		var headersJSON string

		err := rows.Scan(&req.ID, &req.Method, &req.URL, &req.Host, &req.Path, &req.Timestamp, &req.Body,
			&headersJSON, &req.NamespaceID, &req.NamespaceName, &req.EndpointID)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		if headersJSON != "" {
			err = json.Unmarshal([]byte(headersJSON), &req.Headers)
			if err != nil {
				log.Println(err)
				return nil, err
			}
		}

		requests = append(requests, req)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return requests, nil
}

func (s *RequestService) Add(request bobber.Request) (*bobber.Request, error) {
	request.ID = s.Gen.UUID().String()
	request.Timestamp = s.Gen.Now().String()

	headersJSON, err := json.Marshal(request.Headers)
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := s.DB.conn.Prepare(`INSERT INTO requests (id, method, url, host, path, timestamp, body, headers, namespace_id, endpoint_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Printf("Error preparing statement - Request %v : %v", request, err)
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(request.ID, request.Method, request.URL, request.Host, request.Path, request.Timestamp, request.Body, string(headersJSON), request.NamespaceID, request.EndpointID)
	if err != nil {
		log.Printf("Error saving request to database - Request %v : %v", request, err)
		return nil, err
	}

	return &request, nil
}

func (s *RequestService) DeleteById(id string) (*bobber.Request, error) {
	req, err := s.GetById(id)
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

// Match takes in a request method and path in this format /requests/{namespace}/{endpoint}
// and returns a matching namespace, endpoint, and response if exists
func (s *RequestService) Match(method string, path string) (namespaceID, endpointID, response *string) {
	log.Printf("Matching request with method %s and path %s", method, path)
	parts := strings.SplitN(path, "/", 4)

	namespaceID = s.matchNamespace(parts[2])
	if namespaceID == nil {
		log.Printf("Request with method %s and path %s did not match a namespace", method, path)
		return
	}

	endpointID, response = s.matchEndpoint(*namespaceID, method, "/"+parts[3])
	log.Printf("Request with method %s and path %s matched the following namespace %s and endpoint %s", method, path, namespaceID, endpointID)
	return
}

func (s *RequestService) matchNamespace(namespace string) (namespaceID *string) {
	log.Printf("Looking for match for namespace %s", namespace)

	var id string
	err := s.DB.conn.Get(&id, " SELECT id FROM namespaces WHERE slug = ?", namespace)
	if err != nil {
		log.Print(err)
		return nil
	}
	return &id
}

func (s *RequestService) matchEndpoint(namespaceID, method, path string) (endpointID, response *string) {
	var endpoint bobber.Endpoint
	err := s.DB.conn.Get(&endpoint, " SELECT id, response FROM endpoints WHERE namespace_id = ? AND method = ? AND PATH = ?", namespaceID, method, path)
	if err != nil {
		// TODO make sure this catches when there is no match
		log.Print(err)
		err := s.DB.conn.Get(&endpoint, " SELECT id, response FROM endpoints WHERE namespace_id = ? AND method = ? AND PATH = '*'", namespaceID, method)
		if err != nil {
			log.Print(err)
			return nil, nil
		}
	}
	return &endpoint.ID, &endpoint.Response
}

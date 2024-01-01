package sqlite

import (
	"database/sql"
	"encoding/json"
	"github.com/tylerdimon/bobber"
	"log"
	"strings"
)

type RequestService struct {
	DB  *DB
	Gen bobber.Generator
}

type Scannable interface {
	Scan(dest ...any) error
}

func scan(rows Scannable) (*bobber.Request, error) {
	var r bobber.Request
	var headersJSON string
	var namespaceId sql.NullString
	var namespaceName sql.NullString
	var endpointID sql.NullString
	var ts string

	err := rows.Scan(&r.ID, &r.Method, &r.URL, &r.Host, &r.Path, &ts, &r.Body,
		&headersJSON, &namespaceId, &endpointID, &namespaceName)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if headersJSON != "" {
		err = json.Unmarshal([]byte(headersJSON), &r.Headers)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}

	r.NamespaceID = Unwrap(namespaceId)
	r.NamespaceName = Unwrap(namespaceName)
	r.EndpointID = Unwrap(endpointID)

	timestamp, err := ParseTime(ts)
	if err != nil {
		log.Printf("Error parsing timestamp for request: %s", err)
	}
	r.Timestamp = timestamp

	return &r, nil
}

func (s *RequestService) GetById(id string) (*bobber.Request, error) {
	query := `
SELECT r.id, r.method, r.url, r.host, r.path, r.timestamp, r.body,
	   r.headers, r.namespace_id, r.endpoint_id, n.name 
FROM requests r 
LEFT JOIN namespaces n on r.namespace_id = n.id 
WHERE r.id = ?;`

	r, err := scan(s.DB.conn.QueryRow(query, id))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return r, nil
}

func (s *RequestService) GetAll() ([]*bobber.Request, error) {
	query := `
SELECT r.id, r.method, r.url, r.host, r.path, r.timestamp, r.body,
	   r.headers, r.namespace_id, r.endpoint_id, n.name
FROM requests r
LEFT JOIN namespaces n on r.namespace_id = n.id
ORDER BY r.timestamp DESC;`

	rows, err := s.DB.conn.Query(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var requests []*bobber.Request
	for rows.Next() {
		r, err := scan(rows)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		requests = append(requests, r)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return requests, nil
}

func (s *RequestService) Add(request bobber.Request) (*bobber.Request, error) {
	request.ID = s.Gen.UUID().String()
	request.Timestamp = s.Gen.Now()

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

	var namespaceId sql.NullString
	if request.NamespaceID == "" {
		namespaceId.Valid = false
	} else {
		namespaceId.Valid = true
		namespaceId.String = request.NamespaceID
	}

	var endpointId sql.NullString
	if request.NamespaceID == "" {
		endpointId.Valid = false
	} else {
		endpointId.Valid = true
		endpointId.String = request.EndpointID
	}

	_, err = stmt.Exec(request.ID, request.Method, request.URL, request.Host, request.Path, request.Timestamp, request.Body, string(headersJSON), namespaceId, endpointId)
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
func (s *RequestService) Match(method string, path string) (namespaceID, endpointID, response string) {
	log.Printf("Matching request with method %s and path %s", method, path)
	parts := strings.SplitN(path, "/", 4)

	namespaceID = s.matchNamespace(parts[2])
	if namespaceID == "" {
		log.Printf("Request with method %s and path %s did not match a namespace", method, path)
		return
	}

	endpointID, response = s.matchEndpoint(namespaceID, method, "/"+parts[3])
	log.Printf("Request with method %s and path %s matched the following namespace %s and endpoint %s", method, path, namespaceID, endpointID)
	return
}

func (s *RequestService) matchNamespace(slug string) (namespaceID string) {
	log.Printf("Looking for namespace match for slug %s", slug)

	var id string
	err := s.DB.conn.Get(&id, " SELECT id FROM namespaces WHERE slug = ?", slug)
	if err != nil {
		log.Print(err)
		return ""
	}
	return id
}

func (s *RequestService) matchEndpoint(namespaceID, method, path string) (endpointID, response string) {
	var endpoint bobber.Endpoint
	err := s.DB.conn.Get(&endpoint, " SELECT id, response FROM endpoints WHERE namespace_id = ? AND method = ? AND PATH = ?", namespaceID, method, path)
	if err != nil {
		// TODO make sure this catches when there is no match
		log.Print(err)
		err := s.DB.conn.Get(&endpoint, " SELECT id, response FROM endpoints WHERE namespace_id = ? AND method = ? AND PATH = '*'", namespaceID, method)
		if err != nil {
			log.Print(err)
			return "", ""
		}
	}
	return endpoint.ID, endpoint.Response
}

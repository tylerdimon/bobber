package bobber

import (
	"database/sql"
	"fmt"
	"time"
)

type Request struct {
	ID          string `json:"id" db:"id"` // uuid.UUID format
	Method      string `json:"method" db:"method"`
	URL         string `json:"url" db:"url"`
	Host        string `json:"host" db:"host"`
	Path        string `json:"path" db:"path"`
	Timestamp   string `json:"timestamp" db:"timestamp"` // time.Time default format
	Body        string `json:"body" db:"body"`
	Headers     string `json:"headers" db:"headers"`
	NamespaceID string `json:"namespaceId" db:"namespace_id"`
	EndpointID  string `json:"endpointId" db:"endpoint_id"`
}

type RequestDetail struct {
	ID            string         `json:"id" db:"id"` // uuid.UUID format
	Method        string         `json:"method" db:"method"`
	URL           string         `json:"url" db:"url"`
	Host          string         `json:"host" db:"host"`
	Path          string         `json:"path" db:"path"`
	Timestamp     string         `json:"timestamp" db:"timestamp"` // time.Time default format
	Body          string         `json:"body" db:"body"`
	Headers       string         `json:"headers" db:"headers"`
	NamespaceID   sql.NullString `json:"namespaceId" db:"namespace_id"`
	NamespaceName sql.NullString `json:"namespaceName" db:"namespace_name"`
	EndpointID    sql.NullString `json:"endpointId" db:"endpoint_id"`
	EndpointPath  sql.NullString `json:"endpointPath" db:"endpoint_path"`
}

func (r Request) String() string {
	return fmt.Sprintf("Timestamp: %v\nMethod: %v\nURL: %v\nHost: %v\nPath: %v\nHeaders: %v\nBody: %v",
		time.Now().Format(time.RFC3339), r.Method, r.URL, r.Host, r.Path, r.Headers, r.Body)
}

type RequestService interface {
	GetByID(id string) (*Request, error)
	GetAll() ([]RequestDetail, error)
	Add(request RequestDetail) (*RequestDetail, error)
	DeleteByID(id string) (*Request, error)
	DeleteAll() error
	Match(method string, path string) (namespaceID, endpointID, response string)
}

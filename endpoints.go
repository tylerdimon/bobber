package bobber

import "database/sql"

type Endpoint struct {
	ID          string         `db:"id"`
	Method      string         `db:"method"`
	Path        string         `db:"path"`
	Response    string         `db:"response"`
	CreatedAt   string         `db:"created_at"`
	UpdatedAt   sql.NullString `db:"updated_at"`
	NamespaceID string         `db:"namespace_id"`
}

type EndpointService interface {
	GetAll() ([]Endpoint, error)
	GetAllByNamespace(namespaceID string) ([]Endpoint, error)
	Add(request Endpoint) (*Endpoint, error)
}

package bobber

type Endpoint struct {
	ID          string `db:"id"`          // uuid.UUID format
	MethodPath  string `db:"method_path"` // ex. GET /request/to/api
	Response    string `db:"response"`
	CreatedAt   string `db:"created_at"` // time.Time default format
	UpdatedAt   string `db:"updated_at"` // time.Time default format
	NamespaceID string `db:"namespace_id"`
}

type EndpointService interface {
	GetAll() ([]Endpoint, error)
	Add(request Endpoint) (*Endpoint, error)
}

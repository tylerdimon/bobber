package bobber

type Endpoint struct {
	ID          string `json:"id" db:"id"` // uuid.UUID format
	Path        string `json:"path" db:"path"`
	Response    string `json:"response" db:"response"`
	CreatedAt   string `json:"createdAt" db:"created_at"` // time.Time default format
	UpdatedAt   string `json:"updatedAt" db:"updated_at"` // time.Time default format
	NamespaceID string `json:"namespaceId" db:"namespace_id"`
}

type EndpointService interface {
	GetAll() ([]Endpoint, error)
	Add(request Endpoint) (*Endpoint, error)
}

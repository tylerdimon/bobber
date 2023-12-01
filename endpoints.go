package bobber

type Endpoint struct {
	ID          string `json:"id" db:"id"`
	Path        string `json:"path" db:"path"`
	Response    string `json:"response" db:"response"`
	CreatedAt   string `json:"createdAt" db:"created_at"`
	UpdatedAt   string `json:"updatedAt" db:"updated_at"`
	NamespaceID string `json:"namespaceId" db:"namespace_id"`
}

type EndpointService interface {
	GetAll() ([]Endpoint, error)
	Add(request Endpoint) (*Endpoint, error)
}

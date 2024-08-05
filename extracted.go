package bobber

// An Extracted is a value captured by an Extractor
// for a specific request to an endpoint
type Extracted struct {
	ID         string `db:"id"`
	Path       string `db:"path"`
	Value      string `db:"value"`
	EndpointID string `db:"endpoint_id"`
	RequestID  string `db:"request_id"`
	CreatedAt  string `db:"created_at"`
	UpdatedAt  string `db:"updated_at"`
}

type ExtractedService interface {
	GetAllByEndpoint(endpointId string) ([]*Extracted, error)
	AddAll(extractors []Extracted) (*Extracted, error)
}

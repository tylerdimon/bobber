package bobber

// An Extractor is used to capture values from a request
// A captured value gets stored as an Extracted type
type Extractor struct {
	ID         string `db:"id"`
	Path       string `db:"path"`
	EndpointID string `db:"endpoint_id"`
	CreatedAt  string `db:"created_at"`
	UpdatedAt  string `db:"updated_at"`
}

type ExtractorService interface {
	GetAllByEndpoint(endpointId string) ([]*Extractor, error)
	Add(request Extractor) (*Extractor, error)
	Update(request Extractor) (*Extractor, error)
	DeleteById(id string) (*Extractor, error)
	DeleteAll() error
}

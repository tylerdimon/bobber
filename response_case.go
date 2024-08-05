package bobber

// ResponseCase finds a certain value for a request
// send a different response
// Check different cases in priority ordered, first match is used
type ResponseCase struct {
	ID          string `db:"id"`
	MatchValue  string `db:"match_value"`
	Response    string `db:"response"`
	Priority    int    `db:"priority"`
	ExtractorId string `db:"extractor_id"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
}

type ResponseCaseService interface {
	// Get a match if one exists
	GetMatch(value string) (string, error)
	// Get all response cases sorted by priority
	GetAllForEndpoint(endpointId string) ([]*ResponseCase, error)
	Add(request ResponseCase) (*ResponseCase, error)
	Update(request ResponseCase) (*ResponseCase, error)
	DeleteById(id string) (*ResponseCase, error)
	DeleteAll() error
}

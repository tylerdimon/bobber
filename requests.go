package bobber

type Request struct {
	ID             string   `db:"id"`
	Method         string   `db:"method"`
	URL            string   `db:"url"`
	Host           string   `db:"host"`
	Path           string   `db:"path"`
	Timestamp      string   `db:"timestamp"`
	Body           string   `db:"body"`
	Headers        []Header `db:"headers"`
	NamespaceID    string   `db:"namespace_id"`
	NamespaceName  string   `db:"namespace_name"`
	EndpointID     string   `db:"endpoint_id"`
	EndpointMethod string   `db:"endpoint_method"`
	EndpointPath   string   `db:"endpoint_path"`
}

type Header struct {
	Name  string
	Value string
}

type RequestService interface {
	GetById(id string) (*Request, error)
	GetAll() ([]Request, error)
	Add(request Request) (*Request, error)
	DeleteById(id string) (*Request, error)
	DeleteAll() error
	Match(method string, path string) (namespaceID, endpointID, response string)
}

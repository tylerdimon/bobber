package bobber

type Namespace struct {
	ID        string
	Slug      string
	Name      string
	Timestamp string
}

type NamespaceService interface {
	GetByID(id string) (*Namespace, error)
	GetAll() ([]Namespace, error)
	Add(request Namespace) (*Namespace, error)
	//Update(request Namespace) (Namespace, error)
	//DeleteByID(id string) (Namespace, error)
	//DeleteAll() error
}

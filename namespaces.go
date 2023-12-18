package bobber

type Namespace struct {
	ID        string `db:"id"`
	Slug      string `db:"slug"`
	Name      string `db:"name"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
	Endpoints []Endpoint
}

type NamespaceService interface {
	GetById(id string) (*Namespace, error)
	GetAll() ([]*Namespace, error)
	Add(request Namespace) (*Namespace, error)
	Update(request Namespace) (Namespace, error)
}

package bobber

type Namespace struct {
	ID        string `json:"id" db:"id"` // uuid.UUID format
	Slug      string `json:"slug" db:"slug"`
	Name      string `json:"name" db:"name"`
	CreatedAt string `json:"createdAt" db:"created_at"` // time.Time default format
	UpdatedAt string `json:"updatedAt" db:"updated_at"` // time.Time default format
}

type NamespaceService interface {
	GetByID(id string) (*Namespace, error)
	GetAll() ([]*Namespace, error)
	Add(request Namespace) (*Namespace, error)
	Update(request Namespace) (Namespace, error)
	//DeleteByID(id string) (Namespace, error)
	//DeleteByID(id string) (Namespace, error)
	//DeleteAll() error
}

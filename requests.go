package bobber

import (
	"fmt"
	"time"
)

type Request struct {
	ID        string
	Method    string
	URL       string
	Host      string
	Path      string
	Timestamp string
	Body      string
	Headers   string
}

func (r Request) String() string {
	return fmt.Sprintf("Timestamp: %v\nMethod: %v\nURL: %v\nHost: %v\nPath: %v\nHeaders: %v\nBody: %v",
		time.Now().Format(time.RFC3339), r.Method, r.URL, r.Host, r.Path, r.Headers, r.Body)
}

type RequestService interface {
	GetByID(id string) (Request, error)
	GetAll() ([]Request, error)
	Add(request Request) (*Request, error)
	Update(request Request) (Request, error)
	DeleteByID(id string) (Request, error)
	DeleteAll() error
}

package mock

import (
	"github.com/tylerdimon/bobber"
)

type RequestService struct {
	Requests []bobber.Request
	Gen      bobber.Generator

	GetByIDCalled    int
	GetAllCalled     int
	AddCalled        int
	UpdateCalled     int
	DeleteByIDCalled int
	DeleteAllCalled  int
}

func (s *RequestService) Add(request bobber.Request) (*bobber.Request, error) {
	request.ID = s.Gen.UUID().String()
	request.Timestamp = s.Gen.Now().String()

	s.AddCalled = s.AddCalled + 1
	s.Requests = append(s.Requests, request)
	return nil, nil
}

func (s *RequestService) GetByID(id string) (*bobber.Request, error) {
	s.GetByIDCalled = s.GetByIDCalled + 1
	return nil, nil
}

func (s *RequestService) GetAll() ([]bobber.Request, error) {
	s.GetAllCalled = s.GetAllCalled + 1
	return s.Requests, nil
}

func (s *RequestService) Update(request bobber.Request) (bobber.Request, error) {
	s.UpdateCalled = s.UpdateCalled + 1
	return request, nil
}

func (s *RequestService) DeleteByID(id string) (*bobber.Request, error) {
	s.DeleteByIDCalled = s.DeleteByIDCalled + 1
	return nil, nil
}

func (s *RequestService) DeleteAll() error {
	s.DeleteAllCalled = s.DeleteAllCalled + 1
	return nil
}

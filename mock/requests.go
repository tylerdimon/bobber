package mock

import "github.com/tylerdimon/bobber"

type RequestService struct {
	Requests []bobber.Request

	GetByIDCalled    int
	GetAllCalled     int
	AddCalled        int
	UpdateCalled     int
	DeleteByIDCalled int
	DeleteAllCalled  int
}

func (m *RequestService) Add(request bobber.Request) (*bobber.Request, error) {
	m.AddCalled = m.AddCalled + 1
	m.Requests = append(m.Requests, request)
	return nil, nil
}

func (m *RequestService) GetByID(id string) (bobber.Request, error) {
	m.GetByIDCalled = m.GetByIDCalled + 1
	return bobber.Request{}, nil
}

func (m *RequestService) GetAll() ([]bobber.Request, error) {
	m.GetAllCalled = m.GetAllCalled + 1
	return m.Requests, nil
}

func (m *RequestService) Update(request bobber.Request) (bobber.Request, error) {
	m.UpdateCalled = m.UpdateCalled + 1
	return request, nil
}

func (m *RequestService) DeleteByID(id string) (bobber.Request, error) {
	m.DeleteByIDCalled = m.DeleteByIDCalled + 1
	return bobber.Request{}, nil
}

func (m *RequestService) DeleteAll() error {
	m.DeleteAllCalled = m.DeleteAllCalled + 1
	return nil
}

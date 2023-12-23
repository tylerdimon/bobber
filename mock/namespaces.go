package mock

import (
	"fmt"
	"github.com/stretchr/testify/mock"
	"github.com/tylerdimon/bobber"
)

type NamespaceService struct {
	mock.Mock
}

func (s *NamespaceService) GetById(id string) (*bobber.Namespace, error) {
	args := s.Called(id)
	fmt.Println("ARGS")
	fmt.Println(args)
	return nil, nil
}

func (s *NamespaceService) GetAll() ([]*bobber.Namespace, error) {
	args := s.Called()
	fmt.Println("ARGS")
	fmt.Println(args)
	return nil, nil
}

func (s *NamespaceService) Add(request bobber.Namespace) (*bobber.Namespace, error) {
	args := s.Called(request)
	fmt.Println("ARGS")
	fmt.Println(args)
	return nil, nil
}

func (s *NamespaceService) Update(request bobber.Namespace) (*bobber.Namespace, error) {
	args := s.Called(request)
	fmt.Println("ARGS")
	fmt.Println(args)
	return nil, nil
}

func (s *NamespaceService) DeleteById(id string) error {
	args := s.Called(id)
	fmt.Println("ARGS")
	fmt.Println(args)
	return nil
}

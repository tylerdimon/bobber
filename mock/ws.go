package mock

import "github.com/tylerdimon/bobber"

type WebsocketService struct {
	broadcast chan *bobber.Request
}

func (s *WebsocketService) Init() {
	s.broadcast = make(chan *bobber.Request)
}

func (s *WebsocketService) Broadcast() chan *bobber.Request {
	return s.broadcast
}

func (s *WebsocketService) HandleMessages() {

}

func (s *WebsocketService) AddClient(client bobber.Client) {

}
func (s *WebsocketService) DeleteClient(client bobber.Client) {

}

package mock

import "github.com/tylerdimon/bobber"

type WebsocketService struct {
	broadcast chan string
}

func (s *WebsocketService) Init() {
	s.broadcast = make(chan string)
}

func (s *WebsocketService) Broadcast() chan string {
	return s.broadcast
}

func (s *WebsocketService) HandleMessages() {

}

func (s *WebsocketService) AddClient(client bobber.Client) {

}
func (s *WebsocketService) DeleteClient(client bobber.Client) {

}

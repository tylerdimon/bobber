package ws

import (
	"github.com/gorilla/websocket"
	"github.com/tylerdimon/bobber"
	"github.com/tylerdimon/bobber/static"
	"log"
	"net/http"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WebsocketService struct {
	clients   map[bobber.Client]bool
	broadcast chan *bobber.RequestDetail
}

func (s *WebsocketService) Init() {
	s.clients = make(map[bobber.Client]bool)
	s.broadcast = make(chan *bobber.RequestDetail)
}

func (s *WebsocketService) HandleMessages() {
	for {
		request := <-s.broadcast
		log.Printf("Received a websocket message: %v", request)
		msg, err := static.GetRequestHTML(request)
		if err != nil {
			log.Printf("error getting requst HTML to send over websocket: %v", err)
			return
		}

		for client := range s.clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(s.clients, client)
			}
		}
	}
}

func (s *WebsocketService) Broadcast() chan *bobber.RequestDetail {
	return s.broadcast
}

func (s *WebsocketService) AddClient(client bobber.Client) {
	s.clients[client] = true

}
func (s *WebsocketService) DeleteClient(client bobber.Client) {
	delete(s.clients, client)
}

package ws

import (
	"github.com/gorilla/websocket"
	"github.com/tylerdimon/bobber"
	"log"
	"net/http"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WebsocketService struct {
	clients   map[bobber.Client]bool
	broadcast chan string
}

func (s *WebsocketService) Init() {
	s.clients = make(map[bobber.Client]bool)
	s.broadcast = make(chan string)
}

func (s *WebsocketService) HandleMessages() {
	for {
		msg := <-s.broadcast
		log.Printf("Received a websocket message: %v", msg)
		for client := range s.clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(s.clients, client)
			}
		}
	}
}

func (s *WebsocketService) Broadcast() chan string {
	return s.broadcast
}

func (s *WebsocketService) AddClient(client bobber.Client) {
	s.clients[client] = true

}
func (s *WebsocketService) DeleteClient(client bobber.Client) {
	delete(s.clients, client)
}

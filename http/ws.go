package http

import (
	"github.com/gorilla/mux"
	"github.com/tylerdimon/bobber"
	"github.com/tylerdimon/bobber/ws"
	"log"
	"net/http"
)

type WebsocketHandler struct {
	WebsocketService bobber.WebsocketService
}

func (h *WebsocketHandler) RegisterWebsocketRoutes(r *mux.Router) {
	r.HandleFunc("/ws", h.HandleConnections)

}

func (h *WebsocketHandler) HandleConnections(w http.ResponseWriter, r *http.Request) {
	socket, err := ws.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer socket.Close()

	h.WebsocketService.AddClient(socket)

	for {
		_, _, err := socket.ReadMessage()
		if err != nil {
			h.WebsocketService.DeleteClient(socket)
			break
		}
	}
}

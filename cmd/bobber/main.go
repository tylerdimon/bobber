package main

import (
	http2 "github.com/tylerdimon/bobber/http"
	"github.com/tylerdimon/bobber/sqlite"
	"github.com/tylerdimon/bobber/ws"
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	db := sqlite.NewDB("bobber.sqlite")
	if err := db.Open(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	requestService := &sqlite.RequestService{}
	requestService.DB = db

	websocketService := ws.WebsocketService{}
	websocketService.Init()
	go websocketService.HandleMessages()

	requestHandler := &http2.RequestHandler{}
	requestHandler.Service = requestService
	requestHandler.WebsocketService = &websocketService

	websocketHandler := http2.WebsocketHandler{}
	websocketHandler.WebsocketService = &websocketService

	// WebSockets
	http.HandleFunc("/ws", websocketHandler.HandleConnections)

	// API
	http.HandleFunc("/api/requests/delete", requestHandler.DeleteAllRequestsHandler)
	http.HandleFunc("/api/requests/all", requestHandler.GetAllRequests)

	// UI
	http.HandleFunc("/requests/", requestHandler.AddRequestHandler)

	// Static file server
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	log.Println("Listening on :8000...")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"github.com/joho/godotenv"
	"github.com/tylerdimon/bobber"
	"github.com/tylerdimon/bobber/http"
	"github.com/tylerdimon/bobber/sqlite"
	"github.com/tylerdimon/bobber/ws"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.SetFlags(log.LstdFlags | log.Llongfile)

	db := sqlite.NewDB("bobber.sqlite?parseTime=true")
	if err := db.Open(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	generator := bobber.NewGenerator()

	requestService := &sqlite.RequestService{}
	requestService.DB = db
	requestService.Gen = generator

	namespaceService := &sqlite.NamespaceService{}
	namespaceService.DB = db
	namespaceService.Gen = generator

	endpointService := sqlite.NewEndpointService(db, &generator)

	websocketService := &ws.WebsocketService{}
	websocketService.Init()
	go websocketService.HandleMessages()

	server := &http.Server{}
	server.WebsocketService = websocketService
	server.RequestService = requestService
	server.NamespaceService = namespaceService
	server.EndpointService = endpointService
	server.Init()
	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}
}

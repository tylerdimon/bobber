package main

import (
	"github.com/joho/godotenv"
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

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	db := sqlite.NewDB("bobber.sqlite")
	if err := db.Open(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	requestService := &sqlite.RequestService{}
	requestService.DB = db

	namespaceService := &sqlite.NamespaceService{}
	namespaceService.DB = db

	websocketService := &ws.WebsocketService{}
	websocketService.Init()
	go websocketService.HandleMessages()

	server := &http.Server{}
	server.WebsocketService = websocketService
	server.RequestService = requestService
	server.NamespaceService = namespaceService
	server.Init()
	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}
}

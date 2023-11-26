package http

import (
	"github.com/gorilla/mux"
	"github.com/tylerdimon/bobber"
	"github.com/tylerdimon/bobber/static"
	"log"
	"net/http"
)

type Server struct {
	server *http.Server
	router *mux.Router

	WebsocketService bobber.WebsocketService
	RequestService   bobber.RequestService
}

func (s *Server) Init() {
	s.server = &http.Server{}
	s.router = mux.NewRouter()

	//Our router is wrapped by another function handler to perform some
	//middleware-like tasks that cannot be performed by actual middleware.
	//This includes changing route paths for JSON endpoints & overridding methods.
	s.server.Handler = http.HandlerFunc(s.router.ServeHTTP)

	// Setup error handling routes.
	//s.router.NotFoundHandler = http.HandlerFunc(s.handleNotFound)

	s.router.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.FS(static.Assets))))

	static.ParseHTML()

	requestHandler := &RequestHandler{}
	requestHandler.Service = s.RequestService
	requestHandler.WebsocketService = s.WebsocketService
	requestHandler.RegisterRequestRoutes(s.router)

	websocketHandler := WebsocketHandler{}
	websocketHandler.WebsocketService = s.WebsocketService
	websocketHandler.RegisterWebsocketRoutes(s.router)
}

func (s *Server) Run() error {
	log.Println("Listening on :8000...")
	s.server.Addr = ":8000"
	return s.server.ListenAndServe()
}

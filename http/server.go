package http

import (
	"github.com/gorilla/mux"
	"github.com/tylerdimon/bobber"
	"github.com/tylerdimon/bobber/static"
	"log"
	"net/http"
	"os"
)

type Server struct {
	server *http.Server
	router *mux.Router

	WebsocketService bobber.WebsocketService
	RequestService   bobber.RequestService
	NamespaceService bobber.NamespaceService
	EndpointService  bobber.EndpointService
}

func (s *Server) serveHTTP(w http.ResponseWriter, r *http.Request) {
	// Override method for forms passing "_method" value.
	if r.Method == http.MethodPost {
		switch v := r.PostFormValue("_method"); v {
		case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
			r.Method = v
		}
	}

	s.router.ServeHTTP(w, r)
}

func (s *Server) Init() {
	s.server = &http.Server{}
	s.router = mux.NewRouter()

	//Our router is wrapped by another function handler to perform some
	//middleware-like tasks that cannot be performed by actual middleware.
	//This includes changing route paths for JSON endpoints & overridding methods.
	s.server.Handler = http.HandlerFunc(s.serveHTTP)

	// Setup error handling routes.
	//s.router.NotFoundHandler = http.HandlerFunc(s.handleNotFound)

	staticFilesDebugMode := os.Getenv("DEBUG")
	if staticFilesDebugMode == "True" {
		log.Println("Serving static files from directory...")
		s.router.PathPrefix("/static/assets").
			Handler(http.StripPrefix("/static/assets", http.FileServer(http.Dir("static/assets"))))
	} else {
		log.Println("Serving embedded static files...")
		s.router.PathPrefix("/static/").
			Handler(http.StripPrefix("/static/", http.FileServer(http.FS(static.Assets))))
	}

	static.ParseHTML()

	configHandler := &NamespaceHandler{}
	configHandler.NamespaceService = s.NamespaceService
	configHandler.RegisterNamespaceRoutes(s.router)

	endpointHandler := &EndpointHandler{}
	endpointHandler.EndpointService = s.EndpointService
	endpointHandler.RegisterEndpointRoutes(s.router)

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

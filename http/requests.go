package http

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tylerdimon/bobber"
	"github.com/tylerdimon/bobber/static"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type RequestHandler struct {
	Service          bobber.RequestService
	WebsocketService bobber.WebsocketService
}

func (h *RequestHandler) RegisterRequestRoutes(r *mux.Router) {
	r.HandleFunc("/api/requests/delete", h.DeleteAllRequestsHandler)
	r.HandleFunc("/api/requests/all", h.GetAllRequests)
	r.PathPrefix("/requests/").HandlerFunc(h.RecordRequestHandler)
	r.HandleFunc("/", h.RequestIndexHandler)
}

func (h *RequestHandler) RecordRequestHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
	}
	body := string(bodyBytes)

	var headers []string
	for name, values := range r.Header {
		for _, value := range values {
			headers = append(headers, fmt.Sprintf("%v: %v", name, value))
		}
	}

	request := bobber.Request{
		Method:    r.Method,
		URL:       r.URL.String(),
		Path:      r.URL.Path,
		Host:      r.Host,
		Timestamp: time.Now().Format(time.RFC3339),
		Body:      body,
		Headers:   strings.Join(headers, ", "),
	}

	if _, err := h.Service.Add(request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.WebsocketService.Broadcast() <- &request

	w.Write([]byte("Request received"))
}

func (h *RequestHandler) GetAllRequests(w http.ResponseWriter, r *http.Request) {
	requests, err := h.Service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var strings []string
	for _, req := range requests {
		strings = append(strings, req.String())
	}

	jsonData, err := json.Marshal(strings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (h *RequestHandler) RequestIndexHandler(w http.ResponseWriter, r *http.Request) {
	requests, err := h.Service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageData := struct {
		Title string
		Data  any
	}{
		Title: "Requests",
		Data:  requests,
	}

	err = static.IndexTemplate.Execute(w, pageData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *RequestHandler) DeleteAllRequestsHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.Service.DeleteAll(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "All requests cleared")
}
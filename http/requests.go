package http

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tylerdimon/bobber"
	"github.com/tylerdimon/bobber/static"
	"io"
	"log"
	"net/http"
)

type RequestHandler struct {
	RequestService   bobber.RequestService
	WebsocketService bobber.WebsocketService
}

func (h *RequestHandler) RegisterRequestRoutes(r *mux.Router) {
	r.HandleFunc("/request/{id}", h.deleteHandler).Methods("DELETE")
	r.HandleFunc("/request/{id}", h.detailHandler).Methods("GET")

	r.HandleFunc("/request", h.DeleteAllRequestsHandler).Methods("DELETE")
	r.HandleFunc("/request", h.RequestIndexHandler).Methods("GET")

	r.PathPrefix("/requests/").HandlerFunc(h.RecordRequestHandler)

	r.HandleFunc("/", h.RequestIndexHandler).Methods("GET")
}

func (h *RequestHandler) RecordRequestHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
	}
	body := string(bodyBytes)

	var headers []bobber.Header
	for name, values := range r.Header {
		for _, v := range values {
			headers = append(headers, bobber.Header{
				Name:  name,
				Value: v,
			})
		}
	}

	request := bobber.Request{
		Method:  r.Method,
		Path:    r.URL.Path,
		Host:    r.Host,
		Body:    body,
		Headers: headers,
	}

	namespaceID, endpointID, response := h.RequestService.Match(request.Method, request.Path)
	request.NamespaceID = namespaceID
	request.EndpointID = endpointID
	if response == "" {
		request.Response = "Request received"
	} else {
		request.Response = response
	}

	savedRequest, err := h.RequestService.Add(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.WebsocketService.Broadcast(savedRequest)

	w.Write([]byte(request.Response))
}

func (h *RequestHandler) RequestIndexHandler(w http.ResponseWriter, r *http.Request) {
	requests, err := h.RequestService.GetAll()
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
	if err := h.RequestService.DeleteAll(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("All requests cleared")
	http.Redirect(w, r, "", http.StatusSeeOther)
}

func (h *RequestHandler) detailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	req, err := h.RequestService.GetById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	title := fmt.Sprintf("Request %s %s", req.Method, req.Path)

	pageData := struct {
		Title   string
		Request *bobber.Request
	}{
		Title:   title,
		Request: req,
	}

	err = static.RequestDetailTemplate.Execute(w, pageData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *RequestHandler) deleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if _, err := h.RequestService.DeleteById(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Request %s deleted", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

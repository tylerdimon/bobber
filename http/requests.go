package http

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tylerdimon/bobber"
	"github.com/tylerdimon/bobber/static"
	"io"
	"log"
	"net/http"
	"strings"
)

type RequestHandler struct {
	Service          bobber.RequestService
	WebsocketService bobber.WebsocketService
}

func (h *RequestHandler) RegisterRequestRoutes(r *mux.Router) {
	r.HandleFunc("/api/requests/delete", h.DeleteAllRequestsHandler)
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

	request := bobber.RequestDetail{
		Method:  r.Method,
		URL:     r.URL.String(),
		Path:    r.URL.Path,
		Host:    r.Host,
		Body:    body,
		Headers: strings.Join(headers, "\n"),
	}

	namespaceID, endpointID, response := h.Service.Match(request.Method, request.Path)

	if namespaceID == "" {
		request.NamespaceID = sql.NullString{
			String: "",
			Valid:  false,
		}
	} else {
		request.NamespaceID = sql.NullString{
			String: namespaceID,
			Valid:  true,
		}
	}
	if endpointID == "" {
		request.EndpointID = sql.NullString{
			String: "",
			Valid:  false,
		}
	} else {
		request.EndpointID = sql.NullString{
			String: endpointID,
			Valid:  true,
		}
	}

	savedRequest, err := h.Service.Add(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.WebsocketService.Broadcast() <- savedRequest

	if response == "" {
		w.Write([]byte("Request received"))
	} else {
		w.Write([]byte(response))
	}
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

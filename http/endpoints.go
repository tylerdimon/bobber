package http

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tylerdimon/bobber"
	"github.com/tylerdimon/bobber/static"
	"log"
	"net/http"
	"time"
)

type EndpointHandler struct {
	EndpointService bobber.EndpointService
}

func (h *EndpointHandler) RegisterEndpointRoutes(r *mux.Router) {
	r.HandleFunc("/config/namespace/{id}/endpoint", h.endpointDetailHandler)
	r.HandleFunc("/config/endpoint/{id}", h.endpointDetailHandler)
}

func (h *EndpointHandler) endpointDetailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespaceID := vars["id"]

	if r.Method == "GET" {
		pageData := struct {
			Title       string
			NamespaceID string
		}{
			Title:       "Add Endpoint",
			NamespaceID: namespaceID,
		}

		err := static.EndpointAddTemplate.Execute(w, pageData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing the form", http.StatusInternalServerError)
		return
	}

	endpoint := bobber.Endpoint{
		Name:        r.FormValue("name"),
		Method:      r.FormValue("method"),
		Path:        r.FormValue("path"),
		Response:    r.FormValue("response"),
		CreatedAt:   time.Now().String(),
		NamespaceID: namespaceID,
	}

	added, err := h.EndpointService.Add(endpoint)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error adding to database: %v", err), http.StatusInternalServerError)
		return
	}

	// Print the namespace info. In a real application, you might save it to a database.
	log.Printf("Endpoint Added: %+v", added)

	http.Redirect(w, r, fmt.Sprintf("/namespace/%v", namespaceID), http.StatusSeeOther)
}

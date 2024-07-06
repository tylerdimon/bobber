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
	r.HandleFunc("/namespace/{nId}/endpoint/{eId}/delete", h.deleteById).Methods("GET")
	r.HandleFunc("/namespace/{nId}/endpoint/{eId}", h.detail).Methods("GET")
	r.HandleFunc("/namespace/{nId}/endpoint/{eId}", h.detail).Methods("PUT")
	r.HandleFunc("/namespace/{nId}/endpoint", h.detail).Methods("GET")
	r.HandleFunc("/namespace/{nId}/endpoint", h.detail).Methods("POST")
}

func (h *EndpointHandler) detail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespaceId := vars["nId"]
	endpointId := vars["eId"]
	log.Printf("Request to endpoint detail page for endpoint %s", namespaceId, endpointId)

	if r.Method == "GET" {
		var endpoint *bobber.Endpoint
		var err error
		if endpointId != "" {
			endpoint, err = h.EndpointService.GetById(endpointId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
		pageData := struct {
			Title       string
			NamespaceID string
			Endpoint    *bobber.Endpoint
		}{
			Title:       "Add Endpoint",
			NamespaceID: namespaceId,
			Endpoint:    endpoint,
		}

		err = static.EndpointDetailTemplate.Execute(w, pageData)
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

	if r.Method == "POST" {
		endpoint := bobber.Endpoint{
			Name:        r.FormValue("name"),
			Method:      r.FormValue("method"),
			Path:        r.FormValue("path"),
			Response:    r.FormValue("response"),
			CreatedAt:   time.Now().String(),
			NamespaceID: namespaceId,
		}

		added, err := h.EndpointService.Add(endpoint)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error adding to database: %v", err), http.StatusInternalServerError)
			return
		}

		// Print the namespace info. In a real application, you might save it to a database.
		log.Printf("Endpoint Added: %+v", added)
	} else if r.Method == "PUT" && endpointId != "" {
		endpoint := bobber.Endpoint{
			ID:        endpointId,
			Name:      r.FormValue("name"),
			Method:    r.FormValue("method"),
			Path:      r.FormValue("path"),
			Response:  r.FormValue("response"),
			UpdatedAt: time.Now().String(),
		}

		updated, err := h.EndpointService.Update(endpoint)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error updating endpoint in database: %v", err), http.StatusInternalServerError)
			return
		}

		// Print the namespace info. In a real application, you might save it to a database.
		log.Printf("Endpoint Updated: %+v", updated)
	}

	http.Redirect(w, r, fmt.Sprintf("/namespace/%v", namespaceId), http.StatusSeeOther)
}

func (h *EndpointHandler) deleteById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespaceId := vars["nId"]
	endpointId := vars["eId"]
	log.Printf("Handling delete request for endpoint %s", endpointId)

	err := h.EndpointService.DeleteById(endpointId)
	if err != nil {
		http.Error(w, "Error parsing the form", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/namespace/%v", namespaceId), http.StatusSeeOther)
}

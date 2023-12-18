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

type ConfigHandler struct {
	NamespaceService bobber.NamespaceService
	EndpointService  bobber.EndpointService
}

func (h *ConfigHandler) RegisterConfigRoutes(r *mux.Router) {
	r.HandleFunc("/config", h.configIndexHandler)
	r.HandleFunc("/config/namespace", h.namespaceDetailHandler)
	r.HandleFunc("/config/namespace/{id}", h.namespaceDetailHandler)
	r.HandleFunc("/config/namespace/{id}/endpoint", h.addEndpointHandler)
}

func (h *ConfigHandler) configIndexHandler(w http.ResponseWriter, r *http.Request) {
	namespaces, err := h.NamespaceService.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageData := struct {
		Title string
		Data  any
	}{
		Title: "Config",
		Data:  namespaces,
	}

	err = static.ConfigTemplate.Execute(w, pageData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *ConfigHandler) namespaceDetailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if r.Method == "GET" {
		h.serveNamespaceDetail(w, id)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing the form", http.StatusInternalServerError)
		return
	}

	namespace := bobber.Namespace{
		Slug: r.FormValue("slug"),
		Name: r.FormValue("name"),
	}

	if r.Method == "PUT" {
		namespace.ID = id
		updated, err := h.NamespaceService.Update(namespace)
		if err != nil {
			http.Error(w, "Error updating database", http.StatusInternalServerError)
			return
		}

		log.Printf("Namespace Updated: %+v", updated)

	} else if r.Method == "POST" {
		added, err := h.NamespaceService.Add(namespace)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error adding to database: %v", err), http.StatusInternalServerError)
			return
		}

		log.Printf("Namespace Added: %+v", added)
	}

	http.Redirect(w, r, "/config", http.StatusSeeOther)

}

func (h *ConfigHandler) serveNamespaceDetail(w http.ResponseWriter, id string) {
	var title string
	var namespace *bobber.Namespace
	var endpoints []bobber.Endpoint
	var err error

	if id == "" {
		title = "Add Namespace"
	} else {
		title = "Edit Namespace"
		namespace, err = h.NamespaceService.GetByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		endpoints, err = h.EndpointService.GetAllByNamespace(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	pageData := struct {
		Title     string
		Namespace *bobber.Namespace
		Endpoints []bobber.Endpoint
	}{
		Title:     title,
		Namespace: namespace,
		Endpoints: endpoints,
	}

	err = static.NamespaceAddTemplate.Execute(w, pageData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *ConfigHandler) addEndpointHandler(w http.ResponseWriter, r *http.Request) {
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

	http.Redirect(w, r, fmt.Sprintf("/config/namespace/%v", namespaceID), http.StatusSeeOther)
}

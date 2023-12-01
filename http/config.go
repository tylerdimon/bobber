package http

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tylerdimon/bobber"
	"github.com/tylerdimon/bobber/static"
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
		var title string
		var namespace *bobber.Namespace
		var err error

		if id == "" {
			title = "Add Namespace"
		} else {
			title = "Edit Namespace"
			namespace, err = h.NamespaceService.GetByID(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

		pageData := struct {
			Title     string
			Namespace *bobber.Namespace
		}{
			Title:     title,
			Namespace: namespace,
		}

		err = static.NamespaceAddTemplate.Execute(w, pageData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if r.Method == "PUT" {
		// do update
		fmt.Println("UPDATING NAMESPACE NOT YET IMPLEMENTED SORRY BOUT THAT")
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing the form", http.StatusInternalServerError)
		return
	}

	namespace := bobber.Namespace{
		ID:        r.FormValue("id"),
		Slug:      r.FormValue("slug"),
		Name:      r.FormValue("name"),
		Timestamp: r.FormValue("timestamp"),
	}

	added, err := h.NamespaceService.Add(namespace)
	if err != nil {
		http.Error(w, "Error adding to database", http.StatusInternalServerError)
		return
	}

	// Print the namespace info. In a real application, you might save it to a database.
	fmt.Printf("Namespace Added: %+v", added)

	http.Redirect(w, r, "/config", http.StatusSeeOther)

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
		Path:        r.FormValue("path"),
		Response:    r.FormValue("response"),
		CreatedAt:   time.Now().Format(time.RFC3339),
		NamespaceID: namespaceID,
	}

	added, err := h.EndpointService.Add(endpoint)
	if err != nil {
		http.Error(w, "Error adding to database", http.StatusInternalServerError)
		return
	}

	// Print the namespace info. In a real application, you might save it to a database.
	fmt.Printf("Endpoint Added: %+v", added)

	http.Redirect(w, r, fmt.Sprintf("/config/namespace/%v", namespaceID), http.StatusSeeOther)
}

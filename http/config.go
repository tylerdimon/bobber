package http

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tylerdimon/bobber"
	"github.com/tylerdimon/bobber/static"
	"net/http"
)

type ConfigHandler struct {
	NamespaceService bobber.NamespaceService
}

func (h *ConfigHandler) RegisterConfigRoutes(r *mux.Router) {
	r.HandleFunc("/config", h.configIndexHandler)
	r.HandleFunc("/config/namespace", h.addNamespaceHandler)
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

func (h *ConfigHandler) addNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	pageData := struct {
		Title string
	}{
		Title: "Add Namespace",
	}

	if r.Method == "GET" {
		err := static.NamespaceAddTemplate.Execute(w, pageData)
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

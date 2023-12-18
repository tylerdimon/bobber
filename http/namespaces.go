package http

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tylerdimon/bobber"
	"github.com/tylerdimon/bobber/static"
	"log"
	"net/http"
)

type NamespaceHandler struct {
	NamespaceService bobber.NamespaceService
}

func (h *NamespaceHandler) RegisterNamespaceRoutes(r *mux.Router) {
	r.HandleFunc("/config", h.configIndexHandler)
	r.HandleFunc("/config/namespace", h.namespaceDetailHandler)
	r.HandleFunc("/config/namespace/{id}", h.namespaceDetailHandler)
}

func (h *NamespaceHandler) configIndexHandler(w http.ResponseWriter, r *http.Request) {
	namespaces, err := h.NamespaceService.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageData := struct {
		Title string
		Data  any
	}{
		Title: "Namespace",
		Data:  namespaces,
	}

	err = static.ConfigTemplate.Execute(w, pageData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *NamespaceHandler) namespaceDetailHandler(w http.ResponseWriter, r *http.Request) {
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

func (h *NamespaceHandler) serveNamespaceDetail(w http.ResponseWriter, id string) {
	var title string
	var namespace *bobber.Namespace
	var err error

	if id == "" {
		title = "Add Namespace"
	} else {
		title = "Edit Namespace"
		namespace, err = h.NamespaceService.GetById(id)
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
		Endpoints: namespace.Endpoints,
	}

	err = static.NamespaceAddTemplate.Execute(w, pageData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

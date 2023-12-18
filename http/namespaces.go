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
	r.HandleFunc("/namespace/{id}/delete", h.deleteByIdHandler).Methods("GET")
	r.HandleFunc("/namespace/{id}", h.detailHandler).Methods("GET")
	r.HandleFunc("/namespace/{id}", h.updateHandler).Methods("PUT")
	r.HandleFunc("/namespace", h.detailHandler).Methods("GET")
	r.HandleFunc("/namespace", h.AddHandler).Methods("POST")
	r.HandleFunc("/config", h.getAllHandler)
}

func (h *NamespaceHandler) getAllHandler(w http.ResponseWriter, r *http.Request) {
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

func (h *NamespaceHandler) detailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var title string
	var namespace *bobber.Namespace
	endpoints := make([]bobber.Endpoint, 0)
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
		endpoints = namespace.Endpoints
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

func (h *NamespaceHandler) AddHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing the form", http.StatusInternalServerError)
		return
	}

	namespace := bobber.Namespace{
		Slug: r.FormValue("slug"),
		Name: r.FormValue("name"),
	}

	added, err := h.NamespaceService.Add(namespace)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error adding to database: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Namespace Added: %+v", added)
	http.Redirect(w, r, "/config", http.StatusSeeOther)

}

func (h *NamespaceHandler) updateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing the form", http.StatusInternalServerError)
		return
	}

	namespace := bobber.Namespace{
		Slug: r.FormValue("slug"),
		Name: r.FormValue("name"),
	}

	namespace.ID = id
	updated, err := h.NamespaceService.Update(namespace)
	if err != nil {
		http.Error(w, "Error updating database", http.StatusInternalServerError)
		return
	}

	log.Printf("Namespace Updated: %+v", updated)

	log.Printf("Namespace Added: %+v", updated)
	http.Redirect(w, r, "/config", http.StatusSeeOther)

}

func (h *NamespaceHandler) deleteByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	log.Printf("Handling delete request for Namespace %s", id)

	err := h.NamespaceService.DeleteById(id)
	if err != nil {
		http.Error(w, "Error parsing the form", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/config", http.StatusSeeOther)
}

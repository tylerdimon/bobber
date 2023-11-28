package http

import (
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

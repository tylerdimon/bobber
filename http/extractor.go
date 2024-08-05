package http

import (
	"github.com/gorilla/mux"
	"github.com/tylerdimon/bobber"
	"github.com/tylerdimon/bobber/static"
	"log"
	"net/http"
)

type ExtractorHandler struct {
	ExtractorService bobber.ExtractorService
}

func (h *ExtractorHandler) Register(r *mux.Router) {
	r.HandleFunc("/extractor", h.detail).Methods("GET")
}

func (h *ExtractorHandler) detail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	extractorId := vars["eId"]
	log.Printf("Request to endpoint detail page for endpoint %s", extractorId)

	err := static.ExtractorDetailTemplate.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return
}

package router

import (
	"net/http"
	"webhook/api/handler"

	"github.com/gorilla/mux"
)

func RouteBuilder(h handler.IHandler) *mux.Router {

	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1/he").Subrouter()
	api.HandleFunc("/{type}/submit", h.HandleSubmitForm).Methods(http.MethodPost)

	return r

}

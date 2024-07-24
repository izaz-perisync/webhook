package handler

import (
	"encoding/json"
	"net/http"
	"webhook/dto.go/request"
	"webhook/service"

	"github.com/gorilla/mux"
)

type IHandler interface {
	HandleSubmitForm(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	s service.IService
}

func New(s service.IService) IHandler {

	return &Handler{
		s: s,
	}
}

func (h *Handler) HandleSubmitForm(w http.ResponseWriter, r *http.Request) {

	b := service.SubmitForm{}
	period := mux.Vars(r)["type"]

	if err := r.ParseMultipartForm(32 << 20); err != nil {

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("err in file")
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	if err := request.BindBody(r, &b); err != nil {

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	ctx := r.Context()
	if err := h.s.SubmitForm(ctx, b, file, handler.Filename, period); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("added")
	return
}

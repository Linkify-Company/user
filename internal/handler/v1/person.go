package v1

import (
	"context"
	"encoding/json"
	auth "github.com/Linkify-Company/auth-client"
	"github.com/Linkify-Company/common_utils/errify"
	"github.com/Linkify-Company/common_utils/pointer"
	"github.com/Linkify-Company/common_utils/response"
	"github.com/gorilla/mux"
	"net/http"
	"user/internal/domain"
	hr "user/internal/handler"
)

func initPerson(h *handler, router *mux.Router, mwf ...mux.MiddlewareFunc) {
	person := router.PathPrefix("/person").Subrouter()

	person.HandleFunc("", h.AddMyself).Methods(http.MethodPost)
	person.HandleFunc("", h.GetMyself).Methods(http.MethodGet)

	router.Use(mwf...)
}

func (h *handler) AddMyself(w http.ResponseWriter, r *http.Request) {
	var req domain.Person
	e := json.NewDecoder(r.Body).Decode(&req)
	if e != nil {
		response.Error(w, errify.NewBadRequestError(e.Error(), hr.ValidationError, "AddPerson/NewDecoder"), h.log)
		return
	}
	user, ok := auth.GetAuthData(r)
	if !ok {
		response.Error(w, errify.NewInternalServerError("Auth data in context is empty", "AddPerson/GetAuthData"), h.log)
		return
	}
	req.UserID = pointer.P(user.ID)
	req.Email = pointer.P(user.Email)

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.ContextTimeout)
	defer cancel()

	id, err := h.service.AddPerson(ctx, &req)
	if err != nil {
		response.Error(w, err.JoinLoc("AddPerson"), h.log)
		return
	}
	response.Ok(w, response.NewSend(id, "Created Person successfully", http.StatusCreated), h.log)
}

func (h *handler) GetMyself(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.GetAuthData(r)
	if !ok {
		response.Error(w, errify.NewInternalServerError("Auth data in context is empty", "GetPersonByUserID/GetAuthData"), h.log)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.ContextTimeout)
	defer cancel()

	id, err := h.service.GetPersonByUserID(ctx, user.ID)
	if err != nil {
		response.Error(w, err.JoinLoc("GetPersonByUserID"), h.log)
		return
	}
	response.Ok(w, response.NewSend(id, "Get Person by user id successfully", http.StatusOK), h.log)
}

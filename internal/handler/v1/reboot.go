package v1

import (
	"context"
	"github.com/Linkify-Company/common_utils/response"
	"github.com/gorilla/mux"
	"net/http"
)

func initReboot(h *handler, router *mux.Router, mwf ...mux.MiddlewareFunc) {
	reboot := router.PathPrefix("/reboot").Subrouter()

	reboot.HandleFunc("/system", h.RebootSystem).Methods(http.MethodPost)

	reboot.Use(mwf...)
}

func (h *handler) RebootSystem(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.ContextTimeout)
	defer cancel()

	err := h.service.RebootSystem(ctx)
	if err != nil {
		response.Error(w, err, h.log)
		return
	}
	response.Ok(w, response.NewSend("", "reboot successfully", http.StatusOK), h.log)
}

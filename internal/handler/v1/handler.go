package v1

import (
	"fmt"
	auth "github.com/Linkify-Company/auth-client"
	"github.com/Linkify-Company/common_utils/errify"
	"github.com/Linkify-Company/common_utils/logger"
	"github.com/gorilla/mux"
	"net/http"
	"user/internal/config"
	hr "user/internal/handler"
	"user/internal/service"
)

type handler struct {
	cfg     *config.HandlerConfig
	log     logger.Logger
	service *service.Service
}

func NewHandler(
	cfg *config.HandlerConfig,
	log logger.Logger,
	service *service.Service,
) hr.Handler {
	return &handler{
		cfg:     cfg,
		log:     log,
		service: service,
	}
}

func (h *handler) Init(router *mux.Router) {
	h.log.Infof("Initialization handler V1")

	var authMiddleware = auth.NewMiddleware(h.log, h.service.AuthService)

	version := router.PathPrefix("/v1").Subrouter()
	version.Use(h.panicMiddleware)

	initPerson(h, version, authMiddleware.AuthHandler)
	initReboot(h, version, authMiddleware.AuthHandler, authMiddleware.AuthFuncWithRoles(auth.RoleAdmin))
}

func (h *handler) panicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				h.log.Error(errify.NewInternalServerError(fmt.Sprint(err), r.RequestURI).SetDetails("There was a panic in the router under version No. 1"))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

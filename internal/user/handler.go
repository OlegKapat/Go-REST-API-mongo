package user

import (
	"github.com/OlegKapat/Rest-api-mongo/internal/apperror"
	"github.com/OlegKapat/Rest-api-mongo/pkg/logging"
	"net/http"

	"github.com/OlegKapat/Rest-api-mongo/internal/handlers"
	"github.com/julienschmidt/httprouter"
)

var _ handlers.Handler = &handler{}

type handler struct {
	logger *logging.Logger
}

const (
	usersURL = "/users"
	userURL  = "/users/:uuid"
)

func NewHandler(logger *logging.Logger) handlers.Handler {
	return &handler{
		logger: logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, usersURL, apperror.Middleware(h.GetList))
	router.HandlerFunc(http.MethodGet, userURL, apperror.Middleware(h.GetUserByUUID))
	router.HandlerFunc(http.MethodPost, usersURL, apperror.Middleware(h.CreateUser))
	router.HandlerFunc(http.MethodPut, userURL, apperror.Middleware(h.UpdateUser))
	router.HandlerFunc(http.MethodPatch, userURL, apperror.Middleware(h.PartiallyUpdateUser))
	router.HandlerFunc(http.MethodDelete, userURL, apperror.Middleware(h.DeleteUser))
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("This is list of users"))
	return nil
}

func (h *handler) GetUserByUUID(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("This is get of user by id"))
	return nil
}
func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("This is create  of users"))
	return nil
}
func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("This is update of users"))
	return nil
}
func (h *handler) PartiallyUpdateUser(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("This is partitially update of users"))
	return nil
}
func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	w.Write([]byte("This is delete of users"))
	return nil
}

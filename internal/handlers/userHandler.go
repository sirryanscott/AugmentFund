package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/AugmentFund/internal/data"
	"github.com/AugmentFund/internal/services"
	"github.com/go-chi/chi"
)

//TODO use json marshaller
//TODO use context

// UserHandler is a handler for user routes
type UserHandler struct {
	UserService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users, err := h.UserService.GetUsers(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	user, err := h.UserService.GetUser(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userBytes, err := ioReaderToBytes(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user data.User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := h.UserService.CreateUser(ctx, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func ioReaderToBytes(reader io.ReadCloser) ([]byte, error) {
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

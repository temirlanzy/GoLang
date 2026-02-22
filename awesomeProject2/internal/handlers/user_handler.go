package handlers

import (
	"encoding/json"
	"net/http"
	"practice3/internal/usecase"
	"practice3/pkg/modules"
	"strconv"
)

type UserHandler struct {
	usecase *usecase.UserUsecase
}

func NewUserHandler(u *usecase.UserUsecase) *UserHandler {
	return &UserHandler{usecase: u}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, _ := h.usecase.GetUsers()
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	user, err := h.usecase.GetUserByID(id)
	if err != nil {
		http.Error(w, "user not found", 404)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user modules.User
	json.NewDecoder(r.Body).Decode(&user)
	id, _ := h.usecase.CreateUser(user)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	var user modules.User
	json.NewDecoder(r.Body).Decode(&user)
	err := h.usecase.UpdateUser(id, user)
	if err != nil {
		http.Error(w, "user not found", 404)
		return
	}
	json.NewEncoder(w).Encode(map[string]bool{"updated": true})
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	err := h.usecase.DeleteUser(id)
	if err != nil {
		http.Error(w, "user not found", 404)
		return
	}
	json.NewEncoder(w).Encode(map[string]bool{"deleted": true})
}

package handlers

import (
	"encoding/json"
	"net/http"
	"organizer/internal/service"
	"organizer/pkg/logger"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

type Request struct {
	Password string `json:"password"`
	Name     string `json:"name"`
}

type Response struct {
	Success bool `json:"success"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Debug("Не удалось распарсить json :(")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err := h.service.Register(req.Password, req.Name)
	if err != nil {
		logger.Warn("Неудачная попытка регистрации")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Debug("Удачная регистрация")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{Success: true})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Debug("Не удалось распарсить json :(")
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	_, err := h.service.Login(req.Name, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Здесь позже сгенерить JWT

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Success: true})
}

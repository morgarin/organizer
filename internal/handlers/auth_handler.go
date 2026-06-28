package handlers

import (
	"encoding/json"
	"net/http"
	"organizer/internal/models"
	"organizer/pkg/logger"
)

type AuthHandler struct {
	service models.AuthServiceInterface
}

func NewAuthHandler(service models.AuthServiceInterface) *AuthHandler {
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	err := h.service.Register(req.Password, req.Name)
	if err != nil {
		logger.Warn("Неудачная попытка регистрации")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	_, err := h.service.Login(req.Name, req.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Success: true})
}

package handler

import (
    "encoding/json"
    "net/http"
    "log/slog"
    "goproject/internal/jwt"
    "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
    "goproject/internal/domain"
)

type UserHandler struct {
	uc domain.UserUsecase
	log *slog.Logger
}

func NewUserHandler(uc domain.UserUsecase, log *slog.Logger) *UserHandler {
	return &UserHandler{uc: uc, log: log}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Name string `json:"name"`
        Email string `json:"email"`
        Phone string `json:"phone"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    user, err := h.uc.CreateUser(r.Context(), req.Name, req.Email, req.Phone, req.Password)
    if err != nil {
        h.log.Error("RegisterUser failed", slog.Any("error", err))
        http.Error(w, err.Error(), http.StatusConflict)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    user, err := h.uc.GetUserByEmail(r.Context(), req.Email)
    if err != nil || user == nil {
        http.Error(w, "invalid login", http.StatusUnauthorized)
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        http.Error(w, "invalid login or password", http.StatusUnauthorized)
        return
    }

    token, err := jwt.GenerateToken(user.ID, user.Email)
    if err != nil {
        h.log.Error("GenerateToken failed", slog.Any("error", err))
        http.Error(w, "failed to generate token", http.StatusInternalServerError)
        return
    }

    sessionID := uuid.New().String()
    if err := h.uc.CreateSession(r.Context(), sessionID, user.ID); err != nil {
        h.log.Error("CreateSession failed", slog.Any("error", err))
        http.Error(w, "failed to create session", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"token": token})
}

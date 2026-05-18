package handler

import (
    "encoding/json"
    "net/http"
    "goproject/internal/domain"
)

type UserHandler struct {
    uc domain.UserUsecase
}

func NewUserHandler(uc domain.UserUsecase) *UserHandler {
    return &UserHandler{uc: uc}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Name string `json:"name"`
        Email string `json:"email"`
        Phone string `json:"phone"`
        Password string `json:"password"`
    }
    json.NewDecoder(r.Body).Decode(&req)
    user, err := h.uc.Register(req.Name, req.Email, req.Phone, req.Password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(user)
}

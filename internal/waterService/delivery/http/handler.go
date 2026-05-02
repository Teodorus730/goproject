package http

import (
    "net/http"
	waterService "goproject/internal/waterService"
)

type handler struct {
    uc waterService.UseCase
}

func New(uc waterService.UseCase) *handler {
    return &handler{uc: uc}
}

func (h *handler) Test() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        msg, err := h.uc.GetTest(r.Context())
        if err != nil {
            http.Error(w, "internal error", http.StatusInternalServerError)
			return
        }

		w.WriteHeader(http.StatusOK)
        w.Write([]byte(msg))
    }
}

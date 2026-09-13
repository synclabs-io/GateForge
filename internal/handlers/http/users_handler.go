package http_handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/synclabs-io/GateForge/internal/config"
	domain_errors "github.com/synclabs-io/GateForge/internal/domain/errors"
	dto2 "github.com/synclabs-io/GateForge/internal/service/dto"
)

type UsersService interface {
	Register(ctx context.Context, input dto2.RegisterRequestDTO) (dto2.RegisterResponseDTO, error)
}

type UsersHandler struct {
	svc UsersService
	cfg *config.Config
}

func NewUsersHandler(svc UsersService, cfg *config.Config) *UsersHandler {
	return &UsersHandler{svc: svc, cfg: cfg}
}

func (h *UsersHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto2.RegisterRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.svc.Register(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, domain_errors.ErrUserAlreadyExists):
			SendError(w, http.StatusConflict, err.Error())
		case errors.Is(err, domain_errors.ErrUsernameEmpty),
			errors.Is(err, domain_errors.ErrUsernameInvalid), errors.Is(err, domain_errors.ErrUsernameTooShort), errors.Is(err, domain_errors.ErrUsernameTooLong),
			errors.Is(err, domain_errors.ErrPasswordEmpty), errors.Is(err, domain_errors.ErrPasswordInvalid), errors.Is(err, domain_errors.ErrPasswordTooShort),
			errors.Is(err, domain_errors.ErrPasswordTooLong):
			SendError(w, http.StatusBadRequest, err.Error())
		default:
			SendError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	SetCookies(w, resp, h.cfg.JWT.AccessExpiry, h.cfg.JWT.RefreshExpiry, h.cfg.App.Env == "production")

	SendJSON(w, http.StatusCreated, resp)
}

package http_handlers

import (
	"encoding/json"
	"net/http"
	"time"

	dto2 "github.com/synclabs-io/GateForge/internal/service/dto"
)

func SendJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func SendError(w http.ResponseWriter, status int, message string) {
	SendJSON(w, status, map[string]string{"error": message})
}

func SetCookies(w http.ResponseWriter, response dto2.RegisterResponseDTO, AccessTokenTTL time.Duration, RefreshTokenTTL time.Duration, isSecure bool) {
	accessCookie := http.Cookie{
		Name:     "access_token",
		Value:    response.Access,
		Path:     "/",
		MaxAge:   int(AccessTokenTTL.Seconds()),
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
	}
	refreshCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    response.Refresh,
		Path:     "/",
		MaxAge:   int(RefreshTokenTTL.Seconds()),
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &accessCookie)
	http.SetCookie(w, &refreshCookie)
}

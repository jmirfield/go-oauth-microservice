package server

import (
	"encoding/json"
	"net/http"
	"oauth/internal/errors"
	"oauth/internal/models"
	"strings"
	"time"
)

type response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func writeJSON(w http.ResponseWriter, data any, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func (a *app) registerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		client, err := a.m.RegisterClient(ctx)
		if err != nil {
			writeJSON(w, response{Message: err.Error()}, http.StatusInternalServerError)
			return
		}

		writeJSON(w, client, http.StatusCreated)
	}
}

type tokenResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
}

func (a *app) tokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		client, err := a.validateTokenHandlerRequest(r)
		if err != nil {
			writeJSON(w, response{Message: err.Error()}, http.StatusBadRequest)
			return
		}

		token, err := a.m.GenerateToken(ctx, client)
		if err == errors.ErrInternalServer {
			writeJSON(w, response{Message: err.Error()}, http.StatusInternalServerError)
			return
		} else if err == errors.ErrInvalidClient || err != nil {
			writeJSON(w, response{Message: err.Error()}, http.StatusUnauthorized)
			return
		}

		writeJSON(w, tokenResponse{AccessToken: token.Access, ExpiresAt: token.ExpiresAt}, http.StatusOK)
	}
}

func (a *app) validateTokenHandlerRequest(r *http.Request) (*models.Client, error) {
	gt := r.FormValue("grant_type")
	if gt != "client_credentials" {
		return nil, errors.ErrUnsupportedGrantType
	}

	clientID := r.Form.Get("client_id")
	if clientID == "" {
		return nil, errors.ErrInvalidClient
	}
	clientSecret := r.Form.Get("client_secret")

	return &models.Client{ID: clientID, Secret: clientSecret}, nil
}

func (a *app) tokenValidationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, a.validateBearerToken(r), http.StatusOK)
	}
}

func (a *app) validateBearerToken(r *http.Request) bool {
	ctx := r.Context()
	auth := r.Header.Get("Authorization")
	prefix := "Bearer "
	token := ""
	if len(auth) > 0 && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	} else {
		return false
	}

	return a.m.ValidateToken(ctx, token)
}

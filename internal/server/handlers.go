package server

import (
	"net/http"
	"oauth/internal/errors"
	"oauth/internal/models"
)

func (s *server) registerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		client, err := s.m.RegisterClient(ctx)
		if err != nil {
			WriteJSON(w, Response{Message: err.Error()}, http.StatusInternalServerError)
			return
		}

		WriteJSON(w, client, http.StatusCreated)
	}
}

func (s *server) tokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		client, err := s.validateTokenRequest(r)
		if err != nil {
			WriteJSON(w, Response{Message: err.Error()}, http.StatusBadRequest)
			return
		}

		token, err := s.m.GenerateToken(ctx, client)
		if err == errors.ErrInternalServer {
			WriteJSON(w, Response{Message: err.Error()}, http.StatusInternalServerError)
			return
		} else if err == errors.ErrInvalidClient || err != nil {
			WriteJSON(w, Response{Message: err.Error()}, http.StatusUnauthorized)
			return
		}

		WriteJSON(w, token, http.StatusOK)
	}
}

func (s *server) validateTokenRequest(r *http.Request) (*models.Client, error) {
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

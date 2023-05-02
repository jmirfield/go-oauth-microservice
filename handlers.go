package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/google/uuid"
)

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (a *App) TokenHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		a.osrv.HandleTokenRequest(w, r)
	}
}

type RegisterResponse struct {
	Response
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (a *App) RegisterHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		clientID := uuid.New().String()
		clientSecret := uuid.New().String()
		err := a.cs.Create(&models.Client{
			ID:     clientID,
			Secret: clientSecret,
			Domain: "http://localhost",
		})
		if err != nil {
			WriteJSON(w, Response{Message: err.Error()}, http.StatusInternalServerError)
			return
		}

		WriteJSON(w, RegisterResponse{Response: Response{Message: "Success"}, ClientID: clientID, ClientSecret: clientSecret}, http.StatusCreated)
	}
}

func (a *App) ValidateTokenHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := a.osrv.ValidationBearerToken(r)
		if err != nil {
			WriteJSON(w, Response{Message: "Invalid request"}, http.StatusUnauthorized)
			return
		}
		WriteJSON(w, Response{Message: "Valid"}, http.StatusOK)
	}
}

func (a *App) PublicKeyHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, Response{Message: "Public Key", Data: a.publicPEM}, http.StatusOK)
	}
}

func WriteJSON(w http.ResponseWriter, data any, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

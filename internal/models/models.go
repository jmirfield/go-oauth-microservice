package models

import "time"

type Client struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
}

type Token struct {
	Access    string    `json:"access_token"`
	ExpiresAt time.Time `json:"expires_at"`
}

package main

import (
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	oauth "oauth/api"
	"oauth/pkg/jwt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr = "localhost:3000"

type client struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type token struct {
	Token string `json:"access_token"`
}

func main() {
	// register client to retrieve client id and secret
	c := registerClient()

	// get token using client details
	token := getNewToken(c).Token

	// establish grpc connection
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// request public key using grpc
	grpcClient := oauth.NewAuthClient(conn)
	keyResp, err := grpcClient.GetKey(context.Background(), &oauth.KeyRequest{})
	if err != nil {
		panic(err)
	}
	block, _ := pem.Decode(keyResp.Key)
	k, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	// validate jwt using public key
	validator := jwt.NewValidator(k.(*rsa.PublicKey))
	_, err = validator.Validate(token)
	if err != nil {
		panic(err)
	}

	// validate jwt using grpc
	tokenResp, err := grpcClient.ValidateToken(context.Background(), &oauth.TokenRequest{Token: token})
	if err != nil || !tokenResp.Valid {
		panic("token is not valid")
	}

	fmt.Println(token)
}

func registerClient() *client {
	var c client

	resp, err := http.Post(fmt.Sprintf("http://%s/v1/register", addr), "application/json", bytes.NewBuffer([]byte{}))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&c)
	if err != nil {
		panic(err)
	}

	return &c
}

func getNewToken(c *client) *token {
	var t token

	resp, err := http.Get(fmt.Sprintf("http://%s/v1/token?grant_type=client_credentials&client_id=%s&client_secret=%s", addr, c.ClientID, c.ClientSecret))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&t)
	if err != nil {
		panic(err)
	}

	return &t
}

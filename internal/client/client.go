//go:generate go run github.com/golang/mock/mockgen -package=mocks -source=$GOFILE -destination=../../test/mocks/client.go
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/flan6/rdstation/entity"
)

type client struct {
	httpClient *http.Client
	secret     entity.Secret
}

type Client interface {
	Request(path, method string, data []byte) ([]byte, error)
}

func NewClient(secret entity.Secret, endpoint oauth2.Endpoint) (Client, error) {
	data, err := json.Marshal(secret)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post("https://api.rd.services/auth/token", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var token *entity.Token
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, err
	}

	config := oauth2.Config{
		ClientID:     secret.ClientID,
		ClientSecret: secret.ClientSecret,
		Endpoint:     endpoint,
	}

	httpClient := config.Client(context.Background(), token.Auth2Token())

	return client{
		httpClient: httpClient,
		secret:     secret,
	}, nil
}

func (c client) Request(path, method string, data []byte) ([]byte, error) {
	req, err := http.NewRequest(method, path, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var e RDError
	if response.StatusCode != http.StatusOK {
		err := json.NewDecoder(response.Body).Decode(&e)
		if err != nil {
			return nil, err
		}

		e.Errors.StatusCode = response.StatusCode

		return nil, e
	}

	return io.ReadAll(response.Body)
}

package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"github.com/flan6/rdstation/entity"
)

const (
	authTokenURL = "https://api.rd.services/auth/token"
)

func httpReqMockToken(t *testing.T, statusCode int, token *oauth2.Token) httpmock.Responder {
	return func(request *http.Request) (*http.Response, error) {
		if statusCode == 420 {
			return &http.Response{
				StatusCode: 200,
				Body: io.NopCloser(
					bytes.NewBufferString(string("\xef")),
				),
			}, nil
		}

		if statusCode != 200 {
			return nil, errors.New("")
		}

		tokenRaw, err := json.Marshal(token)
		if err != nil {
			t.Error(err)
		}

		return &http.Response{
			StatusCode: 200,
			Body: io.NopCloser(
				bytes.NewBufferString(string(tokenRaw)),
			),
		}, nil
	}
}

func httpReqMockRequest(t *testing.T, statusCode int, returnString string) httpmock.Responder {
	return func(request *http.Request) (*http.Response, error) {
		if statusCode == 404 {
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body: io.NopCloser(
					bytes.NewBufferString(`{
						"errors": {
							"error_type": "RESOURCE_NOT_FOUND",
							"error_message": "The resource could not be found"
						}
					}`),
				),
			}, nil
		}

		if statusCode == 201 {
			return &http.Response{
				StatusCode: 201,
				Body: io.NopCloser(
					bytes.NewBufferString(returnString),
				),
			}, nil
		}

		if statusCode != 200 {
			return nil, errors.New("")
		}

		return &http.Response{
			StatusCode: 200,
			Body: io.NopCloser(
				bytes.NewBufferString(returnString),
			),
		}, nil
	}
}

func TestClient_NewClient(t *testing.T) {
	secret := entity.Secret{
		ClientID:     "identificador-de-um-cliente-nervoso",
		ClientSecret: "shhhhh",
		RefreshToken: "refresh token",
	}

	endpoint := oauth2.Endpoint{
		AuthURL:   "uma url",
		TokenURL:  "uma url de token",
		AuthStyle: oauth2.AuthStyleInParams,
	}

	httpmock.Activate()
	defer httpmock.Deactivate()

	token := &oauth2.Token{AccessToken: "token", RefreshToken: "token2"}
	t.Run("success", func(t *testing.T) {
		httpmock.RegisterResponder("POST", authTokenURL, httpReqMockToken(t, 200, token))

		cl, err := NewClient(secret, endpoint)
		require.NoError(t, err)
		require.NotNil(t, cl)
	})

	t.Run("error get token", func(t *testing.T) {
		httpmock.RegisterResponder("POST", authTokenURL, httpReqMockToken(t, 404, token))

		cl, err := NewClient(secret, endpoint)
		require.Error(t, err)
		require.Nil(t, cl)
	})

	t.Run("error marshal token", func(t *testing.T) {
		httpmock.RegisterResponder("POST", authTokenURL, httpReqMockToken(t, 420, token))

		cl, err := NewClient(secret, endpoint)
		require.Error(t, err)
		require.Nil(t, cl)
	})
}

func TestClient_Request(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	server := httptest.NewServer(nil)
	defer server.Close()

	secret := entity.Secret{
		ClientID:     "identificador-de-um-cliente-nervoso",
		ClientSecret: "ssshhhhh",
		RefreshToken: "refresh token",
	}

	endpoint := oauth2.Endpoint{
		AuthURL:   fmt.Sprintf("%s/auth", server.URL),
		TokenURL:  fmt.Sprintf("%s/token", server.URL),
		AuthStyle: oauth2.AuthStyleAutoDetect,
	}

	httpmock.Activate()
	defer httpmock.Deactivate()

	token := &oauth2.Token{AccessToken: "token", RefreshToken: "token2"}
	httpmock.RegisterResponder("POST", authTokenURL, httpReqMockToken(t, 200, token))

	cl, err := NewClient(secret, endpoint)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		url := fmt.Sprintf("%s/test/success", server.URL)
		httpmock.RegisterResponder("GET", url, httpReqMockRequest(t, 200, "test"))

		result, err := cl.Request(url, http.MethodGet, nil)
		require.NoError(t, err)
		require.Equal(t, "test", string(result))
	})

	t.Run("error response", func(t *testing.T) {
		url := fmt.Sprintf("%s/test/response", server.URL)
		httpmock.RegisterResponder("GET", url, httpReqMockRequest(t, 400, ""))

		result, err := cl.Request(url, http.MethodGet, nil)
		require.Error(t, err)
		require.Nil(t, result)
	})

	t.Run("error decode", func(t *testing.T) {
		url := fmt.Sprintf("%s/test/response", server.URL)
		httpmock.RegisterResponder("GET", url, httpReqMockRequest(t, 201, "\xef"))

		result, err := cl.Request(url, http.MethodGet, nil)
		require.Error(t, err)
		require.Nil(t, result)
	})

	t.Run("error 404", func(t *testing.T) {
		url := fmt.Sprintf("%s/test/error-404", server.URL)
		httpmock.RegisterResponder("GET", url, httpReqMockRequest(t, 404, ""))

		result, err := cl.Request(url, http.MethodGet, nil)
		e, _ := err.(RDError)
		require.Equal(t, http.StatusNotFound, e.Errors.StatusCode)
		require.Error(t, err)
		require.Nil(t, result)
	})
}

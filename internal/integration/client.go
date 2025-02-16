package integration

import (
	"avito-shop/internal/models"
	"avito-shop/internal/models/requests"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
	_ "github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"testing"
)

func registerClient(t *testing.T, user models.User) *resty.Client {
	b, err := json.Marshal(requests.AuthRequest{Username: user.Username, Password: user.Password})
	require.NoError(t, err)

	client := makeUnauthorisedClient()

	resp, err := http.Post(client.BaseURL+"/auth", "application/json", bytes.NewReader(b))

	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var response map[string]string
	err = json.Unmarshal(bodyBytes, &response)
	require.NoError(t, err)

	token, exists := response["token"]
	require.True(t, exists, "Токен отсутствует в ответе")

	client.SetAuthToken(token)

	return client
}

func makeUnauthorisedClient() *resty.Client {
	client := resty.New()
	endpoint := "http://avito-shop-service:8080/api"
	client.SetBaseURL(endpoint)

	return client
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(b)[:length]
}

func generateUserCredentials() models.User {
	username := "user" + generateRandomString(5)
	password := "pass" + generateRandomString(8)
	return models.User{Username: username, Password: password}
}

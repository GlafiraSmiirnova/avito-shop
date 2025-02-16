package integration

import (
	"avito-shop/config/db"
	"avito-shop/internal/models"
	"avito-shop/internal/models/requests"
	"avito-shop/internal/models/responses"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type AuthSuite struct {
	suite.Suite
	client *resty.Client
	user   models.User
}

func (s *AuthSuite) SetupSuite() {
	db.ConnectTestDB()
	db.ResetTestDB()
	s.user = generateUserCredentials()
	s.client = registerClient(s.T(), s.user)
}

func (s *AuthSuite) TestAuthInvalidToken() {
	resp, err := s.client.SetAuthToken("invalidToken").R().SetResult(new(responses.InfoResponse)).Get("/info")
	s.Require().NoError(err)

	s.Require().Equal(http.StatusUnauthorized, resp.StatusCode(), "getting information should not be available with an invalid token")
}

func (s *AuthSuite) TestAuthWrongPassword() {
	authReq := requests.AuthRequest{
		Username: s.user.Username,
		Password: "wrongPassword",
	}

	resp, err := s.client.R().SetBody(authReq).Post("/auth")
	s.Require().NoError(err)

	s.Require().Equal(http.StatusBadRequest, resp.StatusCode(), "authorising with an invalid password for an existing user should not be allowed")
}

func (s *AuthSuite) TearDownSuite() {
	db.CloseDB()
}

func TestAuth(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}

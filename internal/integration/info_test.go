package integration

import (
	"avito-shop/config/db"
	"avito-shop/internal/models/responses"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type InfoSuite struct {
	suite.Suite
	client       *resty.Client
	startBalance int
}

func (s *InfoSuite) SetupSuite() {
	db.ConnectTestDB()
	db.ResetTestDB()
	s.client = registerClient(s.T(), generateUserCredentials())
	s.startBalance = 1000
}

func (s *InfoSuite) TestInfoNewUser() {
	s.client = registerClient(s.T(), generateUserCredentials())

	resp, err := s.client.R().SetResult(new(responses.InfoResponse)).Get("/info")
	s.Require().NoError(err)

	s.Require().Equal(http.StatusOK, resp.StatusCode())

	info, ok := resp.Result().(*responses.InfoResponse)
	s.Require().True(ok)
	s.Require().Equal(s.startBalance, info.Coins, "balance of a newly created user should be zero")
	s.Require().Empty(info.CoinHistory.Received, "received transactions of a newly created user should be empty")
	s.Require().Empty(info.CoinHistory.Sent, "sent transactions of a newly created user should be empty")
	s.Require().Empty(info.Inventory, "inventory of a newly created user should be empty")
}

func (s *InfoSuite) TestInfoUnauthorised() {
	s.client = makeUnauthorisedClient()

	resp, err := s.client.R().SetResult(new(responses.InfoResponse)).Get("/info")
	s.Require().NoError(err)

	s.Require().Equal(http.StatusUnauthorized, resp.StatusCode(), "getting information should not be available for unauthorised users")
}

func (s *InfoSuite) TearDownSuite() {
	db.CloseDB()
}

func TestInfo(t *testing.T) {
	suite.Run(t, new(InfoSuite))
}

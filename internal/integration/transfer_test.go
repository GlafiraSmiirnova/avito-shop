package integration

import (
	"avito-shop/config/db"
	"avito-shop/internal/models"
	"avito-shop/internal/models/requests"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type TransferSuite struct {
	suite.Suite
	clientA      *resty.Client
	clientB      *resty.Client
	userA        models.User
	userB        models.User
	startBalance int
}

func (s *TransferSuite) SetupSuite() {
	db.ConnectTestDB()
	db.ResetTestDB()
	s.startBalance = 1000
	s.registerUsers()
}

func (s *TransferSuite) TestTransfer() {
	s.registerUsers()

	transactionSent := models.SentTransaction{
		ToUser: s.userB.Username,
		Amount: s.startBalance / 2,
	}

	transactionReceived := models.ReceivedTransaction{
		FromUser: s.userA.Username,
		Amount:   s.startBalance / 2,
	}

	resp, err := s.clientA.R().SetBody(transactionSent).Post("/sendCoin")
	s.Require().NoError(err)

	s.Require().Equal(http.StatusOK, resp.StatusCode())

	infoUserA := getInfo(&s.Suite, s.clientA)
	infoUserB := getInfo(&s.Suite, s.clientB)

	s.Require().Equal(s.startBalance-transactionSent.Amount, infoUserA.Coins)
	s.Require().Equal(s.startBalance+transactionSent.Amount, infoUserB.Coins)

	s.Require().Contains(infoUserA.CoinHistory.Sent, transactionSent)
	s.Require().Contains(infoUserB.CoinHistory.Received, transactionReceived)
}

func (s *TransferSuite) TestTransferOverBalanceLimit() {
	s.registerUsers()

	transaction := requests.TransferRequest{
		ToUser: s.userB.Username,
		Amount: s.startBalance * 2,
	}

	resp, err := s.clientA.R().SetBody(transaction).Post("/sendCoin")
	s.Require().NoError(err)

	s.Require().Equal(http.StatusBadRequest, resp.StatusCode(), "sending coins exceeding user's balance should not be allowed")

	infoUserA := getInfo(&s.Suite, s.clientA)
	infoUserB := getInfo(&s.Suite, s.clientB)

	s.Require().Equal(s.startBalance, infoUserA.Coins)
	s.Require().Equal(s.startBalance, infoUserB.Coins)

	s.Require().Empty(infoUserA.CoinHistory.Sent)
	s.Require().Empty(infoUserB.CoinHistory.Received)
}

func (s *TransferSuite) TestTransferUnauthorised() {
	s.clientA = makeUnauthorisedClient()

	transaction := requests.TransferRequest{
		ToUser: s.userB.Username,
		Amount: s.startBalance / 2,
	}

	resp, err := s.clientA.R().SetBody(transaction).Post("/sendCoin")
	s.Require().NoError(err)

	s.Require().Equal(http.StatusUnauthorized, resp.StatusCode(), "sending coins should not be allowed for unauthorised users")
}

func (s *TransferSuite) TestTransferReceiverNotFound() {
	s.clientA = registerClient(s.T(), s.userA)

	transaction := requests.TransferRequest{
		ToUser: "nonExistentUser",
		Amount: s.startBalance / 2,
	}

	resp, err := s.clientA.R().SetBody(transaction).Post("/sendCoin")
	s.Require().NoError(err)

	s.Require().Equal(http.StatusBadRequest, resp.StatusCode(), "sending coins to non-existent users should not be allowed")

	infoUserA := getInfo(&s.Suite, s.clientA)

	s.Require().Equal(s.startBalance, infoUserA.Coins)

	s.Require().Empty(infoUserA.CoinHistory.Sent)
}

func (s *TransferSuite) TestTransferIncorrectAmount() {
	s.registerUsers()

	transaction := requests.TransferRequest{
		ToUser: s.userB.Username,
		Amount: 0,
	}

	resp, err := s.clientA.R().SetBody(transaction).Post("/sendCoin")
	s.Require().NoError(err)

	s.Require().Equal(http.StatusBadRequest, resp.StatusCode(), "making a transfer with zero amount should not be allowed")

	transaction = requests.TransferRequest{
		ToUser: s.userB.Username,
		Amount: -s.startBalance,
	}

	resp, err = s.clientA.R().SetBody(transaction).Post("/sendCoin")
	s.Require().NoError(err)

	s.Require().Equal(http.StatusBadRequest, resp.StatusCode(), "making a transfer with negative amount should not be allowed")

	infoUserA := getInfo(&s.Suite, s.clientA)
	infoUserB := getInfo(&s.Suite, s.clientB)

	s.Require().Equal(s.startBalance, infoUserA.Coins)
	s.Require().Equal(s.startBalance, infoUserB.Coins)

	s.Require().Empty(infoUserA.CoinHistory.Sent)
	s.Require().Empty(infoUserB.CoinHistory.Received)
}

func (s *TransferSuite) TestTransferMalformedRequest() {
	s.clientA = registerClient(s.T(), s.userA)

	transaction := requests.TransferRequest{
		ToUser: s.userB.Username,
	}

	resp, err := s.clientA.R().SetBody(transaction).Post("/sendCoin")
	s.Require().NoError(err)

	s.Require().Equal(http.StatusBadRequest, resp.StatusCode(), "making a transfer without an amount should not be allowed")

	transaction = requests.TransferRequest{
		Amount: s.startBalance,
	}

	resp, err = s.clientA.R().SetBody(transaction).Post("/sendCoin")
	s.Require().NoError(err)

	s.Require().Equal(http.StatusBadRequest, resp.StatusCode(), "making a transfer without a receiver should not be allowed")
}

func (s *TransferSuite) registerUsers() {
	s.userA = generateUserCredentials()
	s.userB = generateUserCredentials()

	s.clientA = registerClient(s.T(), s.userA)
	s.clientB = registerClient(s.T(), s.userB)
}

func (s *TransferSuite) TearDownSuite() {
	db.CloseDB()
}

func TestTransfer(t *testing.T) {
	suite.Run(t, new(TransferSuite))
}

package integration

import (
	"avito-shop/config/db"
	"avito-shop/internal/models"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type BuyMerchSuite struct {
	suite.Suite
	client       *resty.Client
	realMerch    models.Merch
	startBalance int
}

func (s *BuyMerchSuite) SetupSuite() {
	db.ConnectTestDB()
	db.ResetTestDB()
	s.realMerch = models.Merch{
		Name:  "pen",
		Price: 10,
	}
	s.startBalance = 1000
}

func (s *BuyMerchSuite) TestBuyMerch() {
	s.client = registerClient(s.T(), generateUserCredentials())
	resp, err := s.client.R().SetPathParam("item", s.realMerch.Name).Get("/buy/{item}")
	s.Require().NoError(err)

	s.Require().Equal(http.StatusOK, resp.StatusCode())

	info := getInfo(&s.Suite, s.client)

	s.Require().Equal(s.startBalance-s.realMerch.Price, info.Coins)

	boughtItem := models.InventoryItem{
		Type:     s.realMerch.Name,
		Quantity: 1,
	}
	s.Require().Contains(info.Inventory, boughtItem)
}

func (s *BuyMerchSuite) TestBuyMerchUnauthorised() {
	s.client = makeUnauthorisedClient()

	resp, err := s.client.R().SetPathParam("item", s.realMerch.Name).Get("/buy/{item}")
	s.Require().NoError(err)

	s.Require().Equal(http.StatusUnauthorized, resp.StatusCode(), "buying merch should not be available for unauthorised users")
}

func (s *BuyMerchSuite) TestBuyMerchOverBalanceLimit() {
	s.client = registerClient(s.T(), generateUserCredentials())

	expensiveMerch := models.Merch{
		Name:  "pink-hoody",
		Price: 500,
	}

	resp, err := s.client.R().SetPathParam("item", expensiveMerch.Name).Get("/buy/{item}")
	s.Require().NoError(err)

	s.Require().Equal(http.StatusOK, resp.StatusCode())

	resp, err = s.client.R().SetPathParam("item", expensiveMerch.Name).Get("/buy/{item}")
	s.Require().NoError(err)

	s.Require().Equal(http.StatusOK, resp.StatusCode(), "having a 0 balance should not cause an error")

	resp, err = s.client.R().SetPathParam("item", expensiveMerch.Name).Get("/buy/{item}")
	s.Require().NoError(err)

	s.Require().Equal(http.StatusBadRequest, resp.StatusCode(), "buying an item more expensive than user's balance should no be allowed")

	info := getInfo(&s.Suite, s.client)
	s.Require().Equal(0, info.Coins)

	boughtItems := models.InventoryItem{
		Type:     expensiveMerch.Name,
		Quantity: 2,
	}
	s.Require().Contains(info.Inventory, boughtItems)
}

func (s *BuyMerchSuite) TestBuyMerchWrongItem() {
	s.client = registerClient(s.T(), generateUserCredentials())

	resp, err := s.client.R().SetPathParam("item", "nonExistentItem").Get("/buy/{item}")
	s.Require().NoError(err)

	s.Require().Equal(http.StatusBadRequest, resp.StatusCode(), "buying a non-existent item should not be allowed")
}

func (s *BuyMerchSuite) TearDownSuite() {
	db.CloseDB()
}

func TestBuyMerch(t *testing.T) {
	suite.Run(t, new(BuyMerchSuite))
}

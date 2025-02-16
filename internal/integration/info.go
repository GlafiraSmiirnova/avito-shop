package integration

import (
	"avito-shop/internal/models/responses"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"
	"net/http"
)

func getInfo(s *suite.Suite, client *resty.Client) *responses.InfoResponse {
	resp, err := client.R().SetResult(new(responses.InfoResponse)).Get("/info")
	s.Require().NoError(err)

	s.Require().Equal(http.StatusOK, resp.StatusCode())

	info, ok := resp.Result().(*responses.InfoResponse)
	s.Require().True(ok)

	return info
}

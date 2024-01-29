package http

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/k1nky/gophkeeper/internal/adapter/http/mock"
	"github.com/stretchr/testify/suite"
)

type adapterTestSuite struct {
	suite.Suite
	authService *mock.MockauthService
}

func TestAdapter(t *testing.T) {
	suite.Run(t, new(adapterTestSuite))
}

func (suite *adapterTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.authService = mock.NewMockauthService(ctrl)
}

package adapter

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
	"spotify_challenge.com/src/models"
)

type AdapterUnitTestSuite struct {
	suite.Suite
	Adapter AdapterImpl
}

func TestAdapterUnitTestSuite(t *testing.T) {
	suite.Run(t, new(AdapterUnitTestSuite))
}

func (s *AdapterUnitTestSuite) SetupSuite() {
	client := http.Client{}
	config := models.NewConfig()

	s.Adapter = AdapterImpl{
		client: &client,
		config: config.SpotifyData,
	}
}

func (s *AdapterUnitTestSuite) AfterTest(suiteName, testName string) {

}

func (s *AdapterUnitTestSuite) Test_NewAdapterImpl_Success() {
	//arrange
	client := http.Client{}
	config := models.NewConfig()

	//act
	adapter := NewAdapter(&client, config.SpotifyData)

	//assert
	s.Assert().IsType(AdapterImpl{}, adapter)
}

func (s *AdapterUnitTestSuite) Test_Get_Success() {
	//arrange

	//act
	m, err := s.Adapter.Get("11dFghVXANMlKmJXsNCbNl")

	//assert
	s.Assert().Nil(err)
	s.Assert().NotNil(m)
	s.Assert().Equal(m.Album.Name, "Cut To The Feeling")
}

func (s *AdapterUnitTestSuite) Test_Get_Invalid_ISRC() {
	//arrange

	//act
	m, err := s.Adapter.Get("11dFghVXANMlKmJXsNCbNlzzz")

	//assert
	s.Assert().NotNil(err)
	s.Assert().Nil(m)
	s.Assert().Equal(err.Error(), "400 Bad Request")
}

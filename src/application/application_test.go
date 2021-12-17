package application

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"spotify_challenge.com/src/adapter"
	"spotify_challenge.com/src/connector"
	"spotify_challenge.com/src/models"
	"spotify_challenge.com/src/service"
)

type ApplicationrUnitTestSuite struct {
	suite.Suite
	Application ApplicationImpl
	Adapter     Adapter
	Service     Service
}

func TestApplicationUnitTestSuite(t *testing.T) {
	suite.Run(t, new(ApplicationrUnitTestSuite))
}

func (s *ApplicationrUnitTestSuite) SetupSuite() {
	client := &http.Client{}
	config := models.NewConfig()
	adapter := adapter.NewAdapter(client, config.SpotifyData)

	db := sql.DB{}
	conn := connector.NewConnector(&db)
	service := service.NewService(conn)

	s.Application = ApplicationImpl{
		Adapter: adapter,
		Service: service,
	}
}

func (s *ApplicationrUnitTestSuite) AfterTest(suiteName, testName string) {
	s.Adapter.AssertExpectations(s.T())
	s.Adapter.ExpectedCalls = nil

	s.Service.AssertExpectations(s.T())
	s.Service.ExpectedCalls = nil
}

type Adapter struct {
	mock.Mock
}

func (a *Adapter) Get(ISRC string) (*models.Metadata, error) {
	args := a.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*models.Metadata), args.Error(1)
	}
}

type Service struct {
	mock.Mock
}

func (s *Service) Write(ctx context.Context, metadata *models.Metadata) error {
	args := s.Called()

	return args.Error(0)
}
func (s *Service) ReadByISRC(ctx context.Context, ISRC string) ([]models.DbTracks, error) {
	args := s.Called()

	return args.Get(0).([]models.DbTracks), args.Error(1)
}
func (s *Service) ReadByArtist(ctx context.Context, name string) ([]models.DbTracks, error) {
	args := s.Called()

	return args.Get(0).([]models.DbTracks), args.Error(1)
}

func (s *ApplicationrUnitTestSuite) Test_NewApplication_Success() {
	//arrange
	adapter := adapter.AdapterImpl{}
	service := service.ServiceImpl{}

	//act
	result := NewApplication(adapter, service)

	//assert
	s.Assert().IsType(ApplicationImpl{}, result)
}

func (s *ApplicationrUnitTestSuite) Test_Write_Invalid_ISRC() {
	//arrange
	ctx := context.Background()

	//act
	err := s.Application.Write(ctx, "")

	//assert
	s.Assert().NotNil(err)
	s.Assert().Equal("invalid ISRC value", err.Error())
}

func (s *ApplicationrUnitTestSuite) Test_Write_Adapter_Get_Failed() {
	//arrange
	ctx := context.Background()
	s.Adapter.On("Get").Return(nil, errors.New("fake error"))

	s.Application.Adapter = &s.Adapter

	//act
	err := s.Application.Write(ctx, "11dFghVXANMlKmJXsNCbNl")

	//assert
	s.Assert().NotNil(err)
	s.Assert().Equal("fake error", err.Error())
}

func (s *ApplicationrUnitTestSuite) Test_Write_Service_Write_Failed() {
	//arrange
	ctx := context.Background()

	s.Adapter.On("Get").Return(&models.Metadata{}, nil)
	s.Application.Adapter = &s.Adapter

	s.Service.On("Write").Return(errors.New("fake error"))
	s.Application.Service = &s.Service

	//act
	err := s.Application.Write(ctx, "11dFghVXANMlKmJXsNCbNl")

	//assert
	s.Assert().NotNil(err)
	s.Assert().Equal("fake error", err.Error())
}

func (s *ApplicationrUnitTestSuite) Test_Write_Service_Write_Success() {
	//arrange
	ctx := context.Background()

	s.Adapter.On("Get").Return(&models.Metadata{}, nil)
	s.Application.Adapter = &s.Adapter

	s.Service.On("Write").Return(nil)
	s.Application.Service = &s.Service

	//act
	err := s.Application.Write(ctx, "11dFghVXANMlKmJXsNCbNl")

	//assert
	s.Assert().Nil(err)
}

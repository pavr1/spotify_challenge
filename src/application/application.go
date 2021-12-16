package application

import (
	"errors"

	"spotify_challenge.com/src/adapter"
	"spotify_challenge.com/src/service"
)

type Application interface {
}

type ApplicationImpl struct {
	Adapter adapter.Adapter
	Service service.Service
}

func NewApplication(adapter adapter.Adapter, service service.Service) Application {
	return ApplicationImpl{
		Adapter: adapter,
		Service: service,
	}
}

func (a ApplicationImpl) Write(ISRC string) error {
	if ISRC == "" {
		return errors.New("invalid ISRC value")
	}

	metadata, err := a.Adapter.Get(ISRC)
	if err != nil {
		return err
	}

	return nil
}

func (a ApplicationImpl) ReadByISRC(ISRC string) (string, error) {
	return "", nil
}

func (a ApplicationImpl) ReadByArtist(name string) (string, error) {
	return "", nil
}

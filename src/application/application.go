package application

import (
	"context"
	"errors"

	"spotify_challenge.com/src/adapter"
	"spotify_challenge.com/src/models"
	"spotify_challenge.com/src/service"
)

type Application interface {
	Write(ctx context.Context, ISRC string) error
	ReadByISRC(ctx context.Context, ISRC string) (*models.DbTracks, error)
	ReadByArtist(ctx context.Context, name string) ([]models.DbTracks, error)
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

func (a ApplicationImpl) Write(ctx context.Context, ISRC string) error {
	if ISRC == "" {
		return errors.New("invalid ISRC value")
	}

	m, err := a.Adapter.Get(ISRC)
	if err != nil {
		return err
	}
	err = a.Service.Write(ctx, m)
	if err != nil {
		return err
	}

	return nil
}

func (a ApplicationImpl) ReadByISRC(ctx context.Context, ISRC string) (*models.DbTracks, error) {
	if ISRC == "" {
		return nil, errors.New("invalid ISRC value")
	}

	tracks, err := a.Service.ReadByISRC(ctx, ISRC)
	if err != nil {
		return nil, err
	}

	if len(tracks) == 0 {
		return nil, nil
	}

	return &tracks[0], nil
}

func (a ApplicationImpl) ReadByArtist(ctx context.Context, name string) ([]models.DbTracks, error) {
	if name == "" {
		return nil, errors.New("invalid ISRC value")
	}

	tracks, err := a.Service.ReadByArtist(ctx, name)
	if err != nil {
		return []models.DbTracks{}, nil
	}

	return tracks, nil
}

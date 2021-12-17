package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"spotify_challenge.com/src/connector"
	"spotify_challenge.com/src/models"
)

type Service interface {
	Write(ctx context.Context, metadata *models.Metadata) error
	ReadByISRC(ctx context.Context, ISRC string) ([]models.DbTracks, error)
	ReadByArtist(ctx context.Context, name string) ([]models.DbTracks, error)
}

type ServiceImpl struct {
	conn connector.Connector
}

func NewService(conn connector.Connector) Service {
	return ServiceImpl{
		conn: conn,
	}
}

func (s ServiceImpl) Write(ctx context.Context, metadata *models.Metadata) error {
	if metadata == nil {
		return errors.New("invalid metadata")
	}

	tracks, err := s.ReadByISRC(ctx, metadata.ExternalIds.Isrc)
	if err != nil {
		return err
	}

	if len(tracks) > 0 {
		return errors.New("invalid ISRC, already existent")
	}

	var imageURL string
	if len(metadata.Album.Images) > 0 {
		imageURL = metadata.Album.Images[0].URL
	}
	artistNames := ""
	for i, val := range metadata.Album.Artists {
		artistNames += val.Name

		if i+1 < len(metadata.Album.Artists)-1 {
			artistNames += ", "
		}
	}

	statement := fmt.Sprintf("INSERT INTO Tracks(ISRC, URI, Title, ArtistNames) VALUES ('%s', '%s', '%s', '%s')", metadata.ExternalIds.Isrc, imageURL, metadata.Name, artistNames)

	err = s.conn.Execute(ctx, statement)
	if err != nil {
		return err
	}

	return nil
}

func (s ServiceImpl) ReadByISRC(ctx context.Context, ISRC string) ([]models.DbTracks, error) {
	statement := fmt.Sprintf("SELECT ISRC, URI, Title, ArtistNames FROM Tracks WHERE ISRC = '%s'", ISRC)
	tracks, err := s.conn.Retrieve(ctx, statement)
	if err != nil {
		return nil, err
	}
	return tracks, nil
}

func (s ServiceImpl) ReadByArtist(ctx context.Context, name string) ([]models.DbTracks, error) {
	if name == "" {
		return []models.DbTracks{}, errors.New("invalid name empty value")
	}

	statement := "SELECT ISRC, URI, Title, ArtistNames FROM Tracks"
	tracks, err := s.conn.Retrieve(ctx, statement)
	if err != nil {
		return nil, err
	}

	//this could be improved in sql directly, no time to do so
	results := []models.DbTracks{}
	for _, t := range tracks {
		if strings.Contains(t.ArtistNames, name) {
			results = append(results, t)
		}
	}

	return results, nil
}

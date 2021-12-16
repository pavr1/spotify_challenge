package service

import (
	"context"
	"errors"

	"spotify_challenge.com/src/connector"
)

type Service interface {
}

type ServiceImpl struct {
	conn connector.Connector
}

func NewService(conn connector.Connector) Service {
	return ServiceImpl{
		conn: conn,
	}
}

func (s ServiceImpl) Write(metadata string) error {
	if metadata == "" {
		return errors.New("invalid metadata")
	}

	//statement := ""

	return nil
}

func (a ServiceImpl) ReadByISRC(ctx context.Context, ISRC string) (string, error) {
	if ISRC == "" {
		return "", errors.New("invalid Isrc empty value")
	}

	// statement := fmt.Sprintf("SELECT (ISRC, URI, Title, ArtistNames) FROM Tracks WHERE ISRC = '%s'", ISRC)
	// a.conn.Retrieve(ctx)

	return "", nil
}

func (a ServiceImpl) ReadByArtist(name string) (string, error) {
	return "", nil
}

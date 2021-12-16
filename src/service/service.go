package service

import (
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

	statement := ""

	return nil
}

func (a ServiceImpl) ReadByISRC(ISRC string) (string, error) {
	return "", nil
}

func (a ServiceImpl) ReadByArtist(name string) (string, error) {
	return "", nil
}

package connector

import (
	"context"
	"database/sql"

	"spotify_challenge.com/src/models"
)

type Connector interface {
	Execute(ctx context.Context, sqlStatement string) error
	Retrieve(ctx context.Context, statement string) ([]models.DbTracks, error)
}

type ConnectorImpl struct {
	db *sql.DB
}

func NewConnector(db *sql.DB) Connector {
	return ConnectorImpl{
		db: db,
	}
}

func (c ConnectorImpl) Execute(ctx context.Context, sqlStatement string) error {
	var err error

	err = c.db.PingContext(ctx)
	if err != nil {
		return err
	}

	query, err := c.db.Prepare(sqlStatement)
	if err != nil {
		return err
	}

	defer query.Close()
	_, err = query.ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c ConnectorImpl) Retrieve(ctx context.Context, statement string) ([]models.DbTracks, error) {
	ctx1 := context.Background()
	err := c.db.PingContext(ctx1)
	if err != nil {
		return nil, err
	}

	rows, err := c.db.QueryContext(ctx, statement)
	if err != nil {
		return nil, err
	}

	result := []models.DbTracks{}

	for rows.Next() {
		var isrc, uri, title, artistNames string

		err = rows.Scan(&isrc, &uri, &title, &artistNames)
		if err != nil {
			return nil, err
		}

		result = append(result, models.DbTracks{
			ISRC:        isrc,
			URI:         uri,
			Title:       title,
			ArtistNames: artistNames,
		})
	}

	return result, nil
}

package connector

import (
	"context"
	"database/sql"

	"msbeer.com/src/models"
)

type Connector interface {
	Execute(ctx context.Context, db *sql.DB, sqlStatement string) error
	Retrieve(ctx context.Context, db *sql.DB, statement string) ([]models.BeerItem, error)
}

type ConnectorImpl struct {
}

func NewConnectorImpl() Connector {
	return ConnectorImpl{}
}

func (ConnectorImpl) Execute(ctx context.Context, db *sql.DB, sqlStatement string) error {
	var err error

	err = db.PingContext(ctx)
	if err != nil {
		return err
	}

	query, err := db.Prepare(sqlStatement)
	if err != nil {
		return err
	}

	defer query.Close()
	newRecord := query.QueryRowContext(ctx)

	var newID int64
	err = newRecord.Scan(&newID)
	if err != nil {
		return err
	}

	return nil
}

func (ConnectorImpl) Retrieve(ctx context.Context, db *sql.DB, statement string) ([]models.BeerItem, error) {
	ctx1 := context.Background()
	err := db.PingContext(ctx1)
	if err != nil {
		return nil, err
	}

	rows, err := db.QueryContext(ctx, statement)
	if err != nil {
		return nil, err
	}

	result := []models.BeerItem{}

	for rows.Next() {
		var name, brewery, country, currency string
		var id int
		var price float64

		err = rows.Scan(&id, &name, &brewery, &country, &price, &currency)
		if err != nil {
			return nil, err
		}

		result = append(result, models.BeerItem{
			ID:       id,
			Name:     name,
			Brewery:  brewery,
			Country:  country,
			Price:    price,
			Currency: currency,
		})
	}

	return result, nil
}

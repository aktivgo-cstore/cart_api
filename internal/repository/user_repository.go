package repository

import (
	"cart_api/internal/models"
	"github.com/jmoiron/sqlx"
)

type CartRepository struct {
	MySqlConn *sqlx.DB
}

func NewCartRepository(mySqlConn *sqlx.DB) *CartRepository {
	return &CartRepository{
		MySqlConn: mySqlConn,
	}
}

func (cr *CartRepository) GetCarts() ([]*models.Cart, error) {
	sql := `
		SELECT * FROM carts
	`

	var carts []*models.Cart
	if err := cr.MySqlConn.Select(&carts, sql); err != nil {
		return nil, err
	}

	return carts, nil
}

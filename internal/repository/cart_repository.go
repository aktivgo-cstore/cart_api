package repository

import (
	"cart_api/internal/dtos"
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

func (cr *CartRepository) GetCarts() ([]*dtos.Cart, error) {
	var carts []*dtos.Cart

	sql := `
		SELECT * FROM carts
	`

	var cartsModel []*models.Cart
	if err := cr.MySqlConn.Select(&cartsModel, sql); err != nil {
		return nil, err
	}

	for _, cartModel := range cartsModel {
		var products []*dtos.CartProduct

		sql = `
			SELECT products.id, title, price, image 
			FROM products
			JOIN cart_to_product 
			WHERE cart_to_product.product_id = ?         
			AND products.id = cart_to_product.product_id 
		`

		if err := cr.MySqlConn.Select(&products, sql, cartModel.ID); err != nil {
			return nil, err
		}

		cart := &dtos.Cart{
			ID:     cartModel.ID,
			UserID: cartModel.UserID,
			Items:  products,
		}

		carts = append(carts, cart)
	}

	return carts, nil
}

func (cr *CartRepository) GetCart(id int) (*dtos.Cart, error) {
	var cart *dtos.Cart

	sql := `
		SELECT * FROM carts
		WHERE id = ?
	`

	var cartsModel []*models.Cart
	if err := cr.MySqlConn.Select(&cartsModel, sql, id); err != nil {
		return nil, err
	}

	if len(cartsModel) <= 0 {
		return cart, nil
	}

	var products []*dtos.CartProduct

	sql = `
		SELECT products.id, title, price, image 
		FROM products
		JOIN cart_to_product 
		WHERE cart_to_product.product_id = ?         
		AND products.id = cart_to_product.product_id 
	`

	if err := cr.MySqlConn.Select(&products, sql, cartsModel[0].ID); err != nil {
		return nil, err
	}

	cart = &dtos.Cart{
		ID:     cartsModel[0].ID,
		UserID: cartsModel[0].UserID,
		Items:  products,
	}

	return cart, nil
}

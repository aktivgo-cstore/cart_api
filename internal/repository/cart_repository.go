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
			WHERE cart_to_product.cart_id = ?         
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

func (cr *CartRepository) GetCartByUserID(userID int) (*dtos.Cart, error) {
	var cart *dtos.Cart

	sql := `
		SELECT * FROM carts
		WHERE user_id = ?
	`

	var cartModel []*models.Cart
	if err := cr.MySqlConn.Select(&cartModel, sql, userID); err != nil {
		return nil, err
	}

	if len(cartModel) <= 0 {
		return cart, nil
	}

	var products []*dtos.CartProduct

	sql = `
		SELECT products.id, title, price, image 
		FROM products
		JOIN cart_to_product 
		WHERE cart_to_product.cart_id = ?         
		AND products.id = cart_to_product.product_id 
	`

	if err := cr.MySqlConn.Select(&products, sql, cartModel[0].ID); err != nil {
		return nil, err
	}

	cart = &dtos.Cart{
		ID:     cartModel[0].ID,
		UserID: cartModel[0].UserID,
		Items:  products,
	}

	return cart, nil
}

func (cr *CartRepository) AddProduct(cartID int, productID int) error {
	sql := `
		INSERT INTO cart_to_product
		(product_id, cart_id)
		VALUE (?, ?)
	`

	if _, err := cr.MySqlConn.Exec(sql, productID, cartID); err != nil {
		return err
	}

	return nil
}

func (cr *CartRepository) CreateCart(userID int) (int, error) {
	sql := `
		INSERT INTO carts
		(user_id)
		VALUE (?)
		
	`

	res, err := cr.MySqlConn.Exec(sql, userID)
	if err != nil {
		return -1, err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(lastInsertID), nil
}

func (cr *CartRepository) GetProduct(cartID int, productID int) (*models.CartToProduct, error) {
	var product []*models.CartToProduct

	sql := `
		SELECT *
		FROM cart_to_product
		WHERE cart_id = ?         
		AND product_id = ?
	`

	if err := cr.MySqlConn.Select(&product, sql, cartID, productID); err != nil {
		return nil, err
	}

	if product == nil {
		return nil, nil
	}

	return product[0], nil
}

func (cr *CartRepository) DeleteProduct(cartID int, productID int) error {
	sql := `
		DELETE FROM cart_to_product
		WHERE cart_id = ? 
		AND product_id = ?
	`

	if _, err := cr.MySqlConn.Exec(sql, cartID, productID); err != nil {
		return err
	}

	return nil
}

func (cr *CartRepository) Delete(cartID int) error {
	sql := `
		DELETE FROM cart_to_product
		WHERE cart_id = ? 
	`

	if _, err := cr.MySqlConn.Exec(sql, cartID); err != nil {
		return err
	}

	return nil
}

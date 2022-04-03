package models

type Cart struct {
	ID     int `db:"id" json:"id"`
	UserID int `db:"user_id" json:"user_id"`
}

type CartToProduct struct {
	ID        int `db:"id" json:"id"`
	ProductID int `db:"product_id" json:"product_id"`
	CartID    int `db:"cart_id" json:"cart_id"`
}

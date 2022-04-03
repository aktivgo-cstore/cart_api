package dtos

type Cart struct {
	ID     int            `json:"id"`
	UserID int            `json:"user_id"`
	Items  []*CartProduct `json:"items"`
}

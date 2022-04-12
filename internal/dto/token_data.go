package dto

type TokenData struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	IsActivated int    `json:"is_activated"`
	Role        string `json:"role"`
}

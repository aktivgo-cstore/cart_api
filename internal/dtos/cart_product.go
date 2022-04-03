package dtos

type CartProduct struct {
	ID    int         `json:"id"`
	Title string      `json:"title"`
	Price int         `json:"price"`
	Image interface{} `json:"image"`
}

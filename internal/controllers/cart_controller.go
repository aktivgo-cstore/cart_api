package controllers

import (
	"cart_api/internal/helpers"
	"cart_api/internal/service"
	"encoding/json"
	"log"
	"net/http"
)

type CartController struct {
	CartService *service.CartService
}

func NewCartController(cartService *service.CartService) *CartController {
	return &CartController{
		CartService: cartService,
	}
}

func (cc *CartController) GetCarts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	carts, er := cc.CartService.GetCarts()
	if er != nil {
		helpers.ErrorResponse(w, er.Message, er.Status)
		return
	}

	encode, err := json.Marshal(carts)
	if err != nil {
		log.Println("unable to encode carts: " + err.Error())
		helpers.ErrorResponse(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(encode)
}

package controllers

import (
	"cart_api/internal/helpers"
	"cart_api/internal/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
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

	token, er := service.GetToken(r.Header)
	if er != nil {
		helpers.ErrorResponse(w, er.Message, er.Status)
		return
	}

	if er = service.CheckAccess(token); er != nil {
		helpers.ErrorResponse(w, er.Message, er.Status)
		return
	}

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

func (cc *CartController) GetCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token, er := service.GetToken(r.Header)
	if er != nil {
		helpers.ErrorResponse(w, er.Message, er.Status)
		return
	}

	idStr, ok := mux.Vars(r)["id"]
	if !ok {
		log.Println("unable to encode var [id]")
		helpers.ErrorResponse(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("unable to convert id")
		helpers.ErrorResponse(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	if er = service.CheckAccess(token); er != nil {
		if er = service.CheckUserAccess(token, id); er != nil {
			helpers.ErrorResponse(w, er.Message, er.Status)
			return
		}
	}

	cart, er := cc.CartService.GetCart(id)
	if er != nil {
		helpers.ErrorResponse(w, er.Message, er.Status)
		return
	}

	encode, err := json.Marshal(cart)
	if err != nil {
		log.Println("unable to encode carts: " + err.Error())
		helpers.ErrorResponse(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(encode)
}

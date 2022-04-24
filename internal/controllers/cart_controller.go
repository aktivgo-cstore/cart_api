package controllers

import (
	"cart_api/internal/dto"
	"cart_api/internal/dtos"
	"cart_api/internal/helpers"
	"cart_api/internal/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
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

	userIDStr, ok := mux.Vars(r)["user_id"]
	if !ok {
		log.Println("unable to encode var [user_id]")
		helpers.ErrorResponse(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Println("unable to convert userID")
		helpers.ErrorResponse(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	if er = service.CheckAccess(token); er != nil {
		if er = service.CheckUserAccess(token, userID); er != nil {
			helpers.ErrorResponse(w, er.Message, er.Status)
			return
		}
	}

	cart, er := cc.CartService.GetCartByUserID(userID)
	if er != nil {
		helpers.ErrorResponse(w, er.Message, er.Status)
		return
	}

	if cart == nil {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{}"))
		return
	}

	if len(cart.Items) == 0 {
		cart.Items = make([]*dtos.CartProduct, 0)
	}

	encode, err := json.Marshal(cart)
	if err != nil {
		log.Println("unable to encode cart: " + err.Error())
		helpers.ErrorResponse(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(encode)
}

func (cc *CartController) AddProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token, er := service.GetToken(r.Header)
	if er != nil {
		helpers.ErrorResponse(w, er.Message, er.Status)
		return
	}

	userIDStr, ok := mux.Vars(r)["user_id"]
	if !ok {
		log.Println("unable to encode var [user_id]")
		helpers.ErrorResponse(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Println("unable to convert userID")
		helpers.ErrorResponse(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	if er = service.CheckAccess(token); er != nil {
		if er = service.CheckUserAccess(token, userID); er != nil {
			helpers.ErrorResponse(w, er.Message, er.Status)
			return
		}
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("unable to read request body: " + err.Error())
		helpers.ErrorResponse(w, "Некорректный запрос", http.StatusInternalServerError)
		return
	}

	var productData *dto.ProductData
	if err = json.Unmarshal(body, &productData); err != nil {
		log.Println("unable to decode request body: " + err.Error())
		helpers.ErrorResponse(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	if er = cc.CartService.AddProduct(userID, productData.ProductID); er != nil {
		helpers.ErrorResponse(w, er.Message, er.Status)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (cc *CartController) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token, er := service.GetToken(r.Header)
	if er != nil {
		helpers.ErrorResponse(w, er.Message, er.Status)
		return
	}

	userIDStr, ok := mux.Vars(r)["user_id"]
	if !ok {
		log.Println("unable to encode var [user_id]")
		helpers.ErrorResponse(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Println("unable to convert userID")
		helpers.ErrorResponse(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	if er = service.CheckAccess(token); er != nil {
		if er = service.CheckUserAccess(token, userID); er != nil {
			helpers.ErrorResponse(w, er.Message, er.Status)
			return
		}
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("unable to read request body: " + err.Error())
		helpers.ErrorResponse(w, "Некорректный запрос", http.StatusInternalServerError)
		return
	}

	var productData *dto.ProductData
	if err = json.Unmarshal(body, &productData); err != nil {
		log.Println("unable to decode request body: " + err.Error())
		helpers.ErrorResponse(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	if er = cc.CartService.DeleteProduct(userID, productData.ProductID); er != nil {
		helpers.ErrorResponse(w, er.Message, er.Status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (cc *CartController) Clear(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token, er := service.GetToken(r.Header)
	if er != nil {
		helpers.ErrorResponse(w, er.Message, er.Status)
		return
	}

	userIDStr, ok := mux.Vars(r)["user_id"]
	if !ok {
		log.Println("unable to encode var [user_id]")
		helpers.ErrorResponse(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Println("unable to convert userID")
		helpers.ErrorResponse(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	if er = service.CheckAccess(token); er != nil {
		if er = service.CheckUserAccess(token, userID); er != nil {
			helpers.ErrorResponse(w, er.Message, er.Status)
			return
		}
	}

	if er = cc.CartService.Clear(userID); er != nil {
		helpers.ErrorResponse(w, er.Message, er.Status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

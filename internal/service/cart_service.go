package service

import (
	"cart_api/internal/dtos"
	"cart_api/internal/errors"
	"cart_api/internal/repository"
	"fmt"
	"log"
)

type CartService struct {
	CartRepository *repository.CartRepository
}

func NewCartService(cartRepository *repository.CartRepository) *CartService {
	return &CartService{
		CartRepository: cartRepository,
	}
}

func (cs *CartService) GetCarts() ([]*dtos.Cart, *errors.ApiError) {
	carts, err := cs.CartRepository.GetCarts()
	if err != nil {
		log.Println("unable to get carts: " + err.Error())
		return nil, errors.InternalServerError(err)
	}

	return carts, nil
}

func (cs *CartService) GetCartByUserID(userID int) (*dtos.Cart, *errors.ApiError) {
	cart, err := cs.CartRepository.GetCartByUserID(userID)
	if err != nil {
		log.Println("unable to get cart: " + err.Error())
		return nil, errors.InternalServerError(err)
	}

	return cart, nil
}

func (cs *CartService) AddProduct(userID int, productID int) *errors.ApiError {
	cart, err := cs.CartRepository.GetCartByUserID(userID)
	if err != nil {
		log.Println("unable to get cart: " + err.Error())
		return errors.InternalServerError(err)
	}

	var cartID int
	if cart == nil {
		cartID, err = cs.CartRepository.CreateCart(userID)
		if err != nil {
			log.Println("unable to create cart: " + err.Error())
			return errors.InternalServerError(err)
		}
	} else {
		cartID = cart.ID
	}

	product, err := cs.CartRepository.GetProduct(cartID, productID)
	if err != nil {
		log.Println("unable to get product: " + err.Error())
		return errors.InternalServerError(err)
	}

	if product != nil {
		return errors.BadRequestError(
			"Данный продукт уже находится в корзине",
			fmt.Errorf("this product esists in the cart"),
		)
	}

	if err = cs.CartRepository.AddProduct(cartID, productID); err != nil {
		log.Println("unable to add product: " + err.Error())
		return errors.InternalServerError(err)
	}

	return nil
}

func (cs *CartService) DeleteProduct(userID int, productID int) *errors.ApiError {
	cart, err := cs.CartRepository.GetCartByUserID(userID)
	if err != nil {
		log.Println("unable to get cart: " + err.Error())
		return errors.InternalServerError(err)
	}

	if cart == nil {
		return errors.BadRequestError(fmt.Sprintf(
			"Корзина с user_id=%d не найдена", userID),
			fmt.Errorf("cart with user_id=%d not found", userID),
		)
	}

	if err = cs.CartRepository.DeleteProduct(cart.ID, productID); err != nil {
		log.Println("unable to delete product: " + err.Error())
		return errors.InternalServerError(err)
	}

	return nil
}

func (cs *CartService) Clear(userID int) *errors.ApiError {
	cart, err := cs.CartRepository.GetCartByUserID(userID)
	if err != nil {
		log.Println("unable to get cart: " + err.Error())
		return errors.InternalServerError(err)
	}

	if cart == nil {
		return errors.BadRequestError(fmt.Sprintf(
			"Корзина с user_id=%d не найдена", userID),
			fmt.Errorf("cart with user_id=%d not found", userID),
		)
	}

	if err = cs.CartRepository.Delete(cart.ID); err != nil {
		log.Println("unable to clear cart: " + err.Error())
		return errors.InternalServerError(err)
	}

	return nil
}

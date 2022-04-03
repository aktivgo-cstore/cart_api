package service

import (
	"cart_api/internal/dtos"
	"cart_api/internal/errors"
	"cart_api/internal/repository"
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

func (cs *CartService) GetCart(id int) (*dtos.Cart, *errors.ApiError) {
	cart, err := cs.CartRepository.GetCart(id)
	if err != nil {
		log.Println("unable to get carts: " + err.Error())
		return nil, errors.InternalServerError(err)
	}

	return cart, nil
}

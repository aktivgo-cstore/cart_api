package service

import (
	"cart_api/internal/errors"
	"cart_api/internal/models"
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

func (cs *CartService) GetCarts() ([]*models.Cart, *errors.ApiError) {
	carts, err := cs.CartRepository.GetCarts()
	if err != nil {
		log.Println("unable to get carts: " + err.Error())
		return nil, errors.InternalServerError(err)
	}

	return carts, nil
}

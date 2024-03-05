package service

import (
	"clean_arch/internal/dto"
	"clean_arch/internal/repository"
	"context"
)

type CartService struct {
	repo repository.Cart
}

func NewCartService(repo repository.Cart) *CartService {
	return &CartService{
		repo: repo,
	}
}

func (s *CartService) GetUserCart(ctx context.Context, userId int) ([]*dto.Cart, error) {
	return s.repo.GetUserCart(ctx, userId)
}

func (s *CartService) UpdateUserCart(ctx context.Context, userId int, product *dto.UpdateCart) error {
	return s.repo.UpdateUserCart(ctx, userId, product)
}

func (s *CartService) DeleteProduct(ctx context.Context, userId int, productUUID string) error {
	return s.repo.DeleteCartProduct(ctx, userId, productUUID)
}

func (s *CartService) ClearUserCart(ctx context.Context, userId int) error {
	return s.repo.ClearUserCart(ctx, userId)
}

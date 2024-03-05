package service

import (
	"clean_arch/internal/dto"
	"clean_arch/internal/repository"
	"context"
)

type FavouriteService struct {
	repo repository.Favourite
}

func NewFavouriteService(repo repository.Favourite) *FavouriteService {
	return &FavouriteService{
		repo: repo,
	}
}

func (s *FavouriteService) GetUserFavourites(ctx context.Context, userId int) ([]*dto.Favourite, error) {
	return s.repo.GetUserFavourites(ctx, userId)
}

func (s *FavouriteService) UpdateUserFavourites(ctx context.Context, userId int, product *dto.UpdateFavourite) error {
	return s.repo.UpdateUserFavourites(ctx, userId, product)
}

func (s *FavouriteService) DeleteFavouriteProduct(ctx context.Context, userId int, productUUID string) error {
	return s.repo.DeleteFavouriteProduct(ctx, userId, productUUID)
}

func (s *FavouriteService) ClearUserFavourite(ctx context.Context, userId int) error {
	return s.repo.ClearUserFavourite(ctx, userId)
}

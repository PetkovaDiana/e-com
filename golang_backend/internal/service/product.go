package service

import (
	"clean_arch/internal/dto"
	"clean_arch/internal/repository"
	"context"
)

type ProductService struct {
	repo repository.Product
}

func NewProductService(repo repository.Product) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) GetAllProductsByParams(ctx context.Context, params *dto.Params) (*dto.Products, error) {
	return s.repo.GetAllProductsByParams(ctx, params)
}

func (s *ProductService) GetProductDetail(ctx context.Context, uuid string) (*dto.ProductInformation, error) {
	return s.repo.GetProductDetail(ctx, uuid)
}

func (s *ProductService) CreateReview(ctx context.Context, reviewDTO *dto.Review, userId int) error {
	return s.repo.CreateReview(ctx, reviewDTO, userId)
}

func (s *ProductService) GetReviews(ctx context.Context, productUUID string, params *dto.Params) (*dto.ProductStatistic, error) {
	return s.repo.GetReviews(ctx, productUUID, params)
}

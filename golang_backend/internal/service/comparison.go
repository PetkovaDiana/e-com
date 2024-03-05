package service

import (
	"clean_arch/internal/dto"
	"clean_arch/internal/repository"
	"context"
)

type ComparisonService struct {
	repo repository.Comparison
}

func NewComparisonService(repo repository.Comparison) *ComparisonService {
	return &ComparisonService{repo: repo}
}

func (s *ComparisonService) GetUserComparison(ctx context.Context, userId int) ([]*dto.Comparison, *dto.Count, error) {
	return s.repo.GetUserComparison(ctx, userId)
}

func (s *ComparisonService) UpdateUserComparison(ctx context.Context, userId int, product *dto.UpdateComparison) error {
	return s.repo.UpdateUserComparison(ctx, userId, product)
}

func (s *ComparisonService) DeleteComparisonProduct(ctx context.Context, userId int, productId string) error {
	return s.repo.DeleteComparisonProduct(ctx, userId, productId)
}

func (s *ComparisonService) DeleteComparisonProductByCategoryUUID(ctx context.Context, userId int, categoryUUID string) error {
	return s.repo.DeleteComparisonProductByCategoryUUID(ctx, userId, categoryUUID)
}

func (s *ComparisonService) ClearUserComparison(ctx context.Context, userId int) error {
	return s.repo.ClearUserComparison(ctx, userId)
}

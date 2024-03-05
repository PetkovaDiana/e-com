package service

import (
	"clean_arch/internal/dto"
	"clean_arch/internal/repository"
	"context"
)

type CategoryService struct {
	repo repository.Category
}

func NewCategoryService(repo repository.Category) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

func (s *CategoryService) GetAllCategories(ctx context.Context) ([]*dto.Category, error) {
	return s.repo.GetAllCategories(ctx)
}

func (s *CategoryService) GetCategoriesById(ctx context.Context, params *dto.CategoryParams) (*dto.CategoryData, error) {
	return s.repo.GetCategoriesById(ctx, params)
}

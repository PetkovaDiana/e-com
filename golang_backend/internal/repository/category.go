package repository

import (
	"clean_arch/internal/dto"
	"clean_arch/internal/models"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sync"
)

type CategoryRepository struct {
	db        *gorm.DB
	log       *logrus.Logger
	mediaRoot string
}

func NewCategoryRepository(db *gorm.DB, log *logrus.Logger, mediaRoot string) *CategoryRepository {
	return &CategoryRepository{
		db:        db,
		log:       log,
		mediaRoot: mediaRoot,
	}
}

func (r *CategoryRepository) GetAllCategories(ctx context.Context) ([]*dto.Category, error) {
	var categoryDTO []*dto.Category
	if err := r.db.WithContext(ctx).Model(&models.Category{}).Where("can_to_view = ?", true).Where("level = 0").Select("uuid, title, image").Order("level asc, title asc").Find(&categoryDTO).Error; err != nil {
		r.log.Errorf("error occurred category: %s", err.Error())
		return nil, err
	}

	wg := &sync.WaitGroup{}

	for _, val := range categoryDTO {
		wg.Add(1)
		go func(val *dto.Category) {
			defer wg.Done()
			val.ImageMediaRoot(r.mediaRoot)
		}(val)
	}
	wg.Wait()
	return categoryDTO, nil
}

func (r *CategoryRepository) GetCategoriesById(ctx context.Context, params *dto.CategoryParams) (*dto.CategoryData, error) {
	var categoryDTO []*dto.Category

	categoryMap := [1]uuid.UUID{}

	rows, _ := r.db.WithContext(ctx).Raw(fmt.Sprintf(`
		select 
		uuid, 
		title,
		image from category 
			where can_to_view = true and uuid = '%s' 
			or parent_uuid = '%s' order by level asc, title asc;`, params.UUID, params.UUID)).Rows()

	for rows.Next() {
		category := &dto.Category{}

		if err := rows.Scan(
			&category.UUID,
			&category.Title,
			&category.Image,
		); err != nil {
			return nil, err
		}
		if categoryMap[0] == uuid.Nil {
			categoryMap[0] = category.UUID
			categoryDTO = append(categoryDTO, category)
		} else {
			categoryDTO[0].SubCategory = append(categoryDTO[0].SubCategory, dto.SubCategory{
				UUID:  category.UUID,
				Title: category.Title,
				Image: category.Image,
			})
		}

		categoryMap[0] = category.UUID
	}

	categoryOrigin, err := r.GetCategoriesOrigin(ctx, params)

	if err != nil {
		return nil, err
	}

	wg := &sync.WaitGroup{}

	for _, val := range categoryDTO[0].SubCategory {
		wg.Add(1)
		go func(val dto.SubCategory) {
			defer wg.Done()
			val.ImageMediaRoot(r.mediaRoot)
		}(val)
	}

	wg.Wait()

	categoryDTO[0].ImageMediaRoot(r.mediaRoot)

	return &dto.CategoryData{
		Categories:   categoryDTO,
		CategoryPath: categoryOrigin,
	}, nil
}

func (r *CategoryRepository) GetCategoriesOrigin(ctx context.Context, params *dto.CategoryParams) ([]*dto.Category, error) {
	var categoryPath []*dto.Category

	rows, _ := r.db.WithContext(ctx).Raw(fmt.Sprintf(`
		WITH RECURSIVE category_tree AS (
    		SELECT uuid, title, can_to_view, parent_uuid, level
    		FROM category
    		WHERE uuid = '%s'
    		UNION ALL
			SELECT c.uuid, c.title, c.can_to_view, c.parent_uuid, c.level
			FROM category c
				JOIN category_tree ct ON c.uuid = ct.parent_uuid
    				WHERE c.level <= ct.level
		)
		SELECT uuid, title
		FROM category_tree where uuid != '%s' order by level asc;`, params.UUID, params.UUID)).Rows()

	for rows.Next() {
		category := &dto.Category{}
		if err := rows.Scan(
			&category.UUID,
			&category.Title,
		); err != nil {
			return nil, err
		}
		categoryPath = append(categoryPath, category)
	}
	return categoryPath, nil
}

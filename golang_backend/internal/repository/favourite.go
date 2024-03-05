package repository

import (
	"clean_arch/internal/dto"
	"clean_arch/internal/models"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FavouriteRepository struct {
	db        *gorm.DB
	log       *logrus.Logger
	mediaRoot string
}

func NewFavouriteRepository(db *gorm.DB, log *logrus.Logger, mediaRoot string) *FavouriteRepository {
	return &FavouriteRepository{
		db:        db,
		log:       log,
		mediaRoot: mediaRoot,
	}
}

func (r *FavouriteRepository) GetUserFavourites(ctx context.Context, userId int) ([]*dto.Favourite, error) {
	var favouriteProducts []*models.FavouriteProduct
	var favouriteProductsDto []*dto.FavouriteProduct
	var favouriteDTO []*dto.Favourite

	r.db.WithContext(ctx).Preload(clause.Associations).
		Joins("inner join favouritesm2ms ug on ug.favourite_product_id = favourite_product.id ").
		Joins("inner join favourite g on g.id= ug.favourite_id ").
		Where("g.user_id = ?", uint(userId)).Find(&favouriteProducts)

	for x := range favouriteProducts {

		result := &dto.FavouriteProduct{
			UUID:     fmt.Sprint(favouriteProducts[x].Product.UUID),
			Title:    favouriteProducts[x].Product.Title,
			Price:    fmt.Sprint(favouriteProducts[x].Product.Price),
			Image:    string(favouriteProducts[x].Product.Image),
			Quantity: fmt.Sprint(favouriteProducts[x].Product.Quantity),
		}
		//result.ImageMediaRoot(r.mediaRoot)
		favouriteProductsDto = append(favouriteProductsDto, result)
	}

	favouriteDTO = append(favouriteDTO, &dto.Favourite{Product: favouriteProductsDto})
	return favouriteDTO, nil
}

func (r *FavouriteRepository) UpdateUserFavourites(ctx context.Context, userId int, product *dto.UpdateFavourite) error {
	var favouriteProducts *models.FavouriteProduct
	var favourite []models.Favourite

	if r.db.WithContext(ctx).Preload(clause.Associations).
		Joins("inner join favouritesm2ms ug on ug.favourite_product_id = favourite_product.id ").
		Joins("inner join favourite g on g.id= ug.favourite_id ").
		Where("g.user_id = ?", uint(userId)).
		Where("product_uuid = ?", product.ProductUUID).Find(&favouriteProducts).
		RowsAffected != 0 {

		return fmt.Errorf("product already in list")
	} else {
		r.db.WithContext(ctx).Where("user_id = ?", userId).Find(&favourite)

		prodUUID, err := uuid.Parse(product.ProductUUID)
		if err != nil {
			return err
		}

		newFavouriteProduct := models.FavouriteProduct{
			ProductUUID: prodUUID,
			Favourites:  favourite,
		}
		r.db.Create(&newFavouriteProduct)
		return nil
	}
}

func (r *FavouriteRepository) DeleteFavouriteProduct(ctx context.Context, userId int, productUUID string) error {
	var favouritesProducts *models.FavouriteProduct

	if err := r.db.WithContext(ctx).Preload(clause.Associations).
		Joins("inner join favouritesm2ms ug on ug.favourite_product_id = favourite_product.id ").
		Joins("inner join favourite g on g.id= ug.favourite_id ").
		Where("g.user_id = ?", uint(userId)).
		Where("product_uuid = ?", productUUID).First(&favouritesProducts).
		Error; err != nil {

		return fmt.Errorf("product not found")

	} else {
		r.db.WithContext(ctx).Model(&favouritesProducts).Association("Favourites").Clear()
		r.db.WithContext(ctx).Unscoped().Delete(&favouritesProducts)
		return nil
	}
}

func (r *FavouriteRepository) ClearUserFavourite(ctx context.Context, userId int) error {
	var favouriteProducts []*models.FavouriteProduct

	if err := r.db.WithContext(ctx).Preload(clause.Associations).
		Joins("inner join favouritesm2ms ug on ug.favourite_product_id = favourite_product.id ").
		Joins("inner join favourite g on g.id= ug.favourite_id ").
		Where("g.user_id = ?", uint(userId)).
		Find(&favouriteProducts).
		Error; err != nil {

		return fmt.Errorf("favourite list already empty")

	} else {
		r.db.WithContext(ctx).Model(&favouriteProducts).Association("Favourites").Clear()
		if err := r.db.WithContext(ctx).Unscoped().Delete(&favouriteProducts).Error; err != nil {
			return fmt.Errorf("favourite list already empty")
		}
		return nil
	}
}

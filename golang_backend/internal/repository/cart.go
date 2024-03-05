package repository

import (
	"clean_arch/internal/dto"
	"clean_arch/internal/models"
	"clean_arch/pkg"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"sync"
)

type CartRepository struct {
	db        *gorm.DB
	log       *logrus.Logger
	mediaRoot string
}

func NewCartRepository(db *gorm.DB, log *logrus.Logger, mediaRoot string) *CartRepository {
	return &CartRepository{
		db:        db,
		mediaRoot: mediaRoot,
	}
}

func (r *CartRepository) GetUserCart(ctx context.Context, userId int) ([]*dto.Cart, error) {
	var cartProductsDto []*dto.CartProduct
	var cartTotalPrice float64
	var cartDTO []*dto.Cart

	r.db.WithContext(ctx).Raw(fmt.Sprintf(`
		select 
		c.id, 
		p.uuid, 
		p.title, 
		p.price, 
		p.image, 
		cart_product.count, 
		sum(p.quantity) as quantity, 
		cart_product.total_price
		from cart_product inner join cartm2ms c2m on cart_product.id = c2m.cart_product_id
			inner join cart c on c.id = c2m.cart_id
			inner join product p on p.uuid = cart_product.product_uuid
				where c.user_id = %d and c.in_order = false 
						group by c.id, p.uuid, p.title, p.price, p.image, cart_product.count, cart_product.total_price, cart_product.id order by cart_product.id`, userId)).
		Scan(&cartProductsDto)

	var wg sync.WaitGroup

	for _, product := range cartProductsDto {
		wg.Add(1)
		go func(p *dto.CartProduct) {
			defer wg.Done()
			newFloatPrice, _ := strconv.ParseFloat(p.TotalPrice, 64)
			cartTotalPrice += newFloatPrice
			//p.ImageMediaRoot(r.mediaRoot)
		}(product)
	}
	wg.Wait()

	cartDTO = append(cartDTO, &dto.Cart{
		TotalPrice: strconv.FormatFloat(cartTotalPrice, 'f', 2, 64),
		Product:    cartProductsDto,
	})

	return cartDTO, nil
}

func (r *CartRepository) UpdateUserCart(ctx context.Context, userId int, product *dto.UpdateCart) error {
	var cartProducts *models.CartProduct
	var cart []models.Cart
	var products *models.Product
	var newTotalPrice float64

	newProductCount, _ := strconv.Atoi(product.Count)
	if err := r.db.Where("uuid = ?", product.ProductUUID).Select("price").First(&products).Error; err != nil {
		return err
	}
	newTotalPrice = float64(newProductCount) * products.Price
	if r.db.Preload(clause.Associations).
		Joins("inner join cartm2ms ug on ug.cart_product_id = cart_product.id").
		Joins("inner join cart g on g.id= ug.cart_id").
		Where("g.in_order = false AND g.user_id = ?", uint(userId)).
		Where("product_uuid = ?", product.ProductUUID).First(&cartProducts).
		RowsAffected != 0 {

		cartProducts.Count = newProductCount
		cartProducts.TotalPrice = pkg.Round(newTotalPrice)
		r.db.Save(&cartProducts)
		return nil
	} else {
		r.db.Where("user_id = ?", userId).Where("in_order=false").Find(&cart)

		prodUUID, err := uuid.Parse(product.ProductUUID)
		if err != nil {
			return err
		}

		newCartProduct := models.CartProduct{
			ProductUUID: prodUUID,
			Count:       newProductCount,
			Carts:       cart,
			TotalPrice:  pkg.Round(newTotalPrice),
		}
		r.db.Create(&newCartProduct)
	}
	return nil
}

func (r *CartRepository) DeleteCartProduct(ctx context.Context, userId int, productUUID string) error {
	var cartProducts *models.CartProduct

	if err := r.db.WithContext(ctx).Preload(clause.Associations).
		Joins("inner join cartm2ms ug on ug.cart_product_id = cart_product.id ").
		Joins("inner join cart g on g.id= ug.cart_id ").
		Where("g.in_order = false AND g.user_id = ?", uint(userId)).
		Where("product_uuid = ?", productUUID).First(&cartProducts).
		Error; err != nil {

		return fmt.Errorf("product not found")

	} else {
		if err := r.db.WithContext(ctx).Model(&cartProducts).Association("Carts").Clear(); err != nil {
			return fmt.Errorf("error occured db")
		}
		r.db.WithContext(ctx).Unscoped().Delete(&cartProducts)
		return nil
	}
}

func (r *CartRepository) ClearUserCart(ctx context.Context, userId int) error {
	var cartProducts []*models.CartProduct

	if err := r.db.WithContext(ctx).Preload(clause.Associations).
		Joins("inner join cartm2ms ug on ug.cart_product_id = cart_product.id ").
		Joins("inner join cart g on g.id= ug.cart_id ").
		Where("g.in_order = false AND g.user_id = ?", uint(userId)).
		Find(&cartProducts).
		Error; err != nil {

		return fmt.Errorf("cart already empty")

	} else {
		r.db.WithContext(ctx).Model(&cartProducts).Association("Carts").Clear()
		if err := r.db.Unscoped().Delete(&cartProducts).Error; err != nil {
			return fmt.Errorf("cart already empty")
		} else {
			return nil
		}

	}
}

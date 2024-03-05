package repository

import (
	"clean_arch/internal/dto"
	"clean_arch/internal/models"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
	"sync"
)

type ComparisonRepository struct {
	db        *gorm.DB
	log       *logrus.Logger
	mediaRoot string
}

func NewComparisonRepository(db *gorm.DB, log *logrus.Logger, mediaRoot string) *ComparisonRepository {
	return &ComparisonRepository{
		db:        db,
		log:       log,
		mediaRoot: mediaRoot,
	}
}

func (r *ComparisonRepository) GetUserComparison(ctx context.Context, userId int) ([]*dto.Comparison, *dto.Count, error) {
	var comparisonDTO []*dto.Comparison
	var comparisonProductsDTO []*dto.ComparisonProduct

	productsByCategory := make(map[string][]dto.ComparisonProduct)
	categoryChar := map[string][]*dto.CategoryCharacteristic{}

	wg := &sync.WaitGroup{}

	errChan1 := make(chan error)
	errChan2 := make(chan error)

	wg.Add(2)

	// Получаем категории для сравнения
	go func() {
		defer wg.Done()
		rows, err := r.db.WithContext(ctx).Debug().Raw(fmt.Sprintf(`
			select
			c2.uuid as category_uuid,
			c2.title as category_title,
			c3.uuid, c3.title
			from comparison_product
    			inner join comparisonm2ms c2m on comparison_product.id = c2m.comparison_product_id
    			inner join comparison c on c.id = c2m.comparison_id
    			inner join category c2 on c2.uuid = comparison_product.category_uuid
    			left join category_characteristic cc on c2.uuid = cc.category_uuid
    			left join characteristic c3 on cc.characteristic_uuid = c3.uuid
        			where c.user_id = %d group by c2.uuid, c3.uuid;`, userId)).Rows()

		defer func() error {
			if err := rows.Close(); err != nil {
				return err
			}
			return nil
		}()
		if err != nil {
			errChan1 <- err
		}

		comparisonMap := map[string]dto.Comparison{}

		for rows.Next() {

			comparison := &dto.Comparison{}
			var categoryChars dto.CategoryCharacteristic

			r.db.WithContext(ctx).ScanRows(rows, &categoryChars)

			r.db.WithContext(ctx).ScanRows(rows, &comparison)

			_, ok := comparisonMap[comparison.CategoryUUID]

			if !ok {
				comparisonMap[comparison.CategoryUUID] = *comparison
				comparisonDTO = append(comparisonDTO, comparison)
			}

			if categoryChars.UUID != "" {
				for _, val := range comparisonDTO {
					if val.CategoryUUID == comparison.CategoryUUID {
						val.Characteristics = append(val.Characteristics, &categoryChars)
						break
					}
				}
			}
		}

		// Создаем мапу категорий с их характеристиками
		for _, category := range comparisonDTO {
			categoryChar[category.CategoryUUID] = category.Characteristics
		}

		close(errChan1)
	}()

	// Получаем все продукты, по которым будем ввести сравнение
	go func() {
		defer wg.Done()
		rows, err := r.db.WithContext(ctx).Debug().Raw(fmt.Sprintf(`
			select 
			comparison_product.product_uuid as uuid,
    		p.title as title,
    		ps.rating,
    		count(r.id) as total_reviews,
    		p.quantity,
    		p.price, 
			'1' as count,
    		p.image,
    		c2.uuid as category_uuid,
    		c2.title as category_title,
   		 	p.base_unit,
			c3.title,
    		c3.uuid,
    		pc.value
			from comparison_product
				inner join comparisonm2ms c2m on comparison_product.id = c2m.comparison_product_id
				inner join comparison c on c.id = c2m.comparison_id
				inner join category c2 on c2.uuid = comparison_product.category_uuid
				inner join product p on comparison_product.product_uuid = p.uuid
				inner join product_statistic ps on p.uuid = ps.product_uuid
				left join product_characteristic pc on p.uuid = pc.product_uuid
				left join characteristic c3 on pc.characteristic_uuid = c3.uuid
				left join review r on p.uuid = r.product_uuid
					where c.user_id = %d 
						group by comparison_product.product_uuid, p.title, p.quantity, p.image, ps.rating, p.price, c2.uuid, c2.title, p.base_unit, c3.uuid, c3.title, pc.value;`, userId)).Rows()
		if err != nil {
			errChan2 <- err
		}

		defer func() error {
			if err := rows.Close(); err != nil {
				return err
			}
			return nil
		}()
		comparisonProductsMap := map[string]dto.ComparisonProduct{}

		for rows.Next() {

			comparisonProductWithCharacteristic := &dto.ComparisonProductWithCharacteristic{}
			r.db.WithContext(ctx).ScanRows(rows, &comparisonProductWithCharacteristic)

			comparisonProduct := comparisonProductWithCharacteristic.ComparisonProduct
			categoryChars := comparisonProductWithCharacteristic.Characteristic

			_, ok := comparisonProductsMap[comparisonProduct.UUID]

			if !ok {
				comparisonProductsMap[comparisonProduct.UUID] = *comparisonProduct
				comparisonProductsDTO = append(comparisonProductsDTO, comparisonProduct)
			}

			if categoryChars.UUID != "" {
				for _, val := range comparisonProductsDTO {
					if val.UUID == comparisonProduct.UUID {
						val.Characteristics = append(val.Characteristics, categoryChars)
						break
					}
				}
			}
		}

		// Сортируем категории и продукты
		for _, product := range comparisonProductsDTO {
			//product.ImageMediaRoot(r.mediaRoot)
			productsByCategory[product.CategoryUUID] = append(productsByCategory[product.CategoryUUID], *product)
		}

		close(errChan2)
	}()

	if err := <-errChan1; err != nil {
		return nil, nil, err
	}

	if err := <-errChan2; err != nil {
		return nil, nil, err
	}

	wg.Wait()

	// Итерируемся по слайсу сравнения
	for _, valComparison := range comparisonDTO {
		// Добавляем продукты данной категории в valComparison.ComparisonProduct
		valComparison.ComparisonProduct = append(valComparison.ComparisonProduct, productsByCategory[valComparison.CategoryUUID]...)

		if len(valComparison.Characteristics) != 0 {
			wg.Add(1)
			go func(valComparison *dto.Comparison) {
				defer wg.Done()

				lwg := sync.WaitGroup{}

				// Итерируемся по слайсу характеристик данной категории
				// Оптимальное кол-во го рутин
				for _, valChar := range categoryChar[valComparison.CategoryUUID] {
					var prodValues []dto.ProductValue

					lwg.Add(1)

					go func() {
						defer lwg.Done()
						// Итерируемся по слайсу продуктов данной категории
						for _, valProduct := range valComparison.ComparisonProduct {
							found := false
							// Итерируемся по характеристикам продукта данной категории
						prodChars:
							for _, productChar := range valProduct.Characteristics {
								// Если у продукта есть характеристика, которая есть у категории, добавляем ее значение в prodValues
								if valChar.UUID == productChar.UUID {
									prodValue := dto.ProductValue{
										ProductUUID: valProduct.UUID,
										Value:       productChar.Value,
									}
									prodValues = append(prodValues, prodValue)
									found = true    // Меняем флаг
									break prodChars // Выходим из цикла
								}
							}
							// Если у продукта нет характеристики, которая есть у категории, добавляем "-" в prodValues
							if !found {
								prodValue := dto.ProductValue{
									ProductUUID: valProduct.UUID,
									Value:       "-",
								}
								prodValues = append(prodValues, prodValue)
							}
						}
					}()
					lwg.Wait()
					// Добавляем параметр категории в valComparison.Params
					valComparison.Params = append(valComparison.Params, dto.Param{
						Title:   valChar.Title,
						Product: prodValues,
					})
				}
			}(valComparison)
		}
	}
	wg.Wait()
	return comparisonDTO, &dto.Count{Count: strconv.Itoa(len(comparisonProductsDTO))}, nil
}

func (r *ComparisonRepository) UpdateUserComparison(ctx context.Context, userId int, product *dto.UpdateComparison) error {
	var comparison []models.Comparison
	var comparisonProduct *models.ComparisonProduct

	if result := r.db.WithContext(ctx).Raw(fmt.Sprintf(`select * from comparison_product
	inner join comparisonm2ms c2m on comparison_product.id = c2m.comparison_product_id
	inner join comparison c on c.id = c2m.comparison_id
	inner join category c2 on c2.uuid = comparison_product.category_uuid
	where comparison_product.product_uuid = '%s' and c.user_id = %d limit 1`, product.ProductUUID, userId)).Scan(&comparisonProduct); result.RowsAffected != 0 {
		return fmt.Errorf("product already in comparison")
	} else {
		r.db.WithContext(ctx).Where("user_id = ?", userId).Find(&comparison)
		var mainCategoryUUID string
		r.db.WithContext(ctx).Raw(fmt.Sprintf(`select c.uuid from category_product
	    inner join category c on c.uuid = category_product.category_uuid
	    where category_product.product_uuid = '%s' and c.level = 0`, product.ProductUUID)).Scan(&mainCategoryUUID)
		prodUUID, err := uuid.Parse(product.ProductUUID)
		if err != nil {
			return err
		}
		categoryUUID, err := uuid.Parse(mainCategoryUUID)
		if err != nil {
			return err
		}
		newComparisonProduct := &models.ComparisonProduct{
			ProductUUID:  prodUUID,
			CategoryUUID: categoryUUID,
			Comparison:   comparison,
		}
		r.db.Create(&newComparisonProduct)
		return nil
	}
}

func (r *ComparisonRepository) DeleteComparisonProduct(ctx context.Context, userId int, productId string) error {
	if result := r.db.WithContext(ctx).Exec(fmt.Sprintf(`delete from comparison_product
			where id = (select comparison_product_id from comparisonm2ms
	       inner join comparison c on c.id = comparisonm2ms.comparison_id
	       inner join "user" u on u.id = c.user_id
	       inner join comparison_product cp on comparisonm2ms.comparison_product_id = cp.id
	       where cp.product_uuid = '%s' and u.id = %d);`, productId, userId)); result.RowsAffected == 0 {
		return fmt.Errorf("product not found")
	}
	return nil

}

func (r *ComparisonRepository) DeleteComparisonProductByCategoryUUID(ctx context.Context, userId int, categoryUUID string) error {
	if result := r.db.WithContext(ctx).Exec(fmt.Sprintf(`delete from comparison_product
		where id in (select cp.id from comparison inner join comparisonm2ms c2m on comparison.id = c2m.comparison_id
		inner join comparison_product cp on cp.id = c2m.comparison_product_id
		where user_id = %d and comparison_product.category_uuid = '%s')`, userId, categoryUUID)); result.RowsAffected == 0 {
		return fmt.Errorf("products not found")
	}
	return nil
}

func (r *ComparisonRepository) ClearUserComparison(ctx context.Context, userId int) error {
	if result := r.db.WithContext(ctx).Exec(fmt.Sprintf(`delete from comparison_product
			where id in (select comparison_product_id from comparisonm2ms
	       inner join comparison c on c.id = comparisonm2ms.comparison_id
	       where c.user_id = %d);`, userId)); result.RowsAffected == 0 {
		return fmt.Errorf("favourite list already empty")
	}
	return nil
}

package repository

import (
	"clean_arch/internal/dto"
	"clean_arch/internal/models"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

const zeroUUID = "00000000-0000-0000-0000-000000000000"

type ProductRepository struct {
	db         *gorm.DB
	log        *logrus.Logger
	locTime    *time.Location
	timeFormat string
	mediaRoot  string
}

func NewProductRepository(db *gorm.DB, log *logrus.Logger, locTime *time.Location, timeFormat, mediaRoot string) *ProductRepository {
	return &ProductRepository{
		db:         db,
		log:        log,
		locTime:    locTime,
		timeFormat: timeFormat,
		mediaRoot:  mediaRoot,
	}
}

func (r *ProductRepository) GetAllProductsByParams(ctx context.Context, params *dto.Params) (*dto.Products, error) {
	var productsDTO []*dto.Product
	var totalCount *dto.Count
	var sortParams *dto.SortParams
	var categoryUUIDRequest, limitRequest, productsUUID, emptyQuantity string

	if len(params.ProductUUID) != 0 {
		productsUUID = "p.uuid in ? and "
	} else {
		productsUUID = "p.uuid not in ? and"
		params.ProductUUID = []string{zeroUUID}
	}
	if len(params.CatId) == 0 {
		categoryUUIDRequest = ""
	} else {
		categoryUUIDRequest = fmt.Sprintf("and cp.category_uuid = %s", "'"+params.CatId+"'")
	}
	if params.Limit == 0 {
		limitRequest = "limit all"
	} else {
		limitRequest = fmt.Sprintf("limit %d", params.Limit)
	}
	if params.NotNull == "false" {
		emptyQuantity = fmt.Sprintf("")
	} else {
		emptyQuantity = fmt.Sprintf("and p.quantity > 0")
	}

	mainSort, secondarySort := r.CustomSort(params.Sort)

	errChan := make(chan error)

	wg := &sync.WaitGroup{}

	wg.Add(3)

	go func() {
		defer wg.Done()
		result := r.db.WithContext(ctx).Debug().Raw(fmt.Sprintf(`
		select * from (
		select 
			p.uuid, 
			p.title, 
			p.vendor_code, 
			p.base_unit, 
			'1' as count, 
			p.image, 
			p.price,
        	count(distinct r.id) as review_count, 
			p.created_at,
			p.quantity, 
			ps.sales_count, 
			ps.rating from product p
				left join review r on p.uuid = r.product_uuid
				left join category_product cp on p.uuid = cp.product_uuid
				left join product_statistic ps on p.uuid = ps.product_uuid
					where p.can_to_view = true %s and %s (p.price between %f and %f) and (ps.rating between %f and %f) %s 
						group by p.uuid, ps.sales_count, ps.rating %s
		) sub_query %s %s offset %d;`, emptyQuantity, productsUUID, params.PriceMin, params.PriceMax, params.RatingMin, params.RatingMax, categoryUUIDRequest, mainSort, secondarySort, limitRequest, (params.Page-1)*params.Limit), params.ProductUUID).Scan(&productsDTO)
		if result.RowsAffected == 0 {
			errChan <- fmt.Errorf("empty products slice")
		}
		errChan <- nil
	}()

	// Получаем макс цену у товаров
	go func() {
		defer wg.Done()
		if result := r.db.WithContext(ctx).Raw(fmt.Sprintf(`
			select 
			p.price as max_price 
			from product p
				left join category_product cp on p.uuid = cp.product_uuid
				left join product_statistic ps on p.uuid = ps.product_uuid
					where p.can_to_view = true %s
						order by p.price desc limit 1;`, categoryUUIDRequest)).Scan(&sortParams); result.RowsAffected == 0 {
			errChan <- fmt.Errorf("can't find max_price in product db")
		}
		errChan <- nil
	}()

	// Рассчитываем кол-во записей для пагинации в отедлньом запросе и потоке
	go func() {
		defer wg.Done()
		result := r.db.WithContext(ctx).Debug().Raw(fmt.Sprintf(`
		select count(distinct p.uuid) from product p
		left join review r on p.uuid = r.product_uuid
		left join product_statistic ps on p.uuid = ps.product_uuid
		left join category_product cp on p.uuid = cp.product_uuid
		where p.can_to_view = true %s and %s (p.price between %f and %f) and (ps.rating between %f and %f) %s`, emptyQuantity, productsUUID, params.PriceMin, params.PriceMax, params.RatingMin, params.RatingMax, categoryUUIDRequest), params.ProductUUID).Scan(&totalCount)
		if result.RowsAffected == 0 {
			errChan <- fmt.Errorf("empty products slice")
		}
		errChan <- nil
	}()

	for i := 0; i < 3; i++ {
		if err := <-errChan; err != nil {
			return nil, err
		}
	}

	wg.Wait()

	//for _, product := range productsDTO {
	//	wg.Add(1)
	//	go func(p *dto.Product) {
	//		defer wg.Done()
	//		p.ImageMediaRoot(r.mediaRoot)
	//	}(product)
	//}
	//wg.Wait()

	productsResultDTO := &dto.Products{
		Product:    productsDTO,
		Count:      totalCount,
		SortParams: sortParams,
	}

	return productsResultDTO, nil
}

func (r *ProductRepository) GetProductDetail(ctx context.Context, uuid string) (*dto.ProductInformation, error) {
	var product *dto.ProductInformation

	var wg sync.WaitGroup

	wg.Add(2)

	errChan1 := make(chan error)
	errChan2 := make(chan error)

	go func() {
		defer wg.Done()
		r.db.WithContext(ctx).Raw(fmt.Sprintf(`
			select
    		p.uuid, '1' as count, p.title, p.description,
    		p.vendor_code, p.image, p.price, p.base_unit, ps.rating,
    		p.quantity from product p
    		    left join product_statistic ps on p.uuid = ps.product_uuid
    		        where p.uuid = '%s'
    		            group by p.uuid, ps.rating;`, uuid)).Scan(&product)

		r.db.WithContext(ctx).Raw(fmt.Sprintf(`
			select product_files.* from product_files 
				inner join product_file pf on product_files.id = pf.product_files_id 
					where pf.product_uuid = '%s';`, uuid)).Scan(&product.Files)

		r.db.WithContext(ctx).Raw(fmt.Sprintf(`
			select characteristic.*, pc.value from characteristic 
				inner join product_characteristic pc on characteristic.uuid = pc.characteristic_uuid 
					where pc.product_uuid = '%s';`, uuid)).Scan(&product.Characteristics)

		r.db.WithContext(ctx).Raw(fmt.Sprintf(`
			select 
			uuid, 
			title from category
    			inner join category_product cp on category.uuid = cp.category_uuid
    				where cp.product_uuid = '%s'
						order by level desc limit 1;`, product.UUID)).Scan(&product.Category)
		//product.ImageMediaRoot(r.mediaRoot)

		for _, x := range product.Files {
			x.DocumentRoot(r.mediaRoot)
		}
		close(errChan1)
	}()

	go func() {
		defer wg.Done()
		if result := r.db.WithContext(ctx).Exec(fmt.Sprintf(`update product_statistic
		set request_detail_count = request_detail_count + 1
		where product_uuid = '%s'`, uuid)); result.RowsAffected == 0 {
			errChan2 <- fmt.Errorf("product statistic doesn't updated")
		}
		close(errChan2)
	}()

	if err := <-errChan1; err != nil {
		return nil, err
	}

	if err := <-errChan2; err != nil {
		return nil, err
	}

	wg.Wait()

	return product, nil
}

func (r *ProductRepository) CreateReview(ctx context.Context, reviewDTO *dto.Review, userId int) error {
	var productDB *dto.ProductReviewDB
	var reviewCount int

	tx := r.db.Begin().WithContext(ctx)

	errChan := make(chan error)
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		r.db.WithContext(ctx).Debug().Raw(fmt.Sprintf(`
			select distinct 
			r.id as review_id, 
			ps.rating, 
			ps.sales_count from product
				inner join cart_product cp on product.uuid = cp.product_uuid
				inner join cartm2ms c2m on cp.id = c2m.cart_product_id
				inner join cart c on c.id = c2m.cart_id
				left join product_statistic ps on product.uuid = ps.product_uuid
				left join review r on product.uuid = r.product_uuid and r.user_id = c.user_id
	    			where c.in_order = true and c.user_id = %d and cp.product_uuid = '%s';`, userId, reviewDTO.ProductUUID)).Scan(&productDB)
		if productDB.ReviewID != 0 {
			errChan <- fmt.Errorf("review already created to this product")
		}
		errChan <- nil
	}()

	go func() {
		defer wg.Done()
		if result := r.db.WithContext(ctx).Raw(fmt.Sprintf(`
			select
			count(review.id) as review_count from product
				inner join review on product.uuid = review.product_uuid
				left join product_statistic ps on product.uuid = ps.product_uuid
					where uuid = '%s'`, reviewDTO.ProductUUID)).Scan(&reviewCount); result.RowsAffected == 0 {
			errChan <- fmt.Errorf("error occured db")
		}
		errChan <- nil
	}()

	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			return err
		}
	}

	wg.Wait()

	newRatingByUser, err := strconv.ParseFloat(reviewDTO.Rating, 64)

	if err != nil {
		return err
	}

	dbRating := math.Round((productDB.Rating*float64(reviewCount) + newRatingByUser) / float64(reviewCount+1))

	tx.Exec(fmt.Sprintf(`update product_statistic
	set rating = %v where product_statistic.product_uuid = '%s'`, dbRating, reviewDTO.ProductUUID))

	reviewRating, err := strconv.Atoi(reviewDTO.Rating)

	if err != nil {
		tx.Rollback()
		return err
	}

	prodUUID, err := uuid.Parse(reviewDTO.ProductUUID)
	if err != nil {
		return err
	}

	var reviewImages []models.ReviewPhotos

	for _, val := range reviewDTO.Image {
		imageBytes, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			return err
		}
		reviewImages = append(reviewImages, models.ReviewPhotos{Image: imageBytes})
	}

	newReview := models.Review{
		Body:         reviewDTO.Body,
		Rating:       reviewRating,
		ProductUUID:  prodUUID,
		UserID:       userId,
		Recommend:    reviewDTO.Recommend,
		CreatedAt:    time.Now().In(r.locTime),
		ReviewPhotos: reviewImages,
	}
	tx.Create(&newReview)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *ProductRepository) GetReviews(ctx context.Context, productUUID string, params *dto.Params) (*dto.ProductStatistic, error) {
	var productDB *dto.ProductDB
	var limitRequest string

	if params.Limit == 0 {
		limitRequest = "limit all"
	} else {
		limitRequest = fmt.Sprintf("limit %d", params.Limit)
	}
	if result := r.db.WithContext(ctx).Raw(fmt.Sprintf(`select
    uuid,
    ps.rating,
    ps.sales_count,
    count(review.id) as review_count,
    (select count(id) from review where rating = 5 and review.product_uuid = uuid) as sum5,
    (select count(id) from review where rating = 4 and review.product_uuid = uuid) as sum4,
    (select count(id) from review where rating = 3 and review.product_uuid = uuid) as sum3,
    (select count(id) from review where rating = 2 and review.product_uuid = uuid) as sum2,
    (select count(id) from review where rating = 1 and review.product_uuid = uuid) as sum1,
    (case
        when (select count(id) from review where recommend = true and review.product_uuid = uuid) = 0 then 0
        else (select count(id) from review where recommend = true and review.product_uuid = uuid) / (count(review.id)) * 100
    end) as recommend_percent
	from product
         inner join review on product.uuid = review.product_uuid
         left join product_statistic ps on product.uuid = ps.product_uuid
	where uuid = '%s'
	group by uuid, ps.rating, ps.sales_count
	order by uuid
	limit 1`, productUUID)).
		Scan(&productDB); result.RowsAffected == 0 {
		return nil, nil
	}

	if result := r.db.WithContext(ctx).Raw(fmt.Sprintf(`
		select 
		review.*, 
		u.name
		from review
			inner join "user" u on review.user_id = u.id 
				where review.product_uuid = '%s' 
					order by review.rating %s offset %d`, productUUID, limitRequest,
		(params.Page-1)*params.Limit)).
		Scan(&productDB.Reviews); result.RowsAffected == 0 {
		return nil, nil
	}

	for i, val := range productDB.Reviews {
		r.db.WithContext(ctx).Raw(fmt.Sprintf(`select image from review_photos where review_id = %d`, val.Id)).Scan(&productDB.Reviews[i].Image)
	}

	reviewsDTO := make([]dto.GetReview, len(productDB.Reviews), cap(productDB.Reviews))

	for i, val := range productDB.Reviews {
		reviewDTO := dto.GetReview{
			Id:        fmt.Sprint(val.Id),
			Body:      val.Body,
			Rating:    fmt.Sprint(val.Rating),
			Name:      val.Name,
			CreatedAt: fmt.Sprint(val.CreatedAt.Format(r.timeFormat)),
		}
		for _, image := range val.Image {
			// декодируем фото
			dst := make([]byte, base64.StdEncoding.EncodedLen(len(image.Image)))
			base64.StdEncoding.Encode(dst, image.Image)
			reviewDTO.Image = append(reviewDTO.Image, string(dst))
		}
		reviewsDTO[i] = reviewDTO
	}
	productStatistic := &dto.ProductStatistic{
		Recommend:   productDB.RecommendPercent,
		Rating:      productDB.Rating,
		ReviewCount: productDB.ReviewCount,
		SalesCount:  productDB.SalesCount,
		StarsStatistic: dto.Stars{
			Star5: productDB.Sum5,
			Star4: productDB.Sum4,
			Star3: productDB.Sum3,
			Star2: productDB.Sum2,
			Star1: productDB.Sum1,
		},
		Reviews: reviewsDTO,
	}
	return productStatistic, nil
}

func (r *ProductRepository) CustomSort(strArr []string) (string, string) {
	// Учитываем что строка не имутабельна
	mainSortOptions := strings.Builder{}
	secondarySortOptions := strings.Builder{}

	sortParamsMap := map[string]map[string]string{
		"popular":      {"main": "ps.sales_count desc", "secondary": "sub_query.sales_count desc"},
		"lower_price":  {"main": "p.price asc", "secondary": "sub_query.price asc"},
		"higher_price": {"main": "p.price desc", "secondary": "sub_query.price desc"},
		"news":         {"main": "p.created_at desc", "secondary": "sub_query.created_at desc"},
		"default":      {"main": "ps.sales_count desc, p.uuid desc", "secondary": "sub_query.sales_count desc, subquery.uuid desc"},
		"discounts":    {},
	}

	if len(strArr) == 0 {
		mainSortOptions.WriteString(fmt.Sprintf(" order by %v", sortParamsMap["default"]["main"]))
	} else {
		mainSortOptions.WriteString(fmt.Sprintf(" order by %v", sortParamsMap[strArr[0]]["main"]))
		if len(strArr) > 1 {
			for i, key := range strArr[1:] {
				var stringToAdd string

				val, ok := sortParamsMap[key]["secondary"]

				if ok {
					if i == 0 {
						stringToAdd = fmt.Sprintf(" order by %v", val)
					} else {
						stringToAdd = fmt.Sprintf(", %v", val)
					}
				}
				secondarySortOptions.WriteString(stringToAdd)
			}
		}
	}
	return mainSortOptions.String(), secondarySortOptions.String()
}

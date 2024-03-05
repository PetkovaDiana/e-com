package repository

import (
	"clean_arch/internal/dto"
	"clean_arch/internal/models"
	"clean_arch/pkg"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math"
	"strconv"
	"sync"
	"time"
)

type OrderRepository struct {
	db         *gorm.DB
	log        *logrus.Logger
	locTime    *time.Location
	timeFormat string
	mediaRoot  string
}

func NewOrderRepository(db *gorm.DB, log *logrus.Logger, locTime *time.Location, timeFormat, mediaRoot string) *OrderRepository {
	return &OrderRepository{
		db:         db,
		log:        log,
		locTime:    locTime,
		timeFormat: timeFormat,
		mediaRoot:  mediaRoot,
	}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, userId int, orderInfo *dto.Order) (*dto.OrderData, error) {

	orderData := &dto.OrderData{}

	cartProducts := []*dto.CartProductsDB{}
	var receiptInfo *dto.Receipt

	var orderProductsInfo []*dto.OrderProduct

	userOrderData := dto.UserOrderData{}

	tx := r.db.WithContext(ctx).Begin()

	errChan := make(chan error)

	var wg sync.WaitGroup
	wg.Add(4)

	// Увеличиваем кол-во продаж у товаров
	go func() {
		defer wg.Done()
		defer func() { errChan <- nil }()
		localTx := r.db.Session(&gorm.Session{Context: ctx})
		if err := r.UpdateProductSalesCountStatistic(localTx, userId, &cartProducts); err != nil {
			localTx.Rollback()
			errChan <- err
		}
	}()

	// Меняем кол-во на складе
	go func() {
		defer wg.Done()
		defer func() { errChan <- nil }()
		localTx := r.db.Session(&gorm.Session{Context: ctx})
		if err := r.UpdateProductQuantity(localTx, userId); err != nil {
			localTx.Rollback()
			errChan <- err
		}
	}()

	// Получаем чек
	go func() {
		defer wg.Done()
		defer func() { errChan <- nil }()
		var err error

		var deliveryType int

		if orderInfo.DeliveryType.CourierDelivery.Address != "" {
			deliveryType = 2
		} else if orderInfo.DeliveryType.SelfDelivery.PickUpPointsID != 0 {
			deliveryType = 1
		} else {
			deliveryType = 3
		}

		receiptInfo, err = r.OrderReceipt(ctx, userId, &dto.GetReceipt{
			PromoCode:    orderInfo.PromoCode,
			DeliveryType: strconv.Itoa(deliveryType),
		})

		if err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		defer func() { errChan <- nil }()
		localTx := r.db.Session(&gorm.Session{Context: ctx})
		if err := r.UserOrderData(localTx, userId, &userOrderData, orderInfo.PaymentMethodID); err != nil {
			errChan <- err
		}
	}()

	for i := 0; i <= 3; i++ {
		if err := <-errChan; err != nil {
			return nil, err
		}
	}

	wg.Wait()

	totalPrice, _ := strconv.ParseFloat(receiptInfo.FinalPrice, 64)

	newOrder := models.Order{
		UserID:          userId,
		CreatedAt:       time.Now().In(r.locTime),
		PaymentMethodID: orderInfo.PaymentMethodID,
		DeliveryType: models.DeliveryType{
			CourierDelivery: models.CourierDelivery{
				Address:         orderInfo.DeliveryType.CourierDelivery.Address,
				ApartmentOffice: orderInfo.DeliveryType.CourierDelivery.ApartmentOffice,
				Index:           orderInfo.DeliveryType.CourierDelivery.Index,
				Entrance:        orderInfo.DeliveryType.CourierDelivery.Entrance,
				Intercom:        orderInfo.DeliveryType.CourierDelivery.Intercom,
				Floor:           orderInfo.DeliveryType.CourierDelivery.Floor,
				Note:            orderInfo.DeliveryType.CourierDelivery.Note,
			},
			SelfDelivery: models.SelfDelivery{
				PickUpPointID: orderInfo.DeliveryType.SelfDelivery.PickUpPointsID,
			},
			CDEKDelivery: models.CDEKDelivery{
				PickUpPointAddress: orderInfo.DeliveryType.CDEKDelivery.PickUpPointAddress,
			},
		},
		CartID:     cartProducts[0].CartID,
		TotalPrice: totalPrice,
	}

	if orderInfo.PaymentMethodID != 2 {
		newOrder.OrderStatusID = 4
	}

	wg.Add(5)

	// Уменьшаем кол-во использований промокода
	go func() {
		defer wg.Done()
		defer func() { errChan <- nil }()
		if orderInfo.PromoCode != "" {
			newOrder.Promo = true
			localTx := r.db.Session(&gorm.Session{Context: ctx})
			if err := r.ReduceNumberOfPromoCodeUses(localTx, orderInfo.PromoCode); err != nil {
				errChan <- err
			}
		}
	}()

	// Помечаем корзину как оформленную в заказе
	go func() {
		defer wg.Done()
		defer func() { errChan <- nil }()
		localTx := r.db.Session(&gorm.Session{Context: ctx})
		localTx.Exec(fmt.Sprintf(`update cart
			set in_order = true
			where id = %d`, cartProducts[0].CartID))
	}()

	// Создаем новый заказ
	go func() {
		defer wg.Done()
		defer func() { errChan <- nil }()
		localTx := r.db.Session(&gorm.Session{Context: ctx})
		if err := localTx.Create(&newOrder).Error; err != nil {
			tx.Rollback()
			errChan <- err
		}
	}()

	// Создаем новую корзину
	go func() {
		defer wg.Done()
		defer func() { errChan <- nil }()
		localTx := r.db.Session(&gorm.Session{Context: ctx})
		if err := localTx.Create(&models.Cart{
			UserID: userId,
		}).Error; err != nil {
			tx.Rollback()
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		defer func() { errChan <- nil }()
		localTx := r.db.Session(&gorm.Session{Context: ctx})
		if result := localTx.Raw(fmt.Sprintf(`
		select
		p.uuid,
    	p.title,
    	p.price,
		p.vendor_code as article,
		p.base_unit as unit,
    	cart_product.count as quantity,
    	cart_product.total_price
    	from cart_product
    		inner join cartm2ms c2m on cart_product.id = c2m.cart_product_id
    		inner join cart c on c.id = c2m.cart_id
    		inner join product p on p.uuid = cart_product.product_uuid
				where cart_id = %d
					group by p.uuid, p.title, p.price, cart_product.count, cart_product.total_price;`, cartProducts[0].CartID)).Scan(&orderProductsInfo); result.RowsAffected == 0 {
			tx.Rollback()
			errChan <- fmt.Errorf("error ocсured db")
		}
	}()

	for i := 0; i < 5; i++ {
		if err := <-errChan; err != nil {
			return nil, err
		}
	}

	wg.Wait()

	smsOrderDTO := &dto.SMSOrder{
		Product:    orderProductsInfo,
		TotalPrice: fmt.Sprint(totalPrice),
		Phone:      userOrderData.Phone,
		OrderID:    strconv.Itoa(newOrder.Id),
	}

	var emailStatic *models.EmailStatic

	r.db.WithContext(ctx).First(&emailStatic)
	emailStaticDTO := &dto.EmailStatic{
		CarImage:     r.mediaRoot + emailStatic.CarImage,
		CartImage:    r.mediaRoot + emailStatic.CartImage,
		LikeImage:    r.mediaRoot + emailStatic.LikeImage,
		LogoImage:    r.mediaRoot + emailStatic.LogoImage,
		CourierEmail: emailStatic.CourierEmail,
	}

	emailOrderDTO := &dto.EmailOrder{
		Product:         orderProductsInfo,
		CourierDelivery: orderInfo.DeliveryType.CourierDelivery,
		EmailStatic:     emailStaticDTO,
		TotalPrice:      fmt.Sprint(totalPrice),
		Phone:           userOrderData.Phone,
		OrderID:         strconv.Itoa(newOrder.Id),
		PaymentMethod:   userOrderData.PaymentMethod,
		Inn:             userOrderData.Inn,
		Email:           userOrderData.Email,
		DeliveryPrice:   receiptInfo.DeliveryPrice,
		Sale:            receiptInfo.Sale,
		FIO:             userOrderData.Name,
		ManagerName:     userOrderData.ManagerName,
	}

	var selfDeliveryTitle string
	r.db.WithContext(ctx).Raw(fmt.Sprintf(`select address from pick_up_point where id = %d`, orderInfo.DeliveryType.SelfDelivery.PickUpPointsID)).Scan(&selfDeliveryTitle)

	if len(orderInfo.DeliveryType.CourierDelivery.Address) != 0 {
		emailOrderDTO.Address = fmt.Sprintf("Доставка по адресу: "+"%s", orderInfo.DeliveryType.CourierDelivery.Address)
	} else if orderInfo.DeliveryType.SelfDelivery.PickUpPointsID != 0 {
		emailOrderDTO.Address = fmt.Sprintf("Самовывоз по адресу: "+"%s", selfDeliveryTitle)
	} else {
		emailOrderDTO.Address = fmt.Sprintf("Точка CDEK по адресу: "+"%s", orderInfo.DeliveryType.CDEKDelivery.PickUpPointAddress)
	}

	invoiceData := &dto.InvoiceData{}
	order1C := &dto.Order1C{}
	var productsUUID []string
	// TODO запускаем через го рутины
	// Оформляем чек для счета если способ оплаты == 3

	wg.Add(3)

	go func() {
		defer wg.Done()
		if orderInfo.PaymentMethodID == 3 {
			invoiceData = r.CreateInvoice(orderInfo.DeliveryType.CourierDelivery.Address != "", orderProductsInfo, totalPrice, &userOrderData, &newOrder)
		}
	}()

	// Заполняем uuid продуктов
	// ....

	go func() {
		defer wg.Done()
		productsUUID = r.CreateSearchServiceData(orderProductsInfo)
	}()

	// Заполняем заказ 1С
	// ....

	go func() {
		defer wg.Done()
		order1C = r.Create1COrder(orderProductsInfo, &userOrderData, selfDeliveryTitle)
	}()

	// Заполняем Email
	// ....

	// Заполняем Phone
	// ....

	wg.Wait()
	//Order 1C Id of filial
	order1C.Data.Branch = orderInfo.DeliveryType.SelfDelivery.PickUpPointsID
	order1C.Data.Id = newOrder.Id
	order1C.Data.CreatedAt = newOrder.CreatedAt

	orderData.EmailOrder = *emailOrderDTO
	orderData.SMSOrder = *smsOrderDTO
	orderData.Order1C = *order1C
	orderData.ProductUUIDs = productsUUID
	orderData.PaymentID = userOrderData.PaymentID
	orderData.NewOrderID = newOrder.Id
	orderData.InvoiceData = *invoiceData

	orderData.ReceiptData.ClientEmail = userOrderData.Email
	orderData.ReceiptData.ClientPhone = userOrderData.Phone
	orderData.ReceiptData.ClientID = userOrderData.ID
	orderData.ReceiptData.PayAmount = newOrder.TotalPrice

	tx.Commit()
	return orderData, nil
}

func (r *OrderRepository) GetUserOrders(ctx context.Context, userId int, params *dto.Params) ([]*dto.GetOrder, *dto.Count, error) {
	var userOrders []*dto.GetOrder

	//TODO optimize it

	var limit string

	if params.Limit == 0 {
		limit = "limit all"
	} else {
		limit = fmt.Sprintf("limit %d", params.Limit)
	}

	if result := r.db.WithContext(ctx).Raw(fmt.Sprintf(`select "order".id as id, "order".created_at as created_at_db, "order".cancel, "order".total_price, os.name as order_status, count(*) over() as total_orders from "order"
	join cart c on "order".cart_id = c.id inner join cartm2ms c2m on c.id = c2m.cart_id inner join cart_product cp on cp.id = c2m.cart_product_id inner join order_status os on os.id = "order".order_status_id 
	where "order".user_id = %d group by "order".id, os.name order by "order".id desc %s offset %d`, userId, limit, (params.Page-1)*params.Limit)).
		Scan(&userOrders); result.RowsAffected == 0 {
		return nil, nil, nil
	}

	for x := range userOrders {
		if result := r.db.WithContext(ctx).Raw(fmt.Sprintf(`select cp.product_uuid as uuid, cp.count as quantity, p.title, p.image, p.price, cp.total_price, cp.count, r.rating, '1' as count
	from "order" join cart c on "order".cart_id = c.id inner join cartm2ms c2m on c.id = c2m.cart_id inner join cart_product cp on cp.id = c2m.cart_product_id inner join order_status os on os.id = "order".order_status_id inner join product p on p.uuid = cp.product_uuid left 
	join review r on p.uuid = r.product_uuid where "order".user_id = %d and "order".id = %s group by "order".id, cp.product_uuid, p.title, p.image, p.price, cp.total_price, cp.count, r.rating`, userId, userOrders[x].Id)).
			Scan(&userOrders[x].Products); result.RowsAffected == 0 {
			return nil, nil, fmt.Errorf("no orders yet")
		}
	}

	var wg sync.WaitGroup

	for _, order := range userOrders {
		wg.Add(1)
		go func(p *dto.GetOrder) {
			defer wg.Done()
			p.TimeFormatter(r.timeFormat)
			//for _, product := range p.Products {
			//	wg.Add(1)
			//	go func(p *dto.OrderProduct) {
			//		defer wg.Done()
			//		p.ImageMediaRoot(r.mediaRoot)
			//	}(product)
			//}
		}(order)
	}
	wg.Wait()
	return userOrders, &dto.Count{Count: fmt.Sprint(userOrders[0].TotalOrders)}, nil
}

func (r *OrderRepository) CancelCashOrder(ctx context.Context, userId int, orderId int) error {
	if result := r.db.WithContext(ctx).Exec(fmt.Sprintf(`
			update "order"
			set cancel = true, order_status_id = 3 
				where id = %d and user_id = %d;`, orderId, userId)); result.RowsAffected == 0 {
		return fmt.Errorf("order not found")
	}
	return nil
}

// TODO опасный участок логики
func (r *OrderRepository) CancelPayKeeperOrder(userId int, orderId int) error {
	// здесь мы не проверяем условие статус кода у заказа, так как проверили это в функции PossibleCancelPayKeeperOrder
	if result := r.db.Exec(fmt.Sprintf(`
			update "order" o
			set order_status_id = 3
			from pay_keeper_info as pkf
				where pkf.order_id = o.id and 
					o.id = %d and 
					o.user_id = %d`, orderId, userId)); result.RowsAffected == 0 {
		return fmt.Errorf("order not found")
	}
	return nil
}

// Проверяем возможность отмены заказа
func (r *OrderRepository) PossibleCancelPayKeeperOrder(ctx context.Context, userId int, orderId int) (int, error) {
	var paymentId int
	tx := r.db.Begin().WithContext(ctx)
	if result := tx.Raw(fmt.Sprintf(`
			update "order" o
			set cancel = true, order_status_id = 3
			from pay_keeper_info as pkf
				where pkf.order_id = o.id and 
					o.id = %d and 
					o.user_id = %d and 
					o.order_status_id = 1 and
					o.cancel = false
						returning pkf.payment_id`, orderId, userId)).Scan(&paymentId); result.RowsAffected == 0 {
		tx.Rollback()
		return 0, fmt.Errorf("order not found")
	}
	tx.Rollback()
	return paymentId, nil
}

func (r *OrderRepository) CheckPaymentMethod(ctx context.Context, orderId int) (int, error) {
	var paymentMethodId int
	if result := r.db.WithContext(ctx).Raw(fmt.Sprintf(`select payment_method_id from "order" where id = %d;`, orderId)).Scan(&paymentMethodId); result.RowsAffected == 0 {
		return 0, fmt.Errorf("order not found")
	}
	return paymentMethodId, nil
}

func (r *OrderRepository) GetOrdersProducts(ctx context.Context, userId int, params *dto.Params) ([]*dto.UserOrdersProducts, *dto.Count, error) {
	var userProducts []*dto.UserOrdersProducts

	var limit string

	if params.Limit == 0 {
		limit = "limit all"
	} else {
		limit = fmt.Sprintf("limit %d", params.Limit)
	}

	if result := r.db.WithContext(ctx).Raw(fmt.Sprintf(`select cp.product_uuid as uuid, p.title, p.image, '1' as count, cp.total_price, "order".created_at as created_at_db, 
	(select sum(review.rating) from review where review.user_id = "order".user_id and review.product_uuid = cp.product_uuid) as rating, cp.count as quantity, cast(count(*) OVER() as varchar) as full_count 
	from "order" inner join cart c on "order".cart_id = c.id
    join cartm2ms c2m on c.id = c2m.cart_id
    join cart_product cp on cp.id = c2m.cart_product_id
    join product p on p.uuid = cp.product_uuid
    left join review r on p.uuid = r.product_uuid
    where "order".user_id = %d
    group by cp.product_uuid, p.title, p.image, cp.total_price, "order".created_at, cp.count, "order".user_id order by cp.product_uuid %s offset %d`, userId, limit, (params.Page-1)*params.Limit)).
		Scan(&userProducts); result.RowsAffected == 0 {
		return nil, nil, nil
	}

	for _, product := range userProducts {
		//p.ImageMediaRoot(r.mediaRoot)
		product.TimeFormatter(r.timeFormat)
	}

	return userProducts, &dto.Count{Count: userProducts[0].FullCount}, nil
}

func (r *OrderRepository) GetPaymentMethods(ctx context.Context) []*dto.PaymentMethod {
	var paymentMethods []*dto.PaymentMethod

	r.db.WithContext(ctx).Raw(fmt.Sprintf(`select id, title, description, payment_method.icon, payment_method.image from payment_method;`)).
		Scan(&paymentMethods)

	for _, x := range paymentMethods {
		x.ImageMediaRoot(r.mediaRoot)
	}
	return paymentMethods
}

func (r *OrderRepository) GetPickUpPoints(ctx context.Context) []*dto.PickUpPoint {
	var pickUpPointsDTOArray []*dto.PickUpPoint

	//TODO fix this, all is bad

	rows, err := r.db.WithContext(ctx).Raw(fmt.Sprintf(`select pick_up_point.id as id, pick_up_point.phone1, pick_up_point.phone2, pick_up_point.phone3, pick_up_point.email1, 
		pick_up_point.email2, pick_up_point.address, pick_up_point.coordinate_x, pick_up_point.coordinate_y, pupt.mon, pupt.tue, pupt.wen, pupt.thu,
       pupt.fri, pupt.sat, pupt.sun from pick_up_point
		inner join pick_up_point_time pupt on pupt.id = pick_up_point.pick_up_point_time_id
		order by pick_up_point.id, pick_up_point.address;`)).Rows()

	if err != nil {
		return nil
	}

	defer func() error {
		if err := rows.Close(); err != nil {
			return err
		}
		return nil
	}()
	for rows.Next() {
		var pickUpPointsDTO dto.PickUpPoint
		var pickUpPointsTimesDTO dto.PickUpPointTime
		var pickUpPointsStockTitleDTO dto.PickUpPointStockTitle
		var pickUpPointsStockDescriptionDTO dto.PickUpPointStockDescription
		var coordinatesDTO dto.Coordinates

		r.db.WithContext(ctx).ScanRows(rows, &pickUpPointsDTO)
		r.db.WithContext(ctx).ScanRows(rows, &pickUpPointsTimesDTO)
		r.db.WithContext(ctx).ScanRows(rows, &coordinatesDTO)

		pickUpPointsDTO.PickUpPointTime = pickUpPointsTimesDTO

		pickUpPointsDTO.Coordinates = coordinatesDTO

		rowss, err := r.db.WithContext(ctx).Raw(fmt.Sprintf(`select * from pick_up_point_stock_title
    	inner join pick_up_point on pick_up_point_stock_title.id = pick_up_point.pick_up_point_stock_title_id
		left join pick_up_point_stock_description pupsd on pick_up_point_stock_title.id = pupsd.pick_up_point_stock_title_id
		where pick_up_point.id = %s;`, pickUpPointsDTO.Id)).Rows()

		if err != nil {
			return nil
		}

		for rowss.Next() {
			r.db.ScanRows(rowss, &pickUpPointsStockTitleDTO)
			r.db.ScanRows(rowss, &pickUpPointsStockDescriptionDTO)

			if pickUpPointsStockTitleDTO.Title != "" {
				pickUpPointsStockTitleDTO.Description = append(pickUpPointsStockTitleDTO.Description, pickUpPointsStockDescriptionDTO)
				pickUpPointsDTO.PickUpPointStock = pickUpPointsStockTitleDTO
			}
		}
		pickUpPointsDTOArray = append(pickUpPointsDTOArray, &pickUpPointsDTO)
	}
	return pickUpPointsDTOArray
}

func (r *OrderRepository) PromoCodeValidator() {
	r.db.Exec(fmt.Sprintf(`DELETE FROM promo_code WHERE expires_at <= now();`))
}

func (r *OrderRepository) OrderReceipt(ctx context.Context, userId int, receiptInfo *dto.GetReceipt) (*dto.Receipt, error) {
	var receiptDTO *dto.Receipt

	if result := r.db.Raw(fmt.Sprintf(`
		select 
		sum(cart_product.total_price) as cart_price, 
		sum(cart_product.count) as product_count,
       		case
       		    when promo_code.discount_sum is not null then round(sum(cart_product.total_price) - promo_code.discount_sum +
					(
					case when sum(cart_product.total_price) > 10000 then 0 else dtp.delivery_price end), 2
					)
       		    when promo_code.discount_percent is not null then round(sum(cart_product.total_price) * (100-promo_code.discount_percent)/cast(100 as integer), 2) +
					(
					case when sum(cart_product.total_price) > 10000 then 0 else dtp.delivery_price end
					)
       		    else sum(cart_product.total_price) +
					(
					case when sum(cart_product.total_price) > 10000 then 0 else dtp.delivery_price end
					)
       		    end
       		                              as final_price,
       		sum(cart_product.total_price) + dtp.delivery_price - (
				case
					when promo_code.discount_sum is not null then round(sum(cart_product.total_price) -promo_code.discount_sum + dtp.delivery_price, 2)
					when promo_code.discount_percent is not null then round(sum(cart_product.total_price) * (100-promo_code.discount_percent)/cast(100 as integer), 2) + dtp.delivery_price
					else sum(cart_product.total_price) + dtp.delivery_price
       		    end
				) as sale,
       		case
       		    when sum(cart_product.total_price) >= 10000 then 0
       		    else dtp.delivery_price
       		    end as delivery_price
				from cart_product
					inner join cartm2ms c2m on cart_product.id = c2m.cart_product_id
       		  		inner join cart c on c.id = c2m.cart_id
       		  		left join promo_code on promo_code.promo_code = '%s' and promo_code.number_of_uses > 0
       		  		cross join delivery_type_info dtp
						where c.in_order = false and c.user_id = %d and dtp.id = %s
						group by dtp.delivery_price, promo_code.discount_sum, promo_code.discount_percent;`, receiptInfo.PromoCode, userId, receiptInfo.DeliveryType)).
		Scan(&receiptDTO); result.RowsAffected == 0 {
		return nil, fmt.Errorf("cart is empty")
	}
	return receiptDTO, nil
}

func (r *OrderRepository) GetDeliveryType(ctx context.Context) []*dto.DeliveryTypeInfo {
	var deliveryTypes []*dto.DeliveryTypeInfo
	if result := r.db.WithContext(ctx).Raw(fmt.Sprintf(`select delivery_type_info.id, 
	delivery_type_info.title, 
	delivery_type_info.description, 
	delivery_type_info.icon, 
	delivery_type_info.can_delivery from delivery_type_info`)).Scan(&deliveryTypes); result.RowsAffected == 0 {
		return nil
	}

	for _, x := range deliveryTypes {
		x.ImageMediaRoot(r.mediaRoot)
	}

	return deliveryTypes
}

func (r *OrderRepository) OnlinePaymentValidator(ctx context.Context, orderInfo *dto.OnlineOrderChecker) error {

	var orderSum float64

	orderId, err := strconv.Atoi(orderInfo.OrderID)
	paymentId, err := strconv.Atoi(orderInfo.ID)

	if err != nil {
		return fmt.Errorf("failed to parse order id: %w", err)
	}

	userId, err := strconv.Atoi(orderInfo.ClientID)
	if err != nil {
		return fmt.Errorf("failed to parse user id: %w", err)
	}

	rows, _ := r.db.WithContext(ctx).Raw(fmt.Sprintf(`select total_price from "order" where user_id = %d and id = %d`, userId, orderId)).Rows()

	defer func() error {
		if err = rows.Close(); err != nil {
			return err
		}
		return nil
	}()

	for rows.Next() {
		_ = rows.Scan(&orderSum)
	}

	if orderInfo.Sum != orderSum {
		return fmt.Errorf("sum is not valid")
	}

	if err = r.db.WithContext(ctx).Exec(fmt.Sprintf(`
		update "order"
		set order_status_id = 1
			where id = %d`, orderId)).Error; err != nil {
		return err
	}

	if err = r.db.WithContext(ctx).Exec(fmt.Sprintf(`insert into pay_keeper_info (order_id, payment_id) values (%d, %d);`, orderId, paymentId)).Error; err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) ChangeOrderStatus(ctx context.Context, orderId int) error {
	if err := r.db.WithContext(ctx).Exec(fmt.Sprintf(`update "order" set order_status_id = 5, cancel = true where id = %d`, orderId)); err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *OrderRepository) ChangeCancelStatus(orderId int) error {
	if err := r.db.Exec(fmt.Sprintf(`update "order" set order_status_id = 5, cancel = true where id = %d`, orderId)); err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *OrderRepository) UpdateProductSalesCountStatistic(tx *gorm.DB, userId int, cartProducts *[]*dto.CartProductsDB) error {
	if result := tx.Raw(fmt.Sprintf(`update product_statistic
			set sales_count = sales_count + sub_query.count
			from (
	    		select cart_product.product_uuid as uuid, c.id, cart_product.count from cart_product
	       		inner join cartm2ms c2m on cart_product.id = c2m.cart_product_id
	       		inner join cart c on c.id = c2m.cart_id
	    				where c.user_id = %d and c.in_order = false) as sub_query
				where product_uuid = sub_query.uuid
					returning sub_query.uuid as product_uuid, sub_query.id as cart_id;`, userId)).Scan(&cartProducts); result.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("user cart empty")
	}
	return nil
}

func (r *OrderRepository) UpdateProductQuantity(tx *gorm.DB, userId int) error {
	if result := tx.Exec(fmt.Sprintf(`update product
			set quantity = quantity - sub_query.count
			from (
         	select cart_product.product_uuid as uuid, c.id, cart_product.count from cart_product
				inner join cartm2ms c2m on cart_product.id = c2m.cart_product_id
				inner join cart c on c.id = c2m.cart_id
         	where c.user_id = %d and c.in_order = false) as sub_query
				where product.uuid = sub_query.uuid`, userId)); result.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("user cart empty")
	}
	return nil
}

func (r *OrderRepository) UserOrderData(tx *gorm.DB, userId int, userOrderData *dto.UserOrderData, paymentMethodID int) error {
	if result := tx.Raw(fmt.Sprintf(`select
		u.id,
		u.phone, 
		u.inn,
		u.kpp,
		u.name, 
		u.manager_name,
		u.company_name,
		u.company_address,
		p.id as payment_id,
		p.title as payment_method, 
		e.email 
		from "user" u
        	left join email e on e.user_id = u.id
        	cross join payment_method p
            	where u.id = %d and p.id = %d`, userId, paymentMethodID)).Scan(&userOrderData); result.RowsAffected == 0 {
		return fmt.Errorf("error ocсured db")
	}
	return nil
}

func (r *OrderRepository) ReduceNumberOfPromoCodeUses(tx *gorm.DB, promoCode string) error {
	if result := tx.Exec(fmt.Sprintf(`update promo_code
		set number_of_uses = number_of_uses - 1
		where promo_code = '%s';`, promoCode)); result.RowsAffected == 0 {
		return fmt.Errorf("promocode not found")
	}
	return nil
}

func (r *OrderRepository) CreateInvoice(deliveryOk bool, orderProductsInfo []*dto.OrderProduct, totalPrice float64, userOrderData *dto.UserOrderData, newOrder *models.Order) *dto.InvoiceData {
	invoiceData := &dto.InvoiceData{}
	var invoiceCartProducts []dto.InvoiceDataCartProduct
	var totalCartPrice float64

	// Аллоцируем память под слайс если доставка присутствует
	if deliveryOk {
		invoiceCartProducts = make([]dto.InvoiceDataCartProduct, len(orderProductsInfo)+1, cap(orderProductsInfo)+1)
		invoiceCartProducts[len(orderProductsInfo)] = dto.InvoiceDataCartProduct{Id: len(orderProductsInfo) + 1, Title: "Доставка", Count: 1}
	} else {
		// В противном случаи
		invoiceCartProducts = make([]dto.InvoiceDataCartProduct, len(orderProductsInfo), cap(orderProductsInfo))
	}

	// Заполняем продукты
	for i, val := range orderProductsInfo {
		soloPrice, _ := strconv.ParseFloat(val.Price, 64)
		totalCartPrice += soloPrice
		totalPrice, _ := strconv.ParseFloat(val.TotalPrice, 64)
		count, _ := strconv.ParseFloat(val.Quantity, 64)
		invoiceCartProducts[i] = dto.InvoiceDataCartProduct{Id: i + 1, Title: val.Title, Article: val.Article, Count: count, Unit: val.Unit, SoloPrice: soloPrice, FullPrice: totalPrice}
	}
	// Кол-во товаров = длине слайса с продуктами
	invoiceData.Cart.ProductCounter = len(orderProductsInfo)

	// Выщитываем стоимость заказа
	if deliveryOk {
		invoiceCartProducts[len(orderProductsInfo)].SoloPrice = math.Round((totalPrice-totalCartPrice)*100) / 100
		invoiceCartProducts[len(orderProductsInfo)].FullPrice = invoiceCartProducts[len(orderProductsInfo)].SoloPrice
		invoiceData.Cart.ProductCounter++
	}

	// Заполняем остальную ифнормацию
	invoiceData.Id = strconv.Itoa(newOrder.Id)
	invoiceData.CreatedAt = newOrder.CreatedAt.Format("2006-01-02")
	invoiceData.Cart.Nds = pkg.Round(totalPrice * 0.2)
	invoiceData.Customer = dto.InvoiceDataCustomer{Title: userOrderData.CompanyName, Inn: userOrderData.Inn, Kpp: userOrderData.Kpp, Address: userOrderData.CompanyAddress, Email: userOrderData.Email}
	invoiceData.Cart.TotalPrice = totalPrice
	invoiceData.Cart.CartProducts = invoiceCartProducts
	return invoiceData
}

func (r *OrderRepository) CreateSearchServiceData(orderProductsInfo []*dto.OrderProduct) []string {
	productsUUID := make([]string, len(orderProductsInfo), cap(orderProductsInfo))
	for i, _ := range orderProductsInfo {
		productsUUID[i] = orderProductsInfo[i].UUID
	}
	return productsUUID
}

func (r *OrderRepository) Create1COrder(orderProductsInfo []*dto.OrderProduct, userOrderData *dto.UserOrderData, selfDeliveryTitle string) *dto.Order1C {

	order1C := &dto.Order1C{
		Data: struct {
			Id           int                  `json:"id"`
			CreatedAt    time.Time            `json:"created_at"`
			Name         string               `json:"name"`
			Phone        string               `json:"phone"`
			Address      string               `json:"address"`
			Inn          string               `json:"inn"`
			CustomerType string               `json:"customer_type"`
			Products     []dto.Order1CProduct `json:"products"`
			Comment      string               `json:"comment"`
			Branch       int                  `json:"branch"`
		}{
			Products: make([]dto.Order1CProduct, len(orderProductsInfo), cap(orderProductsInfo)),
		},
	}

	for i, val := range orderProductsInfo {
		order1C.Data.Products[i].UUID = val.UUID
		productCount, _ := strconv.ParseFloat(val.Quantity, 64)
		order1C.Data.Products[i].Quantity = productCount
	}
	order1C.Data.CustomerType = "физ лицо"
	order1C.Data.Name = userOrderData.Name
	// TODO переделать
	order1C.Data.Comment = "тест"

	if userOrderData.Inn != "" {
		order1C.Data.Name = userOrderData.ManagerName
		order1C.Data.CustomerType = "юр лицо"
	}

	//order1C.Data.CreatedAt = time.Now()
	order1C.Data.Id = userOrderData.ID
	order1C.Data.Inn = userOrderData.Inn
	order1C.Data.Address = selfDeliveryTitle
	order1C.Data.Phone = userOrderData.Phone
	return order1C
}

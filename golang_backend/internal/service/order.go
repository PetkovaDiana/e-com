package service

import (
	"clean_arch/internal/dto"
	"clean_arch/internal/repository"
	"clean_arch/pkg/client_1c"
	"clean_arch/pkg/email"
	"clean_arch/pkg/order_service_client"
	"clean_arch/pkg/pay_keeper"
	"clean_arch/pkg/phone"
	"clean_arch/pkg/rec_service_client"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"strconv"
	"sync"
	"time"
)

const (
	timeOrderCancelChecker = 1
	attemptCount           = 5
)

type OrderService struct {
	repo               repository.Order
	apiKey             string
	emailClient        email.Email
	log                *logrus.Logger
	searchClient       rec_service_client.RecommendationServiceClient
	client1c           client_1c.Client1C
	payKeeperClient    pay_keeper.PayKeeperClient
	orderServiceClient order_service_client.OrderServiceClient
}

func NewOrderService(repo repository.Order, apiKey string, emailClient email.Email, log *logrus.Logger, searchClient rec_service_client.RecommendationServiceClient, client1c client_1c.Client1C, payKeeperClient pay_keeper.PayKeeperClient, orderServiceClient order_service_client.OrderServiceClient) *OrderService {
	return &OrderService{
		repo:               repo,
		apiKey:             apiKey,
		emailClient:        emailClient,
		log:                log,
		searchClient:       searchClient,
		client1c:           client1c,
		payKeeperClient:    payKeeperClient,
		orderServiceClient: orderServiceClient,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, userId int, securityTkn string, orderInfo *dto.Order) (*dto.ReturnOrder, error) {
	s.repo.PromoCodeValidator()

	orderData, err := s.repo.CreateOrder(ctx, userId, orderInfo)

	if err != nil {
		return nil, err
	}

	//TODO add logger info + add go routine

	wg := &sync.WaitGroup{}

	wg.Add(5)

	// Если способ оплаты - выставление счета
	// Плохо, что всегда запускаем новуюго рутину и тратим на это ресурсы
	// причем в этом нет необходимости
	go func() {
		defer wg.Done()
		if orderData.PaymentID == 3 {
			log.Println("HERE", orderData.InvoiceData)
			url, err := s.orderServiceClient.SendOrderData(ctx, &orderData.InvoiceData)
			if err == nil {
				orderData.InvoiceData.Url = url
			}
		}
	}()

	go func() {
		defer wg.Done()
		if orderInfo.DeliveryType.CourierDelivery.ApartmentOffice != "" {
			s.emailClient.SendCourierOrderEmail(&orderData.EmailOrder)
		}
	}()

	go func() {
		defer wg.Done()
		if orderData.EmailOrder.Email != "" {
			s.emailClient.SendOrderEmail(&orderData.EmailOrder)
		} else {
			phoneClient := phone.NewPhone(s.apiKey)
			phoneClient.PhoneOrder(&orderData.SMSOrder)
		}
	}()

	go func() {
		defer wg.Done()
		s.searchClient.SendOrderData(ctx, orderData.ProductUUIDs)
	}()

	go func() {
		defer wg.Done()
		s.client1c.SendOrderData(ctx, &orderData.Order1C)
	}()

	wg.Wait()

	return &dto.ReturnOrder{
		Id:  strconv.Itoa(orderData.NewOrderID),
		Url: orderData.InvoiceData.Url,
	}, nil
}

func (s *OrderService) CancelOrder(ctx context.Context, userId int, securityTkn string, orderId int) error {

	paymentMethodId, err := s.repo.CheckPaymentMethod(ctx, orderId)

	if err != nil {
		return err
	}

	// Наличными при получении или счетом
	if paymentMethodId == 2 || paymentMethodId == 3 {
		if err = s.repo.CancelCashOrder(ctx, userId, orderId); err != nil {
			return err
		}
		s.client1c.SendCancelOrderData(ctx, orderId)
		return nil
	}

	// Если не наличными при получении:

	// Делаем проверку на то, можно ли отменить заказ в бд
	paymentId, err := s.repo.PossibleCancelPayKeeperOrder(ctx, userId, orderId)

	if err != nil {
		return fmt.Errorf("can't cancel order: %s", err.Error())
	}

	// Если можно, то выполняем отмену на платформе PayKeeper
	if err = s.payKeeperClient.PaymentReverse(ctx, &dto.PayKeeperOrderCancel{
		PaymentID:     paymentId,
		Partial:       false,
		SecurityToken: securityTkn,
	}); err != nil {
		return err
	}

	// Меняем статус заказа на инициализирован возврат
	if err = s.repo.ChangeOrderStatus(ctx, orderId); err != nil {
		return err
	}
	time.AfterFunc(5*time.Second, func() {
		if err = s.CheckOrderStatusWithRetry(paymentId); err != nil {
			// ошибка инициализации отмена заказа на стороне Pay Keeper
			s.repo.ChangeCancelStatus(orderId)
		} else {
			// Окончательно отменяем заказ для пользователя
			if err = s.repo.CancelPayKeeperOrder(userId, orderId); err == nil {
				s.client1c.SendCancelOrderData(ctx, orderId)
			}
		}
	})
	return nil
}

func (s *OrderService) CheckOrderStatusWithRetry(orderID int) error {
	attempt := 0
	for {
		err := s.payKeeperClient.CheckOrderStatus(orderID)
		if err == nil {
			return nil
		}

		attempt++
		if attempt > attemptCount {
			return fmt.Errorf("can't cancel order after 5 attempts")
		}

		time.Sleep(5 * time.Second)
	}
}

func (s *OrderService) GetPaymentMethods(ctx context.Context) []*dto.PaymentMethod {
	return s.repo.GetPaymentMethods(ctx)
}

func (s *OrderService) GetPickUpPoints(ctx context.Context) []*dto.PickUpPoint {
	return s.repo.GetPickUpPoints(ctx)
}

func (s *OrderService) OrderReceipt(ctx context.Context, userId int, receiptInfo *dto.GetReceipt) (*dto.Receipt, error) {
	s.repo.PromoCodeValidator()
	return s.repo.OrderReceipt(ctx, userId, receiptInfo)
}

func (s *OrderService) GetUserOrders(ctx context.Context, userId int, params *dto.Params) ([]*dto.GetOrder, *dto.Count, error) {
	return s.repo.GetUserOrders(ctx, userId, params)
}

func (s *OrderService) GetOrdersProducts(ctx context.Context, userId int, params *dto.Params) ([]*dto.UserOrdersProducts, *dto.Count, error) {
	return s.repo.GetOrdersProducts(ctx, userId, params)
}

func (s *OrderService) GetDeliveryType(ctx context.Context) []*dto.DeliveryTypeInfo {
	return s.repo.GetDeliveryType(ctx)
}

func (s *OrderService) OnlinePaymentValidator(ctx context.Context, orderInfo *dto.OnlineOrderChecker) (string, error) {

	responseStr, err := s.payKeeperClient.SecretKeyMD5Hashing(orderInfo)

	if err != nil {
		return "", err
	}

	if err = s.repo.OnlinePaymentValidator(ctx, orderInfo); err != nil {
		return "", err
	}

	return responseStr, nil
}

func (s *OrderService) GetSecurityToken(ctx context.Context) (string, error) {
	// Получаем токен
	token, err := s.payKeeperClient.GerSecurityToken(ctx)
	if err != nil {
		return "", err
	}
	return token, nil
}

package service

import (
	"clean_arch/internal/dto"
	"clean_arch/internal/repository"
	"clean_arch/pkg/bitrix_client"
	"clean_arch/pkg/client_1c"
	"clean_arch/pkg/email"
	"clean_arch/pkg/order_service_client"
	"clean_arch/pkg/pay_keeper"
	"clean_arch/pkg/rec_service_client"
	"context"
	"github.com/sirupsen/logrus"
)

type Service struct {
	Category
	Product
	User
	Cart
	Favourite
	Comparison
	Order
	CMS
}

type Category interface {
	GetAllCategories(ctx context.Context) ([]*dto.Category, error)
	GetCategoriesById(ctx context.Context, id *dto.CategoryParams) (*dto.CategoryData, error)
}

type Product interface {
	GetAllProductsByParams(ctx context.Context, params *dto.Params) (*dto.Products, error)
	GetProductDetail(ctx context.Context, uuid string) (*dto.ProductInformation, error)
	CreateReview(ctx context.Context, reviewDTO *dto.Review, userId int) error
	GetReviews(ctx context.Context, productUUID string, params *dto.Params) (*dto.ProductStatistic, error)
}

type User interface {
	RegisterUser(ctx context.Context, userDTO *dto.RegisterUser, session string, userId int) (string, error)
	AuthenticateUser(ctx context.Context, userDTO *dto.UserAuth, session string, userId int) (string, error)
	ParseToken(accessToken string) (int, error)
	RegCodeGenerator(ctx context.Context, code *dto.CodeGenerate) error
	AuthCodeGenerator(ctx context.Context, code *dto.CodeGenerate) error
	ValidateSession(sessionKey string) (string, int, error)
	GetUserData(ctx context.Context, userId int) (*dto.UserData, error)
	UpdateEmailUser(ctx context.Context, emailInfo *dto.UpdateEmail, userId int) error
	CanToSendEmail(ctx context.Context, emailInfo *dto.CanToSendEmail, userId int) error
	UpdateManagerName(ctx context.Context, userInfo *dto.UpdateManagerName, id int) error
	CreateSiteReview(ctx context.Context, siteReview *dto.SiteReview) error
}

type Comparison interface {
	GetUserComparison(ctx context.Context, userId int) ([]*dto.Comparison, *dto.Count, error)
	UpdateUserComparison(ctx context.Context, userId int, product *dto.UpdateComparison) error
	DeleteComparisonProduct(ctx context.Context, userId int, productId string) error
	DeleteComparisonProductByCategoryUUID(ctx context.Context, userId int, categoryUUID string) error
	ClearUserComparison(ctx context.Context, userId int) error
}

type Cart interface {
	GetUserCart(ctx context.Context, userId int) ([]*dto.Cart, error)
	UpdateUserCart(ctx context.Context, userId int, product *dto.UpdateCart) error
	DeleteProduct(ctx context.Context, userId int, productUUID string) error
	ClearUserCart(ctx context.Context, userId int) error
}

type Favourite interface {
	GetUserFavourites(ctx context.Context, userId int) ([]*dto.Favourite, error)
	UpdateUserFavourites(ctx context.Context, userId int, product *dto.UpdateFavourite) error
	DeleteFavouriteProduct(ctx context.Context, userId int, productUUID string) error
	ClearUserFavourite(ctx context.Context, userId int) error
}

type Order interface {
	CreateOrder(ctx context.Context, userId int, securityTkn string, orderInfo *dto.Order) (*dto.ReturnOrder, error)
	CancelOrder(ctx context.Context, userId int, securityTkn string, orderId int) error
	GetPaymentMethods(ctx context.Context) []*dto.PaymentMethod
	GetPickUpPoints(ctx context.Context) []*dto.PickUpPoint
	OrderReceipt(ctx context.Context, userId int, receiptInfo *dto.GetReceipt) (*dto.Receipt, error)
	GetUserOrders(ctx context.Context, userId int, params *dto.Params) ([]*dto.GetOrder, *dto.Count, error)
	GetOrdersProducts(ctx context.Context, userId int, params *dto.Params) ([]*dto.UserOrdersProducts, *dto.Count, error)
	GetDeliveryType(ctx context.Context) []*dto.DeliveryTypeInfo
	OnlinePaymentValidator(ctx context.Context, orderInfo *dto.OnlineOrderChecker) (string, error)
	GetSecurityToken(ctx context.Context) (string, error)
}

type CMS interface {
	RequestCall(ctx context.Context, requestDTO *dto.RequestCall, userId int) error
	GetCourierDeliveryInfo(ctx context.Context) *dto.CourierDeliveryInfo
	GetCDEKDeliveryInfo(ctx context.Context) *dto.CDEKDeliveryInfo
	GetAllVacancies(ctx context.Context) []*dto.Vacancy
	RequestVacancy(ctx context.Context, requestInfo *dto.RequestVacancy) error
	GetRequisites(ctx context.Context) *dto.Requisites
	GetPrivacyPolicy(ctx context.Context) *dto.PrivacyPolicy
}

func NewService(repo *repository.Repository, log *logrus.Logger, tokenTTL, apiKey string, emailClient email.Email, bitrixClient bitrix_client.BitrixClient, searchClient rec_service_client.RecommendationServiceClient, client1c client_1c.Client1C, payKeeperClient pay_keeper.PayKeeperClient, orderServiceClient order_service_client.OrderServiceClient) *Service {
	return &Service{
		Category:   NewCategoryService(repo),
		Product:    NewProductService(repo),
		User:       NewUserService(repo, log, tokenTTL, apiKey),
		Cart:       NewCartService(repo),
		Favourite:  NewFavouriteService(repo),
		Comparison: NewComparisonService(repo),
		Order:      NewOrderService(repo, apiKey, emailClient, log, searchClient, client1c, payKeeperClient, orderServiceClient),
		CMS:        NewCMSService(repo, log, bitrixClient),
	}
}

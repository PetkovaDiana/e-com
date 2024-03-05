package repository

import (
	"clean_arch/internal/dto"
	"clean_arch/pkg/cache"
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type Repository struct {
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
	GetCategoriesById(ctx context.Context, params *dto.CategoryParams) (*dto.CategoryData, error)
}

type Product interface {
	GetAllProductsByParams(ctx context.Context, params *dto.Params) (*dto.Products, error)
	GetProductDetail(ctx context.Context, uuid string) (*dto.ProductInformation, error)
	CreateReview(ctx context.Context, reviewDTO *dto.Review, userId int) error
	GetReviews(ctx context.Context, productUUID string, params *dto.Params) (*dto.ProductStatistic, error)
}

type User interface {
	RegisterUser(ctx context.Context, userDTO *dto.RegisterUser, session string, userId int) (int, error)
	AuthenticateUser(ctx context.Context, userDTO *dto.UserAuth, session string, userId int) (int, error)
	RegCodeGenerator(ctx context.Context, code *dto.CodeGenerate, newCode int) error
	AuthCodeGenerator(ctx context.Context, code *dto.CodeGenerate, newCode int) error
	CheckSessionInDb(sessionKey string) (int, string, error)
	CreateSession(sessionKey string) (int, error)
	SessionValidator()
	GetUserData(ctx context.Context, userId int) (*dto.UserData, error)
	UpdateEmailUser(ctx context.Context, emailInfo *dto.UpdateEmail, userId int) error
	CanToSendEmail(ctx context.Context, emailInfo *dto.CanToSendEmail, userId int) error
	UpdateManagerName(ctx context.Context, userInfo *dto.UpdateManagerName, id int) error
	CreateSiteReview(ctx context.Context, siteReview *dto.SiteReview) error
}

type Cart interface {
	GetUserCart(ctx context.Context, userId int) ([]*dto.Cart, error)
	UpdateUserCart(ctx context.Context, userId int, product *dto.UpdateCart) error
	DeleteCartProduct(ctx context.Context, userId int, productUUID string) error
	ClearUserCart(ctx context.Context, userId int) error
}

type Favourite interface {
	GetUserFavourites(ctx context.Context, userId int) ([]*dto.Favourite, error)
	UpdateUserFavourites(ctx context.Context, userId int, product *dto.UpdateFavourite) error
	DeleteFavouriteProduct(ctx context.Context, userId int, productUUID string) error
	ClearUserFavourite(ctx context.Context, userId int) error
}

type Comparison interface {
	GetUserComparison(ctx context.Context, userId int) ([]*dto.Comparison, *dto.Count, error)
	UpdateUserComparison(ctx context.Context, userId int, product *dto.UpdateComparison) error
	DeleteComparisonProduct(ctx context.Context, userId int, productId string) error
	DeleteComparisonProductByCategoryUUID(ctx context.Context, userId int, categoryUUID string) error
	ClearUserComparison(ctx context.Context, userId int) error
}

type Order interface {
	CreateOrder(ctx context.Context, userId int, orderInfo *dto.Order) (*dto.OrderData, error)
	GetPaymentMethods(ctx context.Context) []*dto.PaymentMethod
	GetPickUpPoints(ctx context.Context) []*dto.PickUpPoint
	OrderReceipt(ctx context.Context, userId int, receiptInfo *dto.GetReceipt) (*dto.Receipt, error)
	GetUserOrders(ctx context.Context, userId int, params *dto.Params) ([]*dto.GetOrder, *dto.Count, error)
	GetOrdersProducts(ctx context.Context, userId int, params *dto.Params) ([]*dto.UserOrdersProducts, *dto.Count, error)
	CancelCashOrder(ctx context.Context, userId int, orderId int) error
	CancelPayKeeperOrder(userId int, orderId int) error
	PossibleCancelPayKeeperOrder(ctx context.Context, userId int, orderId int) (int, error)
	PromoCodeValidator()
	GetDeliveryType(ctx context.Context) []*dto.DeliveryTypeInfo
	OnlinePaymentValidator(ctx context.Context, orderInfo *dto.OnlineOrderChecker) error
	ChangeOrderStatus(ctx context.Context, orderId int) error
	ChangeCancelStatus(orderId int) error
	CheckPaymentMethod(ctx context.Context, orderId int) (int, error)
}

type CMS interface {
	RequestCall(ctx context.Context, requestDTO *dto.RequestCall, userId int) error
	GetCourierDeliveryInfo(ctx context.Context) *dto.CourierDeliveryInfo
	GetCDEKDeliveryInfo(ctx context.Context) *dto.CDEKDeliveryInfo
	GetAllVacancies(ctx context.Context) []*dto.Vacancy
	RequestVacancy(ctx context.Context, responseInfo *dto.RequestVacancy) error
	GetRequisites(ctx context.Context) *dto.Requisites
	GetPrivacyPolicy(ctx context.Context) *dto.PrivacyPolicy
}

func NewRepository(db *gorm.DB, cache cache.Cache, log *logrus.Logger, locTime *time.Location, timeFormat, sessionTTL, mediaRoot string) *Repository {
	return &Repository{
		Category:   NewCategoryRepository(db, log, mediaRoot),
		Product:    NewProductRepository(db, log, locTime, timeFormat, mediaRoot),
		User:       NewUserRepository(db, log, cache, sessionTTL),
		Cart:       NewCartRepository(db, log, mediaRoot),
		Favourite:  NewFavouriteRepository(db, log, mediaRoot),
		Comparison: NewComparisonRepository(db, log, mediaRoot),
		Order:      NewOrderRepository(db, log, locTime, timeFormat, mediaRoot),
		CMS:        NewCMSRepository(db, locTime, timeFormat),
	}
}

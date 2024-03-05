package handler

import (
	_ "clean_arch/docs"
	"clean_arch/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"net/http"
)

var (
	ErrInput = fmt.Errorf("wrong json input")
)

type HttpHandler struct {
	service *service.Service
	cookie  *http.Cookie
}

func NewHttpHandler(service *service.Service, cookie *http.Cookie) *HttpHandler {
	return &HttpHandler{
		service: service,
		cookie:  cookie,
	}
}

func (h *HttpHandler) InitRoutes() *gin.Engine {
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	router.Use(
		CORSMiddleware(),
	)

	router.GET("/api/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		api.POST("/code-reg", h.RegCodeGenerator)
		api.POST("/code-auth", h.AuthCodeGenerator)
		api.POST("/payment_validator", h.OnlinePaymentValidator)
		api.POST("/site-review", h.CreateSiteReview)
		v1 := api.Group("/v1", h.UserIdentity)
		{
			auth := v1.Group("/auth")
			{
				auth.POST("/sign-up", h.RegisterUser)
				auth.POST("/sign-in", h.AuthenticateUser)
			}

			cart := v1.Group("/cart")
			{
				cart.GET("/", h.GetUserCart)
				cart.POST("/add", h.UpdateUserCart)
				cart.POST("/del", h.DeleteProduct)
				cart.GET("/clear", h.ClearUserCart)
			}
			favourite := v1.Group("/favourite")
			{
				favourite.GET("/", h.GetUserFavourites)
				favourite.POST("/add", h.UpdateUserFavourites)
				favourite.POST("/del", h.DeleteFavouriteProduct)
				favourite.GET("/clear", h.ClearUserFavourite)

			}
			comparison := v1.Group("/comparison")
			{
				comparison.GET("/", h.GetUserComparison)
				comparison.POST("/add", h.UpdateUserComparison)
				comparison.POST("/del", h.DeleteComparisonProduct)
				comparison.POST("/del_cat", h.DeleteComparisonProductByCategory)
				comparison.GET("/clear", h.ClearUserComparison)
			}
			order := v1.Group("/order")
			{
				orderValidator := order.Group("", h.UserValidator)
				{
					orderSecurityToken := orderValidator.Group("", h.GetSecurityToken)
					{
						orderSecurityToken.POST("/cancel", h.CancelUserOrder)
					}
					orderValidator.POST("/", h.CreateOrder)
					orderValidator.GET("/", h.GetUserOrders)
					orderValidator.GET("/order_products", h.GetOrdersProducts)
				}
				order.POST("/receipt", h.OrderReceipt)
				order.GET("/delivery_types", h.GetDeliveryType)
				order.GET("/payment_methods", h.GetPaymentMethods)
				order.GET("/pick_up_points", h.GetPickUpPoints)
			}
			products := v1.Group("/products")
			{
				products.GET("/", h.GetAllProductsByParams)
				products.POST("/review", h.CreateReview)
				products.GET("/reviews", h.GetReviews)
			}
			v1.GET("/product_detail", h.GetProductDetail)
			cms := v1.Group("/cms")
			{
				cms.POST("/request_call", h.RequestCall)
				cms.GET("/courier_info", h.GetCourierDeliveryInfo)
				cms.GET("/cdek_info", h.GetCDEKDeliveryInfo)
				cms.GET("/vacancies", h.GetAllVacancies)
				cms.POST("/vacancy_request", h.RequestVacancy)
				cms.GET("/requisites", h.GetRequisites)
				cms.GET("/privacy_policy", h.GetPrivacyPolicy)
			}
			categories := v1.Group("/categories")
			{
				categories.GET("/", h.GetAllCategories)
			}
			v1.GET("/category_detail", h.GetCategoriesById)
			user := v1.Group("/user", h.UserValidator)
			{
				user.GET("/info", h.GetUserData)
				user.POST("/update_email", h.UpdateEmailUser)
				user.POST("/update_email_sender", h.CanToSendEmail)
				user.POST("/update_manager_name", h.UpdateManagerName)
			}
		}
	}

	return router
}

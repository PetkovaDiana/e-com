package handler

import "C"
import (
	"clean_arch/internal/dto"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Security ApiKeyAuth
// @Summary Create order
// @Tags Order
// @Description Create order user order
// @Accept json
// @Produce json
// @Param input body dto.Order true "order info"
// @Success 200 {object} dto.DefaultData{data=dto.ReturnOrder}
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/order [post]
func (h *HttpHandler) CreateOrder(c *gin.Context) {
	id, ok := c.Get(userCtx)
	if !ok {
		id = 0
	}
	securityTkn, ok := c.Get(securityToken)
	if !ok {
		securityTkn = ""
	}
	var orderInfo *dto.Order
	if err := c.ShouldBindJSON(&orderInfo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	orderData, err := h.service.CreateOrder(c.Request.Context(), id.(int), securityTkn.(string), orderInfo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": orderData,
	})
}

// @Security ApiKeyAuth
// @Summary Cancel order
// @Tags Order
// @Description Cancel user order
// @Accept json
// @Param order_id query string true "Order id"
// @Success 200
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/order/cancel [post]
func (h *HttpHandler) CancelUserOrder(c *gin.Context) {
	id, ok := c.Get(userCtx)
	if !ok {
		id = 0
	}
	securityTkn, ok := c.Get(securityToken)
	if !ok {
		securityTkn = ""
	}
	orderId := c.Request.URL.Query().Get("order_id")
	if orderId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprint("invalid query param"),
		})
		return
	}
	orderIdDTO, _ := strconv.Atoi(orderId)
	if err := h.service.CancelOrder(c.Request.Context(), id.(int), securityTkn.(string), orderIdDTO); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

// @Security ApiKeyAuth
// @Summary Get user orders
// @Tags Order
// @Description Get all user orders
// @Accept json
// @Produce json
// @Success 200 {object} dto.CountResponse{count=dto.Count,data=[]dto.GetOrder}
// @Param limit query string false "Orders limit in response"
// @Param page query string false "Response page"
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/order [get]
func (h *HttpHandler) GetUserOrders(c *gin.Context) {
	id, _ := c.Get(userCtx)
	allParams := h.ParseUrlParams(c)
	orderInfo, totalOrders, err := h.service.GetUserOrders(c.Request.Context(), id.(int), allParams)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var ordersInfoResponse, totalOrdersResponse any
	if orderInfo == nil {
		ordersInfoResponse = []string{}
		totalOrdersResponse = []string{}
	} else {
		totalOrdersResponse = totalOrders
		ordersInfoResponse = orderInfo
	}
	c.JSON(http.StatusOK, gin.H{
		"count": totalOrdersResponse,
		"data":  ordersInfoResponse,
	})
}

// @Security ApiKeyAuth
// @Summary Get user products
// @Tags Order
// @Description Get all user products
// @Accept json
// @Produce json
// @Success 200 {object} dto.CountResponse{count=dto.Count,data=[]dto.UserOrdersProducts}
// @Param limit query string false "Orders limit in response"
// @Param page query string false "Response page"
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/order_products [get]
func (h *HttpHandler) GetOrdersProducts(c *gin.Context) {
	id, _ := c.Get(userCtx)
	allParams := h.ParseUrlParams(c)
	productsInfo, totalProducts, err := h.service.GetOrdersProducts(c.Request.Context(), id.(int), allParams)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var productsInfoResponse, totalProductsResponse any
	if productsInfo == nil {
		productsInfoResponse = []string{}
		totalProductsResponse = []string{}
	} else {
		productsInfoResponse = productsInfo
		totalProductsResponse = totalProducts
	}
	c.JSON(http.StatusOK, gin.H{
		"count": totalProductsResponse,
		"data":  productsInfoResponse,
	})

}

// @Summary Get payments methods
// @Tags Order
// @Description Get all payments methods
// @Success 200 {object} dto.DefaultData{data=[]dto.PaymentMethod}
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/order/payment_methods [get]
func (h *HttpHandler) GetPaymentMethods(c *gin.Context) {
	paymentMethods := h.service.GetPaymentMethods(c.Request.Context())

	if paymentMethods == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"data": []string{},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": paymentMethods,
	})
}

// @Summary Get pick up points
// @Tags Order
// @Produce json
// @Description Get all pick up points
// @Success 200 {object} dto.DefaultData{data=[]dto.PickUpPoint}
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/order/pick_up_points [get]
func (h *HttpHandler) GetPickUpPoints(c *gin.Context) {
	pickUpPoints := h.service.GetPickUpPoints(c.Request.Context())
	if pickUpPoints == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"data": []string{},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": pickUpPoints,
	})
}

// @Summary Get receipt
// @Tags Order
// @Produce json
// @Description Get order receipt
// @Param promo_code query string false "promo code"
// @Success 200 {object} dto.DefaultData{data=dto.Receipt}
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/order/receipt [post]
func (h *HttpHandler) OrderReceipt(c *gin.Context) {
	id, _ := c.Get(userCtx)
	var receiptInfo *dto.GetReceipt
	if err := c.ShouldBindJSON(&receiptInfo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "wrong json input",
		})
		return
	}
	receipt, err := h.service.OrderReceipt(c.Request.Context(), id.(int), receiptInfo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": receipt,
	})
}

// @Summary All delivery type
// @Tags Order
// @Produce json
// @Description Get all delivery type
// @Produce json
// @Success 200 {object} dto.DefaultData{data=dto.DeliveryTypeInfo}
// @Failure 401 {object} dto.ErrorResponse
// @Router /api/v1/order/delivery_types [get]
func (h *HttpHandler) GetDeliveryType(c *gin.Context) {
	deliveryTypesDTO := h.service.GetDeliveryType(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{
		"data": deliveryTypesDTO,
	})
}

// @Summary Validate online payment
// @Tags Order
// @Router /api/payment_validator [post]
func (h *HttpHandler) OnlinePaymentValidator(c *gin.Context) {
	var orderInfo *dto.OnlineOrderChecker

	id := c.Request.FormValue("id")
	clientID := c.Request.FormValue("clientid")
	sum, _ := strconv.ParseFloat(c.Request.FormValue("sum"), 64)
	orderID := c.Request.FormValue("orderid")
	key := c.Request.FormValue("key")

	orderInfo = &dto.OnlineOrderChecker{
		ID:       id,
		Sum:      sum,
		ClientID: clientID,
		OrderID:  orderID,
		Key:      key,
	}

	responseStr, err := h.service.OnlinePaymentValidator(c.Request.Context(), orderInfo)
	if err != nil {
		fmt.Print(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	_, err = fmt.Fprintf(c.Writer, responseStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	c.Status(http.StatusOK)
}

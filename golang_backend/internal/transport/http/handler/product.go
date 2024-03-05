package handler

import (
	"clean_arch/internal/dto"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary All products
// @Tags products
// @Description Get all products
// @Accept json
// @Produce json
// @Success 200 {object} dto.CountResponse{count=dto.Count,data=[]dto.Product,sort_params=dto.SortParams}
// @Failure 404 {object} dto.ErrorResponse
// @Param limit query string false "Product limit in response"
// @Param page query string false "Response page"
// @Param cat_id query string false "Category uuid"
// @Param price_min query string false "Minimal products price"
// @Param price_max query string false "Maximum products price"
// @Param rating_min query string false "Minimal products rating"
// @Param rating_max query string false "Maximum products rating"
// @Param prod_id query string false "Products uuid"
// @Param sort query string false "Sort settings: popular/lower_price/higher_price/news/discounts/default"
// @Param not_empty query string false "Is quantity are > 0"
// @Router /api/v1/products [get]
func (h *HttpHandler) GetAllProductsByParams(c *gin.Context) {
	allParams := h.ParseUrlParams(c)
	productsDTO, err := h.service.GetAllProductsByParams(c.Request.Context(), allParams)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"count":       "",
			"sort_params": "",
			"data":        []string{},
		})
		return
	}
	if len(productsDTO.Product) == 0 || productsDTO.Product == nil || productsDTO.SortParams == nil {
		c.JSON(http.StatusOK, gin.H{
			"count":       "",
			"sort_params": "",
			"data":        []string{},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"count":       productsDTO.Count,
		"sort_params": productsDTO.SortParams,
		"data":        productsDTO.Product,
	})
}

// @Summary Product detail
// @Tags products
// @Description Get product detail with related products
// @Accept json
// @Produce json
// @Success 200 {object} dto.DefaultData{data=dto.ProductInformation}
// @Failure 404 {object} dto.ErrorResponse
// @Param prod_id query string false "Product ID"
// @Router /api/v1/product_detail [get]
func (h *HttpHandler) GetProductDetail(c *gin.Context) {
	productId := c.Query("prod_id")
	productInfo, err := h.service.GetProductDetail(c.Request.Context(), productId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"data": []string{},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": productInfo,
	})
}

// @Summary Create product review
// @Tags products
// @Description Create product review
// @Accept json
// @Produce json
// @Param input body dto.Review true "review info"
// @Success 200
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/products/review [post]
func (h *HttpHandler) CreateReview(c *gin.Context) {
	reviewDTO := new(dto.Review)
	if err := c.ShouldBindJSON(reviewDTO); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "wrong json input"})
		return
	}

	id, ok := c.Get(userCtx)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user unauthorized"})
		return
	}

	if len(reviewDTO.Image) > 5 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "wrong json input"})
		return
	}

	if err := h.service.CreateReview(c.Request.Context(), reviewDTO, id.(int)); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// @Summary Get product reviews
// @Tags products
// @Description Get all product reviews
// @Accept json
// @Produce json
// @Param prod_id query string true "Product ID"
// @Param limit query string false "Reviews limit in response"
// @Param page query string false "Response page"
// @Success 200 {object} dto.DefaultData{data=dto.ProductStatistic}
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/products/reviews [get]
func (h *HttpHandler) GetReviews(c *gin.Context) {
	productId := c.Query("prod_id")
	if productId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprint("product_id is empty"),
		})
		return
	}
	allParams := h.ParseUrlParams(c)
	reviewsDTO, err := h.service.GetReviews(c.Request.Context(), productId, allParams)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if reviewsDTO == nil {
		c.JSON(http.StatusOK, gin.H{
			"data": []string{},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": reviewsDTO,
	})
}

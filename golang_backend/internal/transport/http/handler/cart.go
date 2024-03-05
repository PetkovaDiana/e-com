package handler

import (
	"clean_arch/internal/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Security ApiKeyAuth
// @Summary User cart
// @Tags cart
// @Description Get actual user cart
// @Produce json
// @Success 200 {object} dto.DefaultData{data=[]dto.Cart}
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/cart [get]
func (h *HttpHandler) GetUserCart(c *gin.Context) {
	id, _ := c.Get(userCtx)
	userCart, err := h.service.GetUserCart(c.Request.Context(), id.(int))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"data": []string{},
		})
		return
	}
	if userCart[0].Product == nil {
		c.JSON(http.StatusOK, gin.H{
			"data": []string{},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": userCart,
	})

}

// @Security ApiKeyAuth
// @Summary Update user cart
// @Tags cart
// @Description Update actual user cart
// @Accept json
// @Param input body dto.UpdateCart true "product info"
// @Success 200
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/cart/add [post]
func (h *HttpHandler) UpdateUserCart(c *gin.Context) {
	id, _ := c.Get(userCtx)

	var product *dto.UpdateCart

	if err := c.ShouldBindJSON(&product); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "wrong json input",
		})
		return
	}
	err := h.service.UpdateUserCart(c.Request.Context(), id.(int), product)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

// @Security ApiKeyAuth
// @Summary Delete product from cart
// @Tags cart
// @Description Delete product from actual user cart
// @Accept json
// @Param input body dto.DeleteCart true "product info"
// @Success 200
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/cart/del [post]
func (h *HttpHandler) DeleteProduct(c *gin.Context) {
	id, _ := c.Get(userCtx)

	var product *dto.DeleteCart

	if err := c.ShouldBindJSON(&product); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "wrong json input",
		})
		return
	}
	err := h.service.DeleteProduct(c.Request.Context(), id.(int), product.ProductUUID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

// @Security ApiKeyAuth
// @Summary Clear user cart
// @Tags cart
// @Description Clear user cart
// @Success 200
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/cart/clear [post]
func (h *HttpHandler) ClearUserCart(c *gin.Context) {
	id, _ := c.Get(userCtx)

	err := h.service.ClearUserCart(c.Request.Context(), id.(int))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

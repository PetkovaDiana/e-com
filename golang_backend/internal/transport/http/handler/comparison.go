package handler

import (
	"clean_arch/internal/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Security ApiKeyAuth
// @Summary User comparison list
// @Tags comparison
// @Description Get actual user comparison list
// @Produce json
// @Success 200 {object} dto.CountResponse{count=dto.Count,data=[]dto.Comparison}
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/comparison [get]
func (h *HttpHandler) GetUserComparison(c *gin.Context) {
	id, _ := c.Get(userCtx)

	userComparison, prodCount, err := h.service.GetUserComparison(c.Request.Context(), id.(int))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(userComparison) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"count": []string{},
			"data":  []string{},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"count": prodCount,
		"data":  userComparison,
	})
}

// @Security ApiKeyAuth
// @Summary Update user comparison list
// @Tags comparison
// @Description Update actual user comparison list
// @Accept json
// @Param input body dto.UpdateComparison true "product info"
// @Success 200
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/comparison/add [post]
func (h *HttpHandler) UpdateUserComparison(c *gin.Context) {
	id, _ := c.Get(userCtx)

	var product *dto.UpdateComparison

	if err := c.ShouldBindJSON(&product); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "wrong json input",
		})
		return
	}
	err := h.service.UpdateUserComparison(c.Request.Context(), id.(int), product)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

// @Security ApiKeyAuth
// @Summary Delete product from comparison list
// @Tags comparison
// @Description Delete product from actual user comparison list
// @Accept json
// @Param input body dto.DeleteComparison true "product info"
// @Success 200
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/comparison/del [post]
func (h *HttpHandler) DeleteComparisonProduct(c *gin.Context) {
	id, _ := c.Get(userCtx)

	var product *dto.DeleteComparison

	if err := c.ShouldBindJSON(&product); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "wrong json input",
		})
		return
	}
	err := h.service.DeleteComparisonProduct(c.Request.Context(), id.(int), product.ProductUUID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

// @Security ApiKeyAuth
// @Summary Delete products from comparison list by category
// @Tags comparison
// @Description Delete products from actual user comparison list by category
// @Accept json
// @Param input body dto.DeleteComparisonByCategory true "category info"
// @Success 200
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/comparison/del_cat [post]
func (h *HttpHandler) DeleteComparisonProductByCategory(c *gin.Context) {
	id, _ := c.Get(userCtx)

	var category *dto.DeleteComparisonByCategory

	if err := c.ShouldBindJSON(&category); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "wrong json input",
		})
		return
	}

	err := h.service.DeleteComparisonProductByCategoryUUID(c.Request.Context(), id.(int), category.CategoryUUID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

// @Security ApiKeyAuth
// @Summary Clear user comparison list
// @Tags comparison
// @Description Clear user comparison list
// @Success 200
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/comparison/clear [get]
func (h *HttpHandler) ClearUserComparison(c *gin.Context) {
	id, _ := c.Get(userCtx)

	err := h.service.ClearUserComparison(c.Request.Context(), id.(int))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

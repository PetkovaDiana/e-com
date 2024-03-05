package handler

import (
	"clean_arch/internal/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Security ApiKeyAuth
// @Summary User favourite list
// @Tags favourite
// @Description Get actual user favourite list
// @Produce json
// @Success 200 {object} dto.DefaultData{data=[]dto.Favourite}
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/favourite [get]
func (h *HttpHandler) GetUserFavourites(c *gin.Context) {
	id, _ := c.Get(userCtx)

	userFavourites, err := h.service.GetUserFavourites(c.Request.Context(), id.(int))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"data": []string{},
		})
		return
	}
	if userFavourites[0].Product == nil {
		c.JSON(http.StatusOK, gin.H{
			"data": []string{},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": userFavourites,
	})
}

// @Security ApiKeyAuth
// @Summary Update user favourite list
// @Tags favourite
// @Description Update actual user favourite list
// @Accept json
// @Param input body dto.UpdateFavourite true "product info"
// @Success 200
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/favourite/add [post]
func (h *HttpHandler) UpdateUserFavourites(c *gin.Context) {
	id, _ := c.Get(userCtx)

	var product *dto.UpdateFavourite

	if err := c.ShouldBindJSON(&product); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "wrong json input",
		})
		return
	}
	err := h.service.UpdateUserFavourites(c.Request.Context(), id.(int), product)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

// @Security ApiKeyAuth
// @Summary Delete product from favourite list
// @Tags favourite
// @Description Delete product from actual user favourite list
// @Accept json
// @Param input body dto.DeleteFavourite true "product info"
// @Success 200
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/favourite/del [post]
func (h *HttpHandler) DeleteFavouriteProduct(c *gin.Context) {
	id, _ := c.Get(userCtx)

	var product *dto.DeleteFavourite

	if err := c.ShouldBindJSON(&product); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "wrong json input",
		})
		return
	}
	err := h.service.DeleteFavouriteProduct(c.Request.Context(), id.(int), product.ProductUUID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

// @Security ApiKeyAuth
// @Summary Clear user favourite list
// @Tags favourite
// @Description Clear user favourite list
// @Success 200
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/favourite/clear [post]
func (h *HttpHandler) ClearUserFavourite(c *gin.Context) {
	id, _ := c.Get(userCtx)

	err := h.service.ClearUserFavourite(c.Request.Context(), id.(int))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

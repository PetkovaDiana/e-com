package handler

import (
	"clean_arch/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// @Summary All categories
// @Tags categories
// @Description Get all categories
// @Accept json
// @Produce json
// @Success 200 {object} dto.DefaultData{data=[]dto.CategorySwagger}
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/categories [get]
func (h *HttpHandler) GetAllCategories(c *gin.Context) {
	categoriesDTO, err := h.service.GetAllCategories(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"data": []string{}})
		return
	}
	if len(categoriesDTO) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"data": []string{},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": categoriesDTO,
	})
}

// @Summary Category layer by id
// @Tags categories
// @Description Get category layer by id
// @Accept json
// @Produce json
// @Param cat_id query string true "ID of the required category"
// @Success 200 {object} dto.CategorySwagger
// @Failure 404,400 {array} error "Error"
// @Router /api/v1/category_detail [get]
func (h *HttpHandler) GetCategoriesById(c *gin.Context) {
	categoryUUID, err := uuid.Parse(c.Query("cat_id"))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	categoryParamsDTO := &dto.CategoryParams{
		UUID: categoryUUID,
	}
	categoryDataDTO, err := h.service.GetCategoriesById(c.Request.Context(), categoryParamsDTO)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(categoryDataDTO.Categories) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"category_path": "",
			"data":          []string{},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"category_path": categoryDataDTO.CategoryPath,
		"data":          categoryDataDTO.Categories,
	})
}

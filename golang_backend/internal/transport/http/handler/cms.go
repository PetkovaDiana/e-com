package handler

import (
	"clean_arch/internal/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary Request call back
// @Tags cms
// @Description Request call back
// @Accept json
// @Produce json
// @Success 200
// @Param input body dto.RequestCall true "request call back info"
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/v1/cms/request_call [post]
func (h *HttpHandler) RequestCall(c *gin.Context) {
	id, _ := c.Get(userCtx)
	var requestCallDTO *dto.RequestCall
	if err := c.ShouldBindJSON(&requestCallDTO); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "wrong json input",
		})
		return
	}
	if err := h.service.RequestCall(c.Request.Context(), requestCallDTO, id.(int)); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

// @Summary Courier delivery info
// @Tags cms
// @Description Full courier delivery info
// @Produce json
// @Success 200 {object} dto.DefaultData{data=dto.CourierDeliveryInfo}
// @Router /api/v1/cms/courier_info [get]
func (h *HttpHandler) GetCourierDeliveryInfo(c *gin.Context) {
	courierDeliveryInfo := h.service.GetCourierDeliveryInfo(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{
		"data": courierDeliveryInfo,
	})
}

// @Summary CDEK delivery info
// @Tags cms
// @Description Full CDEK delivery info
// @Produce json
// @Success 200 {object} dto.DefaultData{data=dto.CDEKDeliveryInfo}
// @Router /api/v1/cms/cdek_info [get]
func (h *HttpHandler) GetCDEKDeliveryInfo(c *gin.Context) {
	cdekInfoDTO := h.service.GetCDEKDeliveryInfo(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{
		"data": cdekInfoDTO,
	})
}

// @Summary Vacancies info
// @Tags cms
// @Description All vacancies info
// @Produce json
// @Success 200 {object} dto.DefaultData{data=[]dto.Vacancy}
// @Router /api/v1/cms/vacancies [get]
func (h *HttpHandler) GetAllVacancies(c *gin.Context) {
	allVacancies := h.service.GetAllVacancies(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{
		"data": allVacancies,
	})
}

// @Summary Request vacancy
// @Tags cms
// @Description Request vacancy
// @Accept json
// @Produce json
// @Success 200
// @Param input body dto.RequestVacancy true "request vacancy info"
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/v1/cms/vacancy_request [post]
func (h *HttpHandler) RequestVacancy(c *gin.Context) {
	var requestVacancyDTO *dto.RequestVacancy
	if err := c.ShouldBindJSON(&requestVacancyDTO); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "wrong json input",
		})
		return
	}
	if err := h.service.RequestVacancy(c.Request.Context(), requestVacancyDTO); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

// @Summary Get requisites
// @Tags cms
// @Description Get requisites
// @Produce json
// @Success 200
// @Success 200 {object} dto.DefaultData{data=dto.Requisites}
// @Router /api/v1/cms/requisites [get]
func (h *HttpHandler) GetRequisites(c *gin.Context) {
	requisitesDTO := h.service.GetRequisites(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{
		"data": requisitesDTO,
	})
}

// @Summary Get privacy policy
// @Tags cms
// @Description Get privacy policy
// @Produce json
// @Success 200
// @Success 200 {object} dto.DefaultData{data=dto.PrivacyPolicy}
// @Router /api/v1/cms/privacy_policy [get]
func (h *HttpHandler) GetPrivacyPolicy(c *gin.Context) {
	privacyPolicyDTO := h.service.GetPrivacyPolicy(c.Request.Context())
	c.JSON(http.StatusOK, gin.H{
		"data": privacyPolicyDTO,
	})
}

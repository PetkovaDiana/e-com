package handler

import (
	"clean_arch/internal/dto"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary Sign-up
// @Tags auth
// @Description Register user
// @ID create-account
// @Accept json
// @Produce json
// @Param input body dto.RegisterUser true "account info"
// @Success 200 {object} dto.Token
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/auth/sign-up [post]
func (h *HttpHandler) RegisterUser(c *gin.Context) {

	id, _ := c.Get(userCtx)
	session, _ := c.Get(unregisteredUserSessionCtx)

	if session == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "user unauthorized",
		})
		return
	}
	if id == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "user unauthorized",
		})
		return
	}
	var userDTO *dto.RegisterUser

	if err := c.ShouldBindJSON(&userDTO); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": ErrInput,
		})
		return
	}
	token, err := h.service.RegisterUser(c.Request.Context(), userDTO, session.(string), id.(int))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": fmt.Sprint(err),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

// @Summary Sign-in
// @Tags auth
// @Description Authenticate user
// @ID authenticate-account
// @Accept json
// @Produce json
// @Param input body dto.UserAuth true "account info"
// @Success 200 {object} dto.Token
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/auth/sign-in [post]
func (h *HttpHandler) AuthenticateUser(c *gin.Context) {
	id, _ := c.Get(userCtx)
	session, _ := c.Get(unregisteredUserSessionCtx)

	if session == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "user unauthorized",
		})
		return
	}
	if id == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "user unauthorized",
		})
		return
	}
	var userDTO *dto.UserAuth
	if err := c.ShouldBindJSON(&userDTO); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": ErrInput,
		})
		return
	}

	token, err := h.service.AuthenticateUser(c.Request.Context(), userDTO, session.(string), id.(int))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": fmt.Sprint(err),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

// @Summary Reg generate sms-code
// @Tags auth
// @Description Send user code for registration
// @ID generate_reg_phone_code
// @Accept json
// @Produce json
// @Param input body dto.CodeGenerate true "user info"
// @Success 200 {object} dto.DefaultData{data=[]dto.CodeResponse}
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/code-reg [post]
func (h *HttpHandler) RegCodeGenerator(c *gin.Context) {
	var code *dto.CodeGenerate
	if err := c.ShouldBindJSON(&code); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": ErrInput,
		})
		return
	}
	err := h.service.RegCodeGenerator(c.Request.Context(), code)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

// @Summary Auth generate sms-code
// @Tags auth
// @Description Send user code for authentication
// @ID generate_auth_phone_code
// @Accept json
// @Produce json
// @Param input body dto.CodeGenerate true "user info"
// @Success 200 {object} dto.DefaultData{data=[]dto.CodeResponse}
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/code-auth [post]
func (h *HttpHandler) AuthCodeGenerator(c *gin.Context) {
	var code *dto.CodeGenerate
	if err := c.ShouldBindJSON(&code); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": ErrInput,
		})
		return
	}
	err := h.service.AuthCodeGenerator(c.Request.Context(), code)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

// @Summary User info
// @Tags user
// @Description All user data
// @Produce json
// @Success 200 {object} dto.DefaultData{data=dto.UserData}
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/user/info [get]
func (h *HttpHandler) GetUserData(c *gin.Context) {
	id, _ := c.Get(userCtx)
	userData, err := h.service.GetUserData(c.Request.Context(), id.(int))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": userData,
	})
}

// @Summary User email info
// @Tags user
// @Description Update user email info
// @Accept json
// @Param input body dto.UpdateEmail true "email info"
// @Success 200
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/user/update_email [post]
func (h *HttpHandler) UpdateEmailUser(c *gin.Context) {
	id, _ := c.Get(userCtx)
	var emailInfo *dto.UpdateEmail
	if err := c.ShouldBindJSON(&emailInfo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "wrong json input",
		})
		return
	}
	if err := h.service.UpdateEmailUser(c.Request.Context(), emailInfo, id.(int)); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

// @Summary Update email sender params
// @Tags user
// @Description Update email sender params
// @Accept json
// @Param input body dto.CanToSendEmail true "email params info"
// @Success 200
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/user/update_email_sender [post]
func (h *HttpHandler) CanToSendEmail(c *gin.Context) {
	id, _ := c.Get(userCtx)
	var emailInfo *dto.CanToSendEmail
	if err := c.ShouldBindJSON(&emailInfo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "wrong json input",
		})
		return
	}
	if err := h.service.CanToSendEmail(c.Request.Context(), emailInfo, id.(int)); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

// @Summary Update manager name
// @Tags user
// @Description Update user manager name
// @Accept json
// @Param input body dto.UpdateManagerName true "manager name info"
// @Success 200
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/v1/user/update_manager_name [post]
func (h *HttpHandler) UpdateManagerName(c *gin.Context) {
	id, _ := c.Get(userCtx)
	var userInfo *dto.UpdateManagerName
	if err := c.ShouldBindJSON(&userInfo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "wrong json input",
		})
		return
	}
	if err := h.service.UpdateManagerName(c.Request.Context(), userInfo, id.(int)); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

// @Summary Site review
// @Tags user
// @Description User review about site
// @Accept json
// @Param input body dto.SiteReview true "review info"
// @Success 200
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/site-review [post]
func (h *HttpHandler) CreateSiteReview(c *gin.Context) {
	var userSiteReview *dto.SiteReview

	if err := c.BindJSON(&userSiteReview); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "wrong json input",
		})
		return
	}

	if userSiteReview.Rating < 1 || userSiteReview.Rating > 5 {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.service.CreateSiteReview(c.Request.Context(), userSiteReview); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}

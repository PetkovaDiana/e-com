package handler

import (
	"clean_arch/internal/dto"
	"clean_arch/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const (
	authorizationHeader        = "Authorization"
	userCtx                    = "userId"
	unregisteredUserSessionCtx = "session_id"
	sessionName                = "secret"
	securityToken              = "security_token"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.GetHeader("Origin") {
		case "http://localhost:3000":
			c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		case "https://dev.ufaelectro.ru":
			c.Writer.Header().Set("Access-Control-Allow-Origin", "https://dev.ufaelectro.ru")
		case "https://ufaelectro.ru":
			c.Writer.Header().Set("Access-Control-Allow-Origin", "https://ufaelectro.ru")
		}
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Set-Cookie")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func (h *HttpHandler) UserIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	sessionId, _ := c.Cookie(sessionName)

	if header != "" {
		headerParts := strings.Split(header, " ")

		if len(headerParts) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "empty auth user"})
		}

		userId, err := h.service.ParseToken(headerParts[1])

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		}

		if sessionId != "" {
			userSession, _, _ := h.service.ValidateSession(sessionId)
			http.SetCookie(c.Writer, &http.Cookie{
				Name:     h.cookie.Name,
				Value:    userSession,
				Path:     h.cookie.Path,
				Domain:   h.cookie.Domain,
				MaxAge:   -1,
				Secure:   h.cookie.Secure,
				HttpOnly: h.cookie.HttpOnly,
				SameSite: h.cookie.SameSite,
			})
		}
		c.Set(userCtx, userId)
	} else {
		userSession, userId, err := h.service.ValidateSession(sessionId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "error occurred saving session in db"})
		}

		http.SetCookie(c.Writer, &http.Cookie{
			Name:     h.cookie.Name,
			Value:    userSession,
			Path:     h.cookie.Path,
			Domain:   h.cookie.Domain,
			MaxAge:   h.cookie.MaxAge,
			Secure:   h.cookie.Secure,
			HttpOnly: h.cookie.HttpOnly,
			SameSite: h.cookie.SameSite,
		})

		c.Set(userCtx, userId)
		c.Set(unregisteredUserSessionCtx, userSession)
	}
	c.Next()
}

func (h *HttpHandler) UserValidator(c *gin.Context) {
	id, _ := c.Get(userCtx)
	session, _ := c.Get(unregisteredUserSessionCtx)

	if session != nil || id == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "user unauthorized",
		})
	}
	c.Next()
}

func (h *HttpHandler) ParseUrlParams(c *gin.Context) *dto.Params {
	limit := 0
	page := 1
	priceMin := 0.01
	priceMax := 100000000.00
	ratingMin := 0.00
	ratingMax := 5.00
	sort := []string{}
	difference := "false"
	notNull := "false"
	var productUUID []string
	var allCategoriesId string
	query := c.Request.URL.Query()
	for key, value := range query {
		queryValue := value[len(value)-1]
		switch key {
		case "limit":
			limit, _ = strconv.Atoi(queryValue)
		case "page":
			page, _ = strconv.Atoi(queryValue)
		case "price_min":
			priceMin, _ = strconv.ParseFloat(queryValue, 64)
		case "price_max":
			priceMax, _ = strconv.ParseFloat(queryValue, 64)
		case "cat_id":
			allCategoriesId = queryValue
		case "rating_min":
			ratingMin, _ = strconv.ParseFloat(queryValue, 64)
		case "rating_max":
			ratingMax, _ = strconv.ParseFloat(queryValue, 64)
		case "sort":
			sort = strings.Split(queryValue, ",")
		case "prod_id":
			productUUID = pkg.ToArray(queryValue, reflect.TypeOf(queryValue)).([]string)
		case "difference":
			difference = queryValue
		case "not_empty":
			notNull = queryValue
		}
	}
	return &dto.Params{
		Limit:       limit,
		Page:        page,
		PriceMin:    priceMin,
		PriceMax:    priceMax,
		CatId:       allCategoriesId,
		RatingMin:   ratingMin,
		RatingMax:   ratingMax,
		Sort:        sort,
		Difference:  difference,
		ProductUUID: productUUID,
		NotNull:     notNull,
	}
}

func (h *HttpHandler) GetSecurityToken(c *gin.Context) {
	token, err := h.service.GetSecurityToken(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.Set(securityToken, token)
	c.Next()
}

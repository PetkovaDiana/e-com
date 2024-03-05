package dto

import "clean_arch/internal/dto"

type ErrorResponse struct {
	Error string `json:"error"`
}

type Token struct {
	Token string `json:"token"`
}

type DefaultData struct {
	Data interface{} `json:"data"`
}

type CountResponse struct {
	Count interface{} `json:"count"`
	Data  interface{} `json:"data"`
}

type CategorySwagger struct {
	UUID        string            `json:"id"`
	Title       string            `json:"title"`
	Image       string            `json:"image"`
	SubCategory []dto.SubCategory `json:"sub_categories"`
}

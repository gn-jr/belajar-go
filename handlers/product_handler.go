package handlers

import (
	"belajar-go/services"
)

type ProductHandler struct{
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}


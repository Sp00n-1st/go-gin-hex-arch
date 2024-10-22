package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go-gin-hex-arch/internal/core/domain"
	"go-gin-hex-arch/internal/core/service"
	"go-gin-hex-arch/internal/dto"
	"strconv"
	"strings"
)

type ProductHandler struct {
	Service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{Service: service}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := h.Service.InsertProduct(product)
	if err != nil {
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			var errors []string
			for _, fieldErr := range validationErr {
				errors = append(errors, fmt.Sprintf("%s is %s", fieldErr.Field(), fieldErr.Tag()))
			}
			c.JSON(400, gin.H{"error": strings.Join(errors, ", ")})
			return
		}

		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, dto.ToProductResponse(product))
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	id, _ := strconv.ParseUint(c.Param("id"), 10, 0)
	product.ProductID = uint(id)

	err := h.Service.UpdateProduct(product)
	if err != nil {
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			var errors []string
			for _, fieldErr := range validationErr {
				errors = append(errors, fmt.Sprintf("%s is %s", fieldErr.Field(), fieldErr.Tag()))
			}
			c.JSON(400, gin.H{"error": strings.Join(errors, ", ")})
			return
		}

		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, dto.ToProductResponse(product))
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err = h.Service.DeleteProduct(uint(id))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Status(204)
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	product, err := h.Service.GetProduct(uint(id))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, dto.ToProductResponse(product))
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	products, err := h.Service.GetProducts()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, dto.ToProductResponses(products))
}

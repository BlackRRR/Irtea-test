package http

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/BlackRRR/Irtea-test/internal/product/app"
	"github.com/BlackRRR/Irtea-test/internal/product/domain"
	"github.com/BlackRRR/Irtea-test/internal/product/interfaces/http/dto"
	"errors"
	"github.com/google/uuid"
	"github.com/BlackRRR/Irtea-test/pkg/validator"
)

type ProductsHandler struct {
	productService *app.ProductService
}

func NewProductsHandler(productService *app.ProductService) *ProductsHandler {
	return &ProductsHandler{
		productService: productService,
	}
}

func (h *ProductsHandler) CreateProduct(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var req dto.CreateProductRequest
	if err := validator.ReadRequest(c, &req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	input := app.CreateProductInput{
		Description: req.Description,
		Tags:        req.Tags,
		Price:       req.Price,
		Quantity:    req.Quantity,
	}

	product, err := h.productService.CreateProduct(ctx, input)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := h.mapProductToResponse(product)
	return c.Status(http.StatusCreated).JSON(response)
}

func (h *ProductsHandler) GetProduct(c *fiber.Ctx) error {
	ctx := c.UserContext()

	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Product ID is required",
		})
	}

	productID, err := h.parseProductID(idParam)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID format",
		})
	}

	product, err := h.productService.GetProduct(ctx, productID)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "Product not found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	response := h.mapProductToResponse(product)
	return c.JSON(response)
}

func (h *ProductsHandler) GetProducts(c *fiber.Ctx) error {
	ctx := c.UserContext()

	limitParam := c.Query("limit", "10")
	offsetParam := c.Query("offset", "0")

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetParam)
	if err != nil || offset < 0 {
		offset = 0
	}

	products, err := h.productService.GetProducts(ctx, limit, offset)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	responses := make([]dto.ProductResponse, 0, len(products))
	for _, product := range products {
		responses = append(responses, h.mapProductToResponse(product))
	}

	return c.JSON(fiber.Map{
		"products": responses,
		"limit":    limit,
		"offset":   offset,
	})
}

func (h *ProductsHandler) UpdatePrice(c *fiber.Ctx) error {
	ctx := c.UserContext()

	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Product ID is required",
		})
	}

	productID, err := h.parseProductID(idParam)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID format",
		})
	}

	var req dto.UpdatePriceRequest
	if err := validator.ReadRequest(c, &req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	input := app.UpdatePriceInput{
		ProductID: productID,
		Price:     req.Price,
	}

	product, err := h.productService.UpdatePrice(ctx, input)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "Product not found",
			})
		}
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := h.mapProductToResponse(product)
	return c.JSON(response)
}

func (h *ProductsHandler) AdjustStock(c *fiber.Ctx) error {
	ctx := c.UserContext()

	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Product ID is required",
		})
	}

	productID, err := h.parseProductID(idParam)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID format",
		})
	}

	var req dto.AdjustStockRequest
	if err := validator.ReadRequest(c, &req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	input := app.AdjustStockInput{
		ProductID: productID,
		Quantity:  req.Quantity,
	}

	product, err := h.productService.AdjustStock(ctx, input)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "Product not found",
			})
		}
		if errors.Is(err, domain.ErrInsufficientStock) {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Insufficient stock",
			})
		}
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := h.mapProductToResponse(product)
	return c.JSON(response)
}

func (h *ProductsHandler) mapProductToResponse(product *domain.Product) dto.ProductResponse {
	return dto.ProductResponse{
		ID:          product.ID.String(),
		Description: product.Description,
		Tags:        product.Tags,
		Price:       product.Price.Amount(),
		Quantity:    product.Inventory.Quantity(),
		CreatedAt:   product.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   product.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func (h *ProductsHandler) parseProductID(s string) (domain.ProductID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return domain.ProductID{}, err
	}

	return domain.ProductID(id), err
}

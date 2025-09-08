package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/BlackRRR/Irtea-test/internal/order/app"
	"github.com/BlackRRR/Irtea-test/internal/order/domain"
	productDomain "github.com/BlackRRR/Irtea-test/internal/product/domain"
	userDomain "github.com/BlackRRR/Irtea-test/internal/user/domain"
	"errors"
	"github.com/BlackRRR/Irtea-test/internal/order/interfaces/http/dto"
	"github.com/BlackRRR/Irtea-test/pkg/consts"
	"github.com/google/uuid"
	"github.com/BlackRRR/Irtea-test/pkg/validator"
)

type OrdersHandler struct {
	orderService *app.OrderService
}

func NewOrdersHandler(orderService *app.OrderService) *OrdersHandler {
	return &OrdersHandler{
		orderService: orderService,
	}
}

func (h *OrdersHandler) PlaceOrder(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var req dto.PlaceOrderRequest
	if err := validator.ReadRequest(c, &req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	userID, err := h.parseUserID(req.UserID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	items := make([]app.OrderItemInput, 0, len(req.Items))
	for _, itemReq := range req.Items {
		productID, err := h.parseProductID(itemReq.ProductID)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid product ID format",
			})
		}

		items = append(items, app.OrderItemInput{
			ProductID: productID,
			Quantity:  itemReq.Quantity,
		})
	}

	input := app.PlaceOrderInput{
		UserID: userID,
		Items:  items,
	}

	order, err := h.orderService.PlaceOrder(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, productDomain.ErrInsufficientStock):
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Insufficient stock for one or more items",
			})
		case errors.Is(err, productDomain.ErrProductNotFound):
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "One or more products not found",
			})
		default:
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}
	}

	response := h.mapOrderToResponse(order)
	return c.Status(http.StatusCreated).JSON(response)
}

func (h *OrdersHandler) GetOrder(c *fiber.Ctx) error {
	ctx := c.UserContext()

	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Order ID is required",
		})
	}

	orderID, err := h.parseOrderID(idParam)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID format",
		})
	}

	order, err := h.orderService.GetOrder(ctx, orderID)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	response := h.mapOrderToResponse(order)
	return c.JSON(response)
}

func (h *OrdersHandler) GetUserOrders(c *fiber.Ctx) error {
	ctx := c.UserContext()

	userIDParam := c.Params("userId")
	if userIDParam == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	userID, err := h.parseUserID(userIDParam)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

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

	orders, err := h.orderService.GetUserOrders(ctx, userID, limit, offset)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	responses := make([]dto.OrderResponse, 0, len(orders))
	for _, order := range orders {
		responses = append(responses, h.mapOrderToResponse(order))
	}

	return c.JSON(fiber.Map{
		"orders": responses,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *OrdersHandler) ConfirmOrder(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Order ID is required",
		})
	}

	orderID, err := h.parseOrderID(idParam)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID format",
		})
	}

	order, err := h.orderService.ConfirmOrder(context.Background(), orderID)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		if errors.Is(err, domain.ErrInvalidOrderStatus) {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot confirm order in current status",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	response := h.mapOrderToResponse(order)
	return c.JSON(response)
}

func (h *OrdersHandler) CancelOrder(c *fiber.Ctx) error {
	ctx := c.UserContext()

	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Order ID is required",
		})
	}

	orderID, err := h.parseOrderID(idParam)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid order ID format",
		})
	}

	order, err := h.orderService.CancelOrder(ctx, orderID)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		if errors.Is(err, domain.ErrInvalidOrderStatus) {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot cancel order in current status",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	response := h.mapOrderToResponse(order)
	return c.JSON(response)
}

func (h *OrdersHandler) mapOrderToResponse(order *domain.Order) dto.OrderResponse {
	items := make([]dto.OrderItemResponse, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, dto.OrderItemResponse{
			ProductID:          item.ProductID.String(),
			ProductDescription: item.ProductDescription,
			ProductPrice:       item.ProductPrice.Amount(),
			Quantity:           item.Quantity,
			TotalPrice:         item.TotalPrice().Amount(),
		})
	}

	return dto.OrderResponse{
		ID:         order.ID.String(),
		UserID:     order.UserID.String(),
		Items:      items,
		Status:     string(order.Status),
		TotalPrice: order.TotalPrice.Amount(),
		CreatedAt:  order.CreatedAt.Format(consts.FormatTimeLayout),
		UpdatedAt:  order.UpdatedAt.Format(consts.FormatTimeLayout),
	}
}

func (h *OrdersHandler) parseOrderID(s string) (domain.OrderID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return domain.OrderID{}, err
	}
	return domain.OrderID(id), err
}

func (h *OrdersHandler) parseUserID(s string) (userDomain.UserID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return userDomain.UserID{}, err
	}
	return userDomain.UserID(id), err
}

func (h *OrdersHandler) parseProductID(s string) (productDomain.ProductID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return productDomain.ProductID{}, err
	}
	return productDomain.ProductID(id), err
}

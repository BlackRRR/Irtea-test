package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/BlackRRR/Irtea-test/internal/user/app"
	"github.com/BlackRRR/Irtea-test/internal/user/domain"
	"github.com/BlackRRR/Irtea-test/internal/user/interfaces/http/dto"
	"errors"
	"github.com/BlackRRR/Irtea-test/pkg/consts"
	"github.com/BlackRRR/Irtea-test/pkg/validator"
)

type UsersHandler struct {
	userService *app.UserService
}

func NewUsersHandler(userService *app.UserService) *UsersHandler {
	return &UsersHandler{
		userService: userService,
	}
}

func (h *UsersHandler) Register(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var req dto.RegisterRequest

	err := validator.ReadRequest(c, &req)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	input := app.RegisterInput{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Age:       req.Age,
		Email:     req.Email,
		IsMarried: req.IsMarried,
		Password:  req.Password,
	}

	user, err := h.userService.Register(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserTooYoung):
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "User must be at least 18 years old",
			})
		case errors.Is(err, domain.ErrInvalidPassword):
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Password must be at least 8 characters long",
			})
		case errors.Is(err, domain.ErrUserAlreadyExists):
			return c.Status(http.StatusConflict).JSON(fiber.Map{
				"error": "User with this email already exists",
			})
		default:
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}
	}

	response := h.mapUserToResponse(user)
	return c.Status(http.StatusCreated).JSON(response)
}

func (h *UsersHandler) GetByID(c *fiber.Ctx) error {
	ctx := c.UserContext()

	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	parsedUUID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}
	userID := domain.UserID(parsedUUID)

	user, err := h.userService.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	response := h.mapUserToResponse(user)
	return c.JSON(response)
}

func (h *UsersHandler) mapUserToResponse(user *domain.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID.String(),
		FirstName: user.FullName.FirstName,
		LastName:  user.FullName.LastName,
		FullName:  user.FullName.String(),
		Age:       user.Age,
		IsMarried: user.IsMarried,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(consts.FormatTimeLayout),
		UpdatedAt: user.UpdatedAt.Format(consts.FormatTimeLayout),
	}
}

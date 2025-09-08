package validator

import (
	"github.com/gofiber/fiber/v2"
	"encoding/json"
	"fmt"
	"errors"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func Init() {
	validate = validator.New()
}

func ReadRequest(c *fiber.Ctx, request interface{}) error {
	if err := c.BodyParser(request); err != nil {
		var jute *json.UnmarshalTypeError
		if errors.As(err, &jute) {
			return fmt.Errorf("field %s must be of type %s", jute.Field, jute.Type.String())
		}

		return err
	}

	return validate.StructCtx(c.UserContext(), request)
}

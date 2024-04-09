package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type (
	contextWrapperService interface {
		Bind(data any) error
	}

	contextWrapper struct {
		Context   echo.Context
		validator *validator.Validate
	}
)

func ContextWrapper(ctx echo.Context) contextWrapperService {
	return &contextWrapper{
		Context:   ctx,
		validator: validator.New(),
	}
}

func (c *contextWrapper) Bind(data any) error {
	if err := c.Context.Bind(data); err != nil {
		log.Printf("Error: Bind data filed: %s", err.Error())
	}

	if err := c.validator.Struct(data); err != nil {
		log.Printf("Error: Validator data failed: %s", err.Error())
	}

	return nil
}

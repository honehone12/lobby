package form

import "github.com/labstack/echo/v4"

type FormData interface {
}

func ProcessFormData[F FormData](c echo.Context, ptr *F) error {
	if err := c.Bind(ptr); err != nil {
		return err
	}

	if err := c.Validate(ptr); err != nil {
		return err
	}

	return nil
}

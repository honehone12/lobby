package errres

import (
	"lobby/logger"
	"net/http"

	"github.com/labstack/echo/v4"
)

func BadRequest(err error, logger logger.Logger) error {
	logger.Warn(err)
	return echo.NewHTTPError(http.StatusBadRequest, "invalid input")
}

func ServiceError(err error, logger logger.Logger) error {
	logger.Error(err)
	return echo.NewHTTPError(http.StatusInternalServerError, "unexpected error")
}

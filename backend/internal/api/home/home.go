package home

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"dev_nikki/internal/authN"
	"dev_nikki/internal/logger"
)

func Home(c echo.Context) error {
	claims, err := authN.GetExtractedCustomClaims(c)
	if err != nil {
		logger.Slog.Error(err.Error())
		return c.JSON(http.StatusUnauthorized, "goodbye")
	}

	logger.Slog.Info("get claims from user request", "claims", claims)
	return c.JSON(http.StatusOK, claims)
}

package index

import (
	"os"

	"github.com/labstack/echo/v4"

	"dev_nikki/internal/logger"
)

type WildCardHandler struct {
}

func (h WildCardHandler) FallbackToIndex(c echo.Context) error {
	if _, err := os.Stat("static/dist" + c.Request().URL.Path); err == nil {
		logger.Slog.Info("URL PATH: " + c.Request().URL.Path)
		logger.Slog.Info(c.Request().URL.RequestURI())
		return c.File("./static/dist" + c.Request().URL.Path)
	}
	logger.Slog.Warn("Not found: " + c.Request().URL.Path)
	return c.File("./static/dist/index.html")
}

package index

import (
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
)

type WildCardHandler struct {
}

func (h WildCardHandler) FallbackToIndex(c echo.Context) error {
	if _, err := os.Stat("static/dist" + c.Request().URL.Path); err == nil {
		slog.Info("URL PATH: " + c.Request().URL.Path)
		slog.Info(c.Request().URL.RequestURI())
		return c.File("./static/dist" + c.Request().URL.Path)
	}
	slog.Warn("Not found: " + c.Request().URL.Path)
	return c.File("./static/dist/index.html")
}

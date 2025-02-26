package home

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"dev_nikki/internal/authN"
	"dev_nikki/internal/logger"
)

var CustomClaimTypeAssertionError = errors.New("failed to type assertion: claim is not CustomClaim")

type homeResponse struct {
	UserID         int               `json:"user_id"`
	Username       string            `json:"username"`
	Email          string            `json:"email"`
	SideBarFolders []string          `json:"sidebar_folders"`
	SideBarFiles   map[string]string `json:"sidebar_files"`
}

func Home(c echo.Context) error {
	t, err := authN.ParseJWTCookie(c)
	if err != nil {
		logger.Slog.Error("cause wrong JWT, can't access home", "JWT", t)
		return c.JSON(http.StatusUnauthorized, "goodbye")
	}

	claim, ok := t.Claims.(*authN.CustomClaim)
	if !ok {
		return CustomClaimTypeAssertionError
	}
	logger.Slog.Info("get claim from user request", "claim", claim)
	return c.JSON(http.StatusOK, claim)
}

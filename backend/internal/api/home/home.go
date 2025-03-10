package home

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"dev_nikki/internal/api/response"
	"dev_nikki/internal/authN"
	"dev_nikki/internal/logger"
	"dev_nikki/internal/models"
)

var (
	unAuthorizedError = errors.New("不正なアクセスです。")
	preHomeError      = errors.New("pre-homeで問題が発生してアクセスできません。")

	preHomeFailedResponse = response.PreHomeResponse{
		Common: response.CommonResponse{
			Status:   "failed",
			UserID:   0,
			Username: "",
			ErrorMsg: preHomeError.Error(),
		},
	}
)

func PreHome(c echo.Context) error {
	claims, err := authN.GetExtractedCustomClaims(c)
	if err != nil {
		logger.Slog.Error(err.Error())
		// jsonで返すエラー内容はunAuthorizedError
		return c.JSON(http.StatusUnauthorized, "goodbye")
	}

	logger.Slog.Info("get claims from user request", "claims", claims)

	// models.GetProjectsでUserIDに紐づいたProjectを全て取得し、frontにかえす。
	_, project, err := models.GetProjects(uint(claims.UserID))
	if err != nil {
		logger.Slog.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, preHomeFailedResponse)
	}

	resp := response.PreHomeResponse{
		Common: response.CommonResponse{
			Status:   "success pre-home",
			UserID:   claims.UserID,
			Username: claims.Username,
			ErrorMsg: "",
		},
		Projects: project,
	}

	return c.JSON(http.StatusOK, resp)
}

func Home(c echo.Context) error {
	claims, err := authN.GetExtractedCustomClaims(c)
	if err != nil {
		logger.Slog.Error(err.Error())
		return c.JSON(http.StatusUnauthorized, "goodbye")
	}

	logger.Slog.Info("get claims from user request", "claims", claims)
	return c.JSON(http.StatusOK, claims)
}

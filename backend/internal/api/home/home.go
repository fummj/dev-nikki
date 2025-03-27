package home

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"

	"dev_nikki/internal/api/response"
	"dev_nikki/internal/authN"
	"dev_nikki/internal/logger"
	"dev_nikki/internal/models"
)

var (
	unAuthorizedError    = errors.New("不正なアクセスです。")
	notExistUserError    = errors.New("このアカウントは存在していません。")
	notMatchProjectError = errors.New("このプロジェクトは存在していません。")
	internalServerError  = errors.New("Internal Server Error")

	homeFailedResponse = response.HomeResponse{
		Common: response.CommonResponse{
			Status:   "failed",
			UserID:   0,
			Username: "",
			Email:    "",
			ErrorMsg: internalServerError.Error(),
		},
	}

	phasePreHome = "prehome"
	phaseHome    = "home"
)

type preHomeResponseData struct {
	ProjectName string `json:"project_name"`
	Description string `json:"description"`
}

func Home(c echo.Context) error {
	claims, err := authN.GetExtractedCustomClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, homeFailedResponse)
	}

	logger.Slog.Info("get claims from user request", "claims", claims)

	u, err := models.GetExistUser(claims.Email)
	if err != nil {
		return c.JSON(http.StatusNotFound, homeFailedResponse)
	}

	_, project, err := models.GetProject(u.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, homeFailedResponse)
	}

	resp := &response.HomeResponse{
		Common: response.CommonResponse{
			Status:   "success home",
			UserID:   claims.UserID,
			Username: claims.Username,
			Email:    claims.Email,
			ErrorMsg: "",
		},
		Phase:   phaseHome,
		Project: project,
	}

	return c.JSON(http.StatusOK, resp)
}

func PreHome(c echo.Context) error {
	claims, err := authN.GetExtractedCustomClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, homeFailedResponse)
	}

	logger.Slog.Info("get claims from user request", "claims", claims)

	u, err := models.GetExistUser(claims.Email)
	if err != nil {
		return c.JSON(http.StatusNotFound, homeFailedResponse)
	}

	// models.GetProjectsでUserIDに紐づいたProjectを全て取得し、frontにかえす。
	_, projects, err := models.GetProjects(u.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, homeFailedResponse)
	}

	resp := &response.PreHomeResponse{
		Common: response.CommonResponse{
			Status:   "success prehome",
			UserID:   claims.UserID,
			Username: claims.Username,
			Email:    claims.Email,
			ErrorMsg: "",
		},
		Phase:    phasePreHome,
		Projects: projects,
	}

	phase := c.Param("phase")
	if phase != phaseHome && phase != "" {
		fmt.Println("phase: ", phase)
		logger.Slog.Info("prehome")
		resp.Phase = phasePreHome
		return c.JSON(http.StatusOK, resp)
	}

	return c.JSON(http.StatusOK, resp)
}

// ユーザーから送られてきたproject_nameに一致するプロジェクトを返す。
func getMatchProject(p []models.Project, name string) (*models.Project, error) {
	for _, pr := range p {
		if pr.Name == name {
			return &pr, nil
		}
	}

	logger.Slog.Info("not exist project_name", "project_name", name)
	return &models.Project{}, notMatchProjectError
}

// ユーザーから送られてきたproject_nameが存在しない場合に新たにプロジェクトを作成する。
func getProject(c echo.Context, u models.User, p []models.Project) (*models.Project, bool, error) {
	var data preHomeResponseData

	var isNew bool

	contentType := c.Request().Header.Get("Content-Type")
	if contentType == "application/json" {
		content, err := io.ReadAll(c.Request().Body)

		err = json.Unmarshal(content, &data)
		if err != nil {
			logger.Slog.Error(err.Error())
			return &models.Project{}, isNew, err
		}
	}

	project, err := getMatchProject(p, data.ProjectName)
	if err != nil {
		_, project, err := models.CreateProject(data.ProjectName, data.Description, u.ID)
		if err != nil {
			logger.Slog.Error(err.Error())
			return &models.Project{}, isNew, err
		}
		isNew = true
		return project, isNew, nil
	}

	isNew = false
	return project, isNew, nil
}

// project_nameに紐づくfolderを全て取得する。なければjson responseを返す。
func getFolders(u models.User, p models.Project) ([]models.Folder, error) {
	_, folders, err := models.GetFolders(u.ID, p.ID)
	if err != nil {
		logger.Slog.Error(err.Error())
		return []models.Folder{}, err
	}

	return folders, nil
}

func getFilesPerFolder(u models.User, p models.Project, folders []models.Folder) (map[string][]models.File, error) {
	fpf := map[string][]models.File{}

	for _, f := range folders {
		_, files, err := models.GetFiles(u.ID, p.ID, f.ID)
		if err != nil {
			logger.Slog.Error(err.Error())
			return fpf, err
		}

		if len(files) != 0 {
			fpf[f.Name] = files
		}
	}

	return fpf, nil
}

// ユーザーから届いたproject_nameをもとにプロジェクトを生成or取得し関連するfolder, fileを一緒にhomeに流す。
func PostPreHome(c echo.Context) error {
	claims, err := authN.GetExtractedCustomClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, homeFailedResponse)
	}

	logger.Slog.Info("get claims from user request", "claims", claims)

	u, err := models.GetExistUser(claims.Email)
	if err != nil {
		return c.JSON(http.StatusNotFound, homeFailedResponse)
	}

	_, projects, err := models.GetProjects(u.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, homeFailedResponse)
	}

	project, isNew, err := getProject(c, *u, projects)
	if err != nil {
		if errors.Is(err, models.AlreadyExistProjectError) {
			homeFailedResponse.Common.ErrorMsg = err.Error()
			return c.JSON(http.StatusNotFound, homeFailedResponse)
		}

		return c.JSON(http.StatusNotFound, homeFailedResponse)
	}

	// isNewがtrueだったらfolder, fileを探す必要ないのでユーザーに返す。
	if isNew {
		resp := response.NewHomeResponse(u.ID, "success home", u.Username, u.Email, "", phaseHome, *project, []models.Folder{}, map[string][]models.File{})

		logger.Slog.Info("access to home with new project", "response", resp, "isNew", isNew)

		return c.JSON(http.StatusOK, resp)
	}

	// project_nameに紐づいているfolder, fileを全て取得して返す
	folders, err := getFolders(*u, *project)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, homeFailedResponse)
	}

	fpf, err := getFilesPerFolder(*u, *project, folders)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, homeFailedResponse)
	}

	resp := response.NewHomeResponse(u.ID, "success home", u.Username, u.Email, "", phaseHome, *project, folders, fpf)

	logger.Slog.Info("access to home with already exist project", "response", resp)
	return c.JSON(http.StatusOK, resp)
}

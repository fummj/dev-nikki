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
	hasExpiredJWTError   = errors.New("認証情報の期限が切れています。")

	unAuthorizedErrorResponse = response.HomeResponse{
		Common: response.CommonResponse{
			Status:   "failed",
			ErrorMsg: unAuthorizedError.Error(),
		},
	}

	notExistUserErrorResponse = response.HomeResponse{
		Common: response.CommonResponse{
			Status:   "failed",
			ErrorMsg: notExistUserError.Error(),
		},
	}

	notMathProjectErrorResponse = response.HomeResponse{
		Common: response.CommonResponse{
			Status:   "failed",
			ErrorMsg: notMatchProjectError.Error(),
		},
	}

	insernalServerErrorResponse = response.HomeResponse{
		Common: response.CommonResponse{
			Status:   "failed",
			ErrorMsg: notMatchProjectError.Error(),
		},
	}

	hasExpiredJWTErrorResponse = response.HomeResponse{
		Common: response.CommonResponse{
			Status:   "failed",
			ErrorMsg: hasExpiredJWTError.Error(),
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
		return c.JSON(http.StatusUnauthorized, hasExpiredJWTErrorResponse)
	}

	logger.Slog.Info("get claims from user request", "claims", claims)

	u, err := models.GetExistUser(models.DBC.DB, claims.Email)
	if err != nil {
		return c.JSON(http.StatusNotFound, notExistUserErrorResponse)
	}

	_, project, err := models.GetProject(models.DBC.DB, u.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, notMathProjectErrorResponse)
	}

	fs, err := getFolders(u.ID, project.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, insernalServerErrorResponse)
	}

	fpf, err := getFilesPerFolder(u.ID, project.ID, fs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, insernalServerErrorResponse)
	}

	resp := &response.HomeResponse{
		Common: response.CommonResponse{
			Status:   "success home",
			UserID:   claims.UserID,
			Username: claims.Username,
			Email:    claims.Email,
			ErrorMsg: "",
		},
		Phase:          phaseHome,
		Project:        project,
		ProjectFolders: fs,
		FilesPerFolder: fpf,
	}

	return c.JSON(http.StatusOK, resp)
}

func PreHome(c echo.Context) error {
	claims, err := authN.GetExtractedCustomClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, hasExpiredJWTErrorResponse)
	}

	logger.Slog.Info("get claims from user request", "claims", claims)

	u, err := models.GetExistUser(models.DBC.DB, claims.Email)
	if err != nil {
		return c.JSON(http.StatusNotFound, notMathProjectErrorResponse)
	}

	// models.GetProjectsでUserIDに紐づいたProjectを全て取得し、frontにかえす。
	_, projects, err := models.GetProjects(models.DBC.DB, u.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, notMathProjectErrorResponse)
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
		_, project, err := models.CreateProject(models.DBC.DB, data.ProjectName, data.Description, u.ID)
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
func getFolders(u uint, p uint) ([]models.Folder, error) {
	_, folders, err := models.GetFolders(models.DBC.DB, u, p)
	if err != nil {
		logger.Slog.Error(err.Error())
		return []models.Folder{}, err
	}

	return folders, nil
}

func getFilesPerFolder(u uint, p uint, folders []models.Folder) (map[string][]models.File, error) {
	fpf := map[string][]models.File{}

	for _, f := range folders {
		files, err := models.GetFiles(models.DBC.DB, u, p, f.ID)
		if err != nil {
			logger.Slog.Error(err.Error())
			return fpf, err
		}

		if len(files) != 0 {
			fpf[f.Name] = files
		}
	}

	// folder_idがnullのfileを取得。
	nff, err := models.GetNoFolderFiles(models.DBC.DB, u, p)
	if err != nil {
		logger.Slog.Error(err.Error())
		return fpf, err
	}

	if len(nff) != 0 {
		fpf["null"] = nff
	}

	return fpf, nil
}

// ユーザーから届いたproject_nameをもとにプロジェクトを生成or取得し関連するfolder, fileを一緒にhomeに流す。
func PostPreHome(c echo.Context) error {
	claims, err := authN.GetExtractedCustomClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, hasExpiredJWTErrorResponse)
	}

	logger.Slog.Info("get claims from user request", "claims", claims)

	u, err := models.GetExistUser(models.DBC.DB, claims.Email)
	if err != nil {
		return c.JSON(http.StatusNotFound, notExistUserErrorResponse)
	}

	_, projects, err := models.GetProjects(models.DBC.DB, u.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, notMathProjectErrorResponse)
	}

	project, isNew, err := getProject(c, *u, projects)
	if err != nil {
		if errors.Is(err, models.AlreadyExistProjectError) {
			return c.JSON(http.StatusNotFound, notMathProjectErrorResponse)
		}

		return c.JSON(http.StatusNotFound, notMathProjectErrorResponse)
	}

	// isNewがtrueだったらfolder, fileを探す必要ないのでユーザーに返す。
	if isNew {
		resp := response.NewHomeResponse(u.ID, "success home", u.Username, u.Email, "", phaseHome, *project, []models.Folder{}, map[string][]models.File{})

		logger.Slog.Info("access to home with new project", "response", resp, "isNew", isNew)

		return c.JSON(http.StatusOK, resp)
	}

	// project_nameに紐づいているfolder, fileを全て取得して返す
	folders, err := getFolders(u.ID, project.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, insernalServerErrorResponse)
	}

	fpf, err := getFilesPerFolder(u.ID, project.ID, folders)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, insernalServerErrorResponse)
	}

	resp := response.NewHomeResponse(u.ID, "success home", u.Username, u.Email, "", phaseHome, *project, folders, fpf)
	return c.JSON(http.StatusOK, resp)
}

func UpdateMarkdown(c echo.Context) error {
	// 送ってきたユーザーのUserID, projectID, FolderID, FileIDが一致するデータのFile.Contentをupdateする。
	claims, err := authN.GetExtractedCustomClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, hasExpiredJWTErrorResponse)
	}

	f := new(models.File)
	if err = c.Bind(f); err != nil {
		logger.Slog.Error(err.Error())
		return c.JSON(http.StatusUnauthorized, unAuthorizedErrorResponse)
	}
	logger.Slog.Info(fmt.Sprintf("%+v request-body\n", f))

	if err := authN.VerifyJWTAgainstRequest(claims.UserID, f.UserID); err != nil {
		return c.JSON(http.StatusInternalServerError, unAuthorizedErrorResponse)
	}

	err = models.UpdateFile(models.DBC.DB, f.ID, f.Content)
	if err != nil {
		logger.Slog.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, insernalServerErrorResponse)
	}

	// updateしたfileを取得。
	file, err := models.GetFile(models.DBC.DB, f.ID)
	if err != nil {
		logger.Slog.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, insernalServerErrorResponse)
	}

	folders, err := getFolders(file.UserID, file.ProjectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, insernalServerErrorResponse)
	}

	fpf, err := getFilesPerFolder(file.UserID, file.ProjectID, folders)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, insernalServerErrorResponse)
	}
	resp := response.FileUpdateResponse{File: file, FilesPerFolder: fpf}

	return c.JSON(http.StatusOK, resp)
}

func CreateNewFolder(c echo.Context) error {
	claims, err := authN.GetExtractedCustomClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, hasExpiredJWTErrorResponse)
	}

	fo := new(models.Folder)
	if err = c.Bind(fo); err != nil {
		logger.Slog.Error(err.Error())
		return c.JSON(http.StatusUnauthorized, unAuthorizedErrorResponse)
	}
	logger.Slog.Info(fmt.Sprintf("%+v request-body\n", fo))

	if err := authN.VerifyJWTAgainstRequest(claims.UserID, fo.UserID); err != nil {
		return c.JSON(http.StatusInternalServerError, unAuthorizedErrorResponse)
	}

	f, err := models.CreateFolder(models.DBC.DB, fo.Name, fo.UserID, fo.ProjectID, fo.ParentFolderID)
	if err != nil {
		logger.Slog.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, insernalServerErrorResponse)
	}
	logger.Slog.Info("folder create success", "folder", f)

	fs, err := getFolders(f.UserID, f.ProjectID)
	if err != nil {
		logger.Slog.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, insernalServerErrorResponse)
	}

	fpf, err := getFilesPerFolder(f.UserID, f.ProjectID, fs)

	resp := response.CreateFolderResponse{ProjectFolders: fs, FilesPerFolder: fpf}
	return c.JSON(http.StatusOK, resp)
}

func CreateNewFile(c echo.Context) error {
	claims, err := authN.GetExtractedCustomClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, hasExpiredJWTErrorResponse)
	}

	f := new(models.File)
	if err = c.Bind(f); err != nil {
		logger.Slog.Error(err.Error())
		return c.JSON(http.StatusUnauthorized, unAuthorizedErrorResponse)
	}
	logger.Slog.Info(fmt.Sprintf("%+v request-body\n", f))

	if err := authN.VerifyJWTAgainstRequest(claims.UserID, f.UserID); err != nil {
		return c.JSON(http.StatusInternalServerError, unAuthorizedErrorResponse)
	}

	file, err := models.CreateFile(models.DBC.DB, f.Name, f.UserID, f.ProjectID, f.FolderID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, insernalServerErrorResponse)
	}
	logger.Slog.Info("file create success", "file", f)

	fs, err := getFolders(file.UserID, file.ProjectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, insernalServerErrorResponse)
	}

	fpf, err := getFilesPerFolder(file.UserID, file.ProjectID, fs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, insernalServerErrorResponse)
	}

	resp := response.CreateFileResponse{File: file, ProjectFolders: fs, FilesPerFolder: fpf}
	return c.JSON(http.StatusOK, resp)
}

func DeleteFolder(c echo.Context) error {
	claims, err := authN.GetExtractedCustomClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, hasExpiredJWTErrorResponse)
	}

	fo := new(models.Folder)
	if err = c.Bind(fo); err != nil {
		logger.Slog.Error(err.Error())
		return c.JSON(http.StatusUnauthorized, unAuthorizedErrorResponse)
	}
	logger.Slog.Info(fmt.Sprintf("%+v request-body\n", fo))

	if err := authN.VerifyJWTAgainstRequest(claims.UserID, fo.UserID); err != nil {
		return c.JSON(http.StatusInternalServerError, unAuthorizedErrorResponse)
	}

	err = models.DeleteFolder(models.DBC.DB, fo.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, insernalServerErrorResponse)
	}
	logger.Slog.Info("folder delete success", "folder", fo)

	fs, err := getFolders(fo.UserID, fo.ProjectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, insernalServerErrorResponse)
	}

	fpf, err := getFilesPerFolder(fo.UserID, fo.ProjectID, fs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, insernalServerErrorResponse)
	}

	resp := response.DeleteFolderResponse{ProjectFolders: fs, FilesPerFolder: fpf}
	return c.JSON(http.StatusOK, resp)
}

func DeleteFile(c echo.Context) error {
	claims, err := authN.GetExtractedCustomClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, hasExpiredJWTErrorResponse)
	}

	f := new(models.File)
	if err = c.Bind(f); err != nil {
		logger.Slog.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, unAuthorizedErrorResponse)
	}
	logger.Slog.Info(fmt.Sprintf("%+v request-body\n", f))

	if err := authN.VerifyJWTAgainstRequest(claims.UserID, f.UserID); err != nil {
		return c.JSON(http.StatusInternalServerError, unAuthorizedErrorResponse)
	}

	err = models.DeleteFile(models.DBC.DB, f.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, insernalServerErrorResponse)
	}
	logger.Slog.Info("file delete success", "file", f)

	fs, err := getFolders(f.UserID, f.ProjectID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, insernalServerErrorResponse)
	}

	fpf, err := getFilesPerFolder(f.UserID, f.ProjectID, fs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, insernalServerErrorResponse)
	}

	resp := response.DeleteFileResponse{File: *f, ProjectFolders: fs, FilesPerFolder: fpf}
	return c.JSON(http.StatusOK, resp)
}

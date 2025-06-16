package response

import (
	"dev_nikki/internal/models"
)

// ここで各APIのカテゴリ別にjsonで返す際のresponse情報をstructでまとめておく。

type CommonResponse struct {
	Status   string `json:"status"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	ErrorMsg string `json:"errorMsg"`
}

type LoginResponse struct {
	Common CommonResponse
}

type SignUpResponse struct {
	Common CommonResponse
}

type PreHomeResponse struct {
	Common   CommonResponse
	Phase    string           `json:"phase"`
	Projects []models.Project `json:"projects"`
}

type HomeResponse struct {
	Common         CommonResponse
	Phase          string                   `json:"phase"`
	Project        models.Project           `json:"project"`
	ProjectFolders []models.Folder          `json:"project_folders"`
	FilesPerFolder map[string][]models.File `json:"files_per_folder"`
}

func NewHomeResponse(id uint, s, n, e, errmsg, phase string, pr models.Project, folders []models.Folder, fpf map[string][]models.File) HomeResponse {
	return HomeResponse{
		Common: CommonResponse{
			Status:   s,
			UserID:   id,
			Username: n,
			Email:    e,
			ErrorMsg: errmsg,
		},
		Phase:          phase,
		Project:        pr,
		ProjectFolders: folders,
		FilesPerFolder: fpf,
	}
}

type FileUpdateResponse struct {
	File           models.File              `json:"file"`
	FilesPerFolder map[string][]models.File `json:"files_per_folder"`
}

type CreateFolderResponse struct {
	ProjectFolders []models.Folder          `json:"project_folders"`
	FilesPerFolder map[string][]models.File `json:"files_per_folder"`
}

type CreateFileResponse struct {
	File           models.File              `json:"file"`
	ProjectFolders []models.Folder          `json:"project_folders"`
	FilesPerFolder map[string][]models.File `json:"files_per_folder"`
}

type DeleteFolderResponse struct {
	ProjectFolders []models.Folder          `json:"project_folders"`
	FilesPerFolder map[string][]models.File `json:"files_per_folder"`
}

type DeleteFileResponse struct {
	File           models.File              `json:"file"`
	ProjectFolders []models.Folder          `json:"project_folders"`
	FilesPerFolder map[string][]models.File `json:"files_per_folder"`
}

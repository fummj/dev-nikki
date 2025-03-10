package response

import (
	"dev_nikki/internal/models"
)

// ここで各APIのカテゴリ別にjsonで返す際のresponse情報をstructでまとめておく。

type CommonResponse struct {
	Status   string `json:"status"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
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
	Projects []models.Project `json:"projects"`
}

type HomeResponse struct {
	Common         CommonResponse
	SideBarFolders []string          `json:"sidebar_folders"`
	SideBarFiles   map[string]string `json:"sidebar_files"`
}

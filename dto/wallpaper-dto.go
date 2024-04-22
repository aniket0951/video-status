package dto

import "github.com/aniket0951/Chatrapati-Maharaj/models"

type GetWallPapersDTO struct {
	Recent []models.WallPaper `json:"recent"`
	Olds   []models.WallPaper `json:"olds"`
}

type GetWallPaperRequest struct {
	IsActive  string `json:"is_active"`
	AppTag    string `json:"app_tag"`
	PageSkip  int    `json:"page_skip"`
	PageLimit int    `json:"page_limit"`
}

package dto

import "github.com/aniket0951/Chatrapati-Maharaj/models"

type GetWallPapersDTO struct {
	Recent []models.WallPaper `json:"recent"`
	Olds   []models.WallPaper `json:"olds"`
}

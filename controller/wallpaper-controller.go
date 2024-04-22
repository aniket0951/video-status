package controller

import (
	"net/http"
	"strings"

	"github.com/aniket0951/Chatrapati-Maharaj/dto"
	"github.com/aniket0951/Chatrapati-Maharaj/helper"
	"github.com/aniket0951/Chatrapati-Maharaj/models"
	"github.com/aniket0951/Chatrapati-Maharaj/services"
	"github.com/gin-gonic/gin"
)

type WallPaperController interface {
	AddWallPaper(ctx *gin.Context)
	GetWallPapers(ctx *gin.Context)
	ActiveInActiveWallPaper(ctx *gin.Context)
	WallPaperLiked(ctx *gin.Context)
}

type controller struct {
	wallPaperService services.WallPaperService
}

func NewWallPaperController(ser services.WallPaperService) WallPaperController {
	return &controller{
		wallPaperService: ser,
	}
}

func (c *controller) AddWallPaper(ctx *gin.Context) {
	file, _, _ := ctx.Request.FormFile("wallpaper")
	title := ctx.Request.PostForm.Get("title")

	if title == "" {
		helper.RequestBodyEmptyResponse(ctx)
		return
	}

	args := models.WallPaper{
		Title:    title,
		IsActive: false,
	}

	err := c.wallPaperService.AddWallPaper(file, args)

	if err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.VIDEO_DATA)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.BuildSuccessResponse("WallPaper has been upload", helper.VIDEO_DATA, helper.EmptyObj{})
	ctx.JSON(http.StatusOK, response)
}

func (c *controller) GetWallPapers(ctx *gin.Context) {
	var getWallPaperRequest dto.GetWallPaperRequest

	if err := ctx.ShouldBindJSON(&getWallPaperRequest); err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, "Invalid Params !", helper.WALLPAPER_DATA)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	result, err := c.wallPaperService.GetWallPapers(getWallPaperRequest)

	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.VIDEO_DATA)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
	}

	response := helper.BuildSuccessResponse(helper.FETCHED_SUCCESS, helper.VIDEO_DATA, result)
	ctx.JSON(http.StatusOK, response)
}

func (c *controller) ActiveInActiveWallPaper(ctx *gin.Context) {
	wallPaperId := ctx.Param("id")
	tag := ctx.Param("tag")

	if strings.TrimSpace(wallPaperId) == "" || strings.TrimSpace(tag) == "" {
		helper.RequestBodyEmptyResponse(ctx)
		return
	}

	var isActive = false

	if tag == "ACTIVE" {
		isActive = true
	}

	err := c.wallPaperService.ActiveInActiveWallPaper(wallPaperId, isActive)

	if err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.WALLPAPER_DATA)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := helper.BuildSuccessResponse("WallPaper status has been changed", helper.WALLPAPER_DATA, nil)
	ctx.JSON(http.StatusOK, response)
}

func (c *controller) WallPaperLiked(ctx *gin.Context) {
	wallPaperId := ctx.Param("id")

	if wallPaperId == "" {
		helper.RequestBodyEmptyResponse(ctx)
		return
	}

	err := c.wallPaperService.WallPaperLiked(wallPaperId)

	if err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.WALLPAPER_DATA)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := helper.BuildSuccessResponse("WallPaper has been Liked", helper.WALLPAPER_DATA, helper.EmptyObj{})
	ctx.JSON(http.StatusOK, response)
}

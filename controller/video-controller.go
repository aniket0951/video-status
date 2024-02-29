package controller

import (
	"net/http"

	"github.com/aniket0951/Chatrapati-Maharaj/dto"
	"github.com/aniket0951/Chatrapati-Maharaj/helper"
	"github.com/aniket0951/Chatrapati-Maharaj/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VideoController interface {
	CreateCategory(ctx *gin.Context)
	UpdateCategory(ctx *gin.Context)
	GetAllCategory(ctx *gin.Context)
	DeleteCategory(ctx *gin.Context)

	AddVideo(ctx *gin.Context)
	GetAllVideos(ctx *gin.Context)
	UpdateVideo(ctx *gin.Context)
	DeleteVideo(ctx *gin.Context)
	GetVideoByID(ctx *gin.Context)

	FetchInActiveVideos(ctx *gin.Context)
	ActiveVideo(ctx *gin.Context)
	IncreaseDownloadCount(ctx *gin.Context)
}

type videocontroller struct {
	service services.VideoService
}

func NewVideoController(ser services.VideoService) VideoController {
	return &videocontroller{
		service: ser,
	}
}

func (c *videocontroller) CreateCategory(ctx *gin.Context) {

	category := dto.CreateVideoCategoriesDTO{}
	ctx.BindJSON(&category)

	if (category == dto.CreateVideoCategoriesDTO{}) {
		helper.RequestBodyEmptyResponse(ctx)
		return
	}

	sv := validator.New()

	if svErr := sv.Struct(category); svErr != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, svErr.Error(), helper.VIDEO_CATEGORY, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	res, err := c.service.CreateCategory(category)

	if err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.VIDEO_CATEGORY, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := helper.BuildSuccessResponse(helper.DATA_INSERTED, helper.VIDEO_CATEGORY, res)
	ctx.JSON(http.StatusOK, response)

}

func (c *videocontroller) UpdateCategory(ctx *gin.Context) {
	category := dto.CreateVideoCategoriesDTO{}
	ctx.BindJSON(&category)

	if (category == dto.CreateVideoCategoriesDTO{}) {
		helper.RequestBodyEmptyResponse(ctx)
		return
	}

	sv := validator.New()

	if svErr := sv.Struct(category); svErr != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, svErr.Error(), helper.VIDEO_DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	res, err := c.service.UpdateCategory(category)

	if err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.VIDEO_DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse(helper.UPDATE_SUCCESS, helper.VIDEO_DATA, res)
	ctx.JSON(http.StatusOK, response)
}

func (c *videocontroller) GetAllCategory(ctx *gin.Context) {
	res, err := c.service.GetAllCategory()

	if err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.VIDEO_DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := helper.BuildSuccessResponse(helper.DATA_FOUND, helper.VIDEO_DATA, res)
	ctx.JSON(http.StatusOK, response)
}

func (c *videocontroller) DeleteCategory(ctx *gin.Context) {
	categoryId := ctx.Request.URL.Query().Get("category_id")

	if len(categoryId) <= 0 {
		helper.RequestBodyEmptyResponse(ctx)
		return
	}

	objId, err := primitive.ObjectIDFromHex(categoryId)

	if err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.VIDEO_DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if !primitive.IsValidObjectID(string(categoryId)) {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, helper.INVALID_ID, helper.VIDEO_DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	delErr := c.service.DeleteCategory(objId)

	if delErr != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, delErr.Error(), helper.VIDEO_DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := helper.BuildSuccessResponse(helper.DELETE_SUCCESS, helper.VIDEO_DATA, helper.EmptyObj{})
	ctx.AbortWithStatusJSON(http.StatusOK, response)
}

func (c *videocontroller) AddVideo(ctx *gin.Context) {
	file, _, _ := ctx.Request.FormFile("video")
	title := ctx.Request.PostForm.Get("title")
	desc := ctx.Request.PostForm.Get("desc")
	videoCatId := ctx.Request.PostForm.Get("video_cat_id")

	objID, _ := primitive.ObjectIDFromHex(videoCatId)

	videoToCreate := dto.CreateVideosDTO{
		VideoTitle:        title,
		VideoDescription:  desc,
		IsVideoActive:     false,
		VideoCategoriesID: objID,
	}

	sv := validator.New()

	if svErr := sv.Struct(&videoToCreate); svErr != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, svErr.Error(), helper.VIDEO_DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err := c.service.AddVideo(videoToCreate, file)

	if err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.VIDEO_DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := helper.BuildSuccessResponse(helper.DATA_INSERTED, helper.VIDEO_DATA, helper.EmptyObj{})
	ctx.AbortWithStatusJSON(http.StatusOK, response)
}

func (c *videocontroller) GetAllVideos(ctx *gin.Context) {
	var getAllVideos dto.GetAllVideosRequestDTO

	if err := ctx.ShouldBindQuery(&getAllVideos); err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, helper.CheckErr(err.Error()), helper.VIDEO_DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	var res []dto.GetVideosDTO
	var err error

	if getAllVideos.Tag == "ACTIVE" {
		res, err = c.service.GetAllVideos()
	}
	if getAllVideos.Tag == "INACTIVE" {
		c.FetchInActiveVideos(ctx)
		return
	}

	if err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.VIDEO_DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := helper.BuildSuccessResponse("Video "+helper.DATA_FOUND, helper.VIDEO_DATA, res)
	ctx.JSON(http.StatusOK, response)
}

func (c *videocontroller) UpdateVideo(ctx *gin.Context) {
	videoToUpdate := dto.UpdateVideoDTO{}

	ctx.BindJSON(&videoToUpdate)

	if (videoToUpdate == dto.UpdateVideoDTO{}) {
		helper.RequestBodyEmptyResponse(ctx)
		return
	}

	err := c.service.UpdateVideo(videoToUpdate)

	if err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.VIDEO_DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := helper.BuildSuccessResponse(helper.UPDATE_SUCCESS, helper.VIDEO_DATA, helper.EmptyObj{})
	ctx.JSON(http.StatusOK, response)

	defer ctx.Request.Body.Close()
}

func (c *videocontroller) DeleteVideo(ctx *gin.Context) {
	defer ctx.Request.Body.Close()
	videoId := ctx.Request.URL.Query().Get("video_id")

	if len(videoId) <= 0 {
		helper.RequestBodyEmptyResponse(ctx)
		return
	}

	if !primitive.IsValidObjectID(videoId) {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, "Invalid video id provided please check video id", helper.VIDEO_DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	objId, objErr := primitive.ObjectIDFromHex(videoId)

	if objErr != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, objErr.Error(), helper.VIDEO_DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err := c.service.DeleteVideo(objId)

	if err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.VIDEO_DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := helper.BuildSuccessResponse(helper.DELETE_SUCCESS, helper.VIDEO_DATA, helper.EmptyObj{})
	ctx.JSON(http.StatusOK, response)
}

func (c *videocontroller) FetchInActiveVideos(ctx *gin.Context) {
	res, err := c.service.FetchInActiveVideos()

	if err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.VIDEO_DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := helper.BuildSuccessResponse("Video "+helper.DATA_FOUND, helper.VIDEO_DATA, res)
	ctx.JSON(http.StatusOK, response)
}

// make a inactive video to active video
func (c *videocontroller) ActiveVideo(ctx *gin.Context) {
	video_id := ctx.Param("videoId")

	if !primitive.IsValidObjectID(video_id) {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, helper.INVALID_ID, helper.USER_DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	objId, objErr := primitive.ObjectIDFromHex(video_id)

	if objErr != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, objErr.Error(), helper.USER_DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err := c.service.ActiveVideo(objId)

	if err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.VIDEO_DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := helper.BuildSuccessResponse("Video has been activited", helper.VIDEO_DATA, helper.EmptyObj{})
	ctx.JSON(http.StatusOK, response)
}

func (c *videocontroller) IncreaseDownloadCount(ctx *gin.Context) {
	video_id := ctx.Param("videoId")

	if !primitive.IsValidObjectID(video_id) {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, helper.INVALID_ID, helper.USER_DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	objId, objErr := primitive.ObjectIDFromHex(video_id)

	if objErr != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, objErr.Error(), helper.USER_DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err := c.service.IncreaseDownloadCount(objId)

	if err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.VIDEO_DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := helper.BuildSuccessResponse("Download count has been increased", helper.VIDEO_DATA, helper.EmptyObj{})
	ctx.JSON(http.StatusOK, response)
}

func (c *videocontroller) GetVideoByID(ctx *gin.Context) {
	video_id := ctx.Param("videoId")

	if !primitive.IsValidObjectID(video_id) {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, helper.INVALID_ID, helper.USER_DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	objId, objErr := primitive.ObjectIDFromHex(video_id)

	if objErr != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, objErr.Error(), helper.USER_DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	video, err := c.service.GetVideoByID(objId)

	if err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.VIDEO_DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := helper.BuildSuccessResponse(helper.FETCHED_SUCCESS, helper.VIDEO_DATA, video)
	ctx.JSON(http.StatusOK, response)
}

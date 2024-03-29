package controller

import (
	"errors"
	"fmt"
	"github.com/aniket0951/Chatrapati-Maharaj/dto"
	"github.com/aniket0951/Chatrapati-Maharaj/helper"
	"github.com/aniket0951/Chatrapati-Maharaj/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
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

	VideoFullDetails(ctx *gin.Context)

	CreateVerificationProcess(videoId primitive.ObjectID, verificationStatus string)
}

type videocontroller struct {
	service                  services.VideoService
	userVideoService         services.UserVideoService
	videoVerificationService services.VideoVerificationService
}

func NewVideoController(ser services.VideoService, userVideoServ services.UserVideoService, verificationService services.VideoVerificationService) VideoController {
	return &videocontroller{
		service:                  ser,
		userVideoService:         userVideoServ,
		videoVerificationService: verificationService,
	}
}

func (c *videocontroller) CreateCategory(ctx *gin.Context) {

	category := dto.CreateVideoCategoriesDTO{}
	_ = ctx.BindJSON(&category)

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
	_ = ctx.BindJSON(&category)

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

	if !primitive.IsValidObjectID(categoryId) {
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
		VideoCategoriesID: objID,
	}

	sv := validator.New()

	if svErr := sv.Struct(&videoToCreate); svErr != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, svErr.Error(), helper.VIDEO_DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	res, err := c.service.AddVideo(videoToCreate, file)

	if err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.VIDEO_DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	userErr := c.userVideoService.AddUserVideo(res)
	if userErr != nil {
		fmt.Println(userErr.Error())
	}

	go func() {
		c.CreateVerificationProcess(res, helper.VERIFICATION_PENDING)
	}()

	response := helper.BuildSuccessResponse(helper.DATA_INSERTED, helper.VIDEO_DATA, helper.EmptyObj{})
	ctx.AbortWithStatusJSON(http.StatusOK, response)
}

func (c *videocontroller) GetAllVideos(ctx *gin.Context) {
	res, err := c.service.GetAllVideos()

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
	_ = ctx.BindJSON(&videoToUpdate)

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

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ctx.Request.Body)
}

func (c *videocontroller) DeleteVideo(ctx *gin.Context) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ctx.Request.Body)
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

func (c *videocontroller) VideoFullDetails(ctx *gin.Context) {
	videoId := ctx.Request.URL.Query().Get("video_id")

	if videoId == "" {
		helper.RequestBodyEmptyResponse(ctx)
		return
	}

	if !primitive.IsValidObjectID(videoId) {
		helper.BuildUnprocessableEntityResponse(ctx, errors.New("invalid input passed"))
		return
	}

	objId, objErr := primitive.ObjectIDFromHex(videoId)

	if objErr != nil {
		helper.BuildUnprocessableEntityResponse(ctx, objErr)
		return
	}

	res, err := c.service.VideoFullDetails(objId)

	if err != nil {
		helper.BuildUnprocessableEntityResponse(ctx, err)
		return
	}

	response := helper.BuildSuccessResponse(helper.FETCHED_SUCCESS, helper.VIDEO_DATA, res)
	ctx.JSON(http.StatusOK, response)

}

func (c *videocontroller) CreateVerificationProcess(videoId primitive.ObjectID, verificationStatus string) {
	var videoVerification dto.CreateVideoVerificationDTO

	userId, _ := primitive.ObjectIDFromHex(helper.USER_ID)

	videoVerification.VideoId = videoId
	videoVerification.UserId = userId
	videoVerification.VerificationStatus = verificationStatus

	err := c.videoVerificationService.CreateVerification(videoVerification)

	fmt.Println(err)

}

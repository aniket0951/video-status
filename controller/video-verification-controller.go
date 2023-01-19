package controller

import (
	"errors"
	"fmt"
	"github.com/aniket0951/Chatrapati-Maharaj/dto"
	"github.com/aniket0951/Chatrapati-Maharaj/helper"
	"github.com/aniket0951/Chatrapati-Maharaj/models"
	"github.com/aniket0951/Chatrapati-Maharaj/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type VideoVerificationController interface {
	CreateVerification(ctx *gin.Context)
	GetAllVideosVerification(ctx *gin.Context)
	ApproveOrDeniedVideo(ctx *gin.Context)
	VideosForVerification(ctx *gin.Context)

	PublishedVideo(ctx *gin.Context)
	GetAllPublishData(ctx *gin.Context)

	CreateVerificationNotification(ctx *gin.Context)
	GetUserVerificationNotification(ctx *gin.Context)

	BuildVerificationNotificationData(title string, videoStatus string, videoId primitive.ObjectID, reason string) dto.CreateVerificationNotificationDTO
}

type videoVerificationController struct {
	verificationService services.VideoVerificationService
	videoService        services.VideoService
}

func NewVideoVerificationController(service services.VideoVerificationService, videoService services.VideoService) VideoVerificationController {
	return &videoVerificationController{
		verificationService: service,
		videoService:        videoService,
	}
}

func (c *videoVerificationController) CreateVerification(ctx *gin.Context) {
	verificationToCreate := dto.CreateVideoVerificationDTO{}
	if bindErr := ctx.BindJSON(&verificationToCreate); bindErr != nil {
		helper.RequestBodyEmptyResponse(ctx)
		return
	}

	sv := validator.New()

	if svErr := sv.Struct(&verificationToCreate); svErr != nil {
		helper.BuildUnprocessableEntityResponse(ctx, svErr)
		return
	}

	err := c.verificationService.CreateVerification(verificationToCreate)

	if err != nil {
		helper.BuildUnprocessableEntityResponse(ctx, err)
		return
	}

	response := helper.BuildSuccessResponse(helper.DATA_INSERTED, helper.VIDEO_VERIFICATION, helper.EmptyObj{})
	ctx.JSON(http.StatusOK, response)
}
func (c *videoVerificationController) GetAllVideosVerification(ctx *gin.Context) {
	res, err := c.verificationService.GetAllVideosVerification()

	if err != nil {
		helper.BuildUnprocessableEntityResponse(ctx, err)
		return
	}

	response := helper.BuildSuccessResponse(helper.DATA_INSERTED, helper.VIDEO_VERIFICATION, res)
	ctx.JSON(http.StatusOK, response)

}
func (c *videoVerificationController) ApproveOrDeniedVideo(ctx *gin.Context) {
	videoId := ctx.Request.URL.Query().Get("video_id")
	videoStatus := ctx.Request.URL.Query().Get("video_status")
	reason := ctx.Request.URL.Query().Get("reason")

	if videoStatus == "" || len(videoStatus) == 0 && len(videoId) == 0 {
		helper.RequestBodyEmptyResponse(ctx)
		return
	}

	if !primitive.IsValidObjectID(videoId) {
		helper.BuildUnprocessableEntityResponse(ctx, errors.New(helper.INVALID_ID))
		return
	}

	objId, objErr := primitive.ObjectIDFromHex(videoId)

	if objErr != nil {
		helper.BuildUnprocessableEntityResponse(ctx, objErr)
		return
	}

	err := c.verificationService.ApproveOrDeniedVideo(objId, videoStatus)

	if err != nil {
		helper.BuildUnprocessableEntityResponse(ctx, err)
		return
	}

	msg := fmt.Sprintf("video %s successful", videoStatus)

	go func() {
		notification := c.BuildVerificationNotificationData(msg, videoStatus, objId, reason)
		_ = c.verificationService.CreateVerificationNotification(notification)
	}()

	if videoStatus == helper.VERIFICATION_APPROVE {
		go func() {
			video := models.Videos{
				ID:            objId,
				IsVideoActive: false,
				IsVerified:    true,
				IsPublished:   false,
			}

			_ = c.videoService.UpdateVideoVerification(video)
		}()
	}

	response := helper.BuildSuccessResponse(msg, helper.VIDEO_VERIFICATION, helper.EmptyObj{})
	ctx.JSON(http.StatusOK, response)

}
func (c *videoVerificationController) VideosForVerification(ctx *gin.Context) {
	tag := ctx.Request.URL.Query().Get("tag")

	if tag == "" || len(tag) <= 0 {
		helper.RequestBodyEmptyResponse(ctx)
		return
	}

	res, err := c.verificationService.VideosForVerification(tag)

	if err != nil {
		helper.BuildUnprocessableEntityResponse(ctx, err)
		return
	}

	response := helper.BuildSuccessResponse(helper.FETCHED_SUCCESS, helper.VIDEO_VERIFICATION, res)
	ctx.JSON(http.StatusOK, response)
}

func (c *videoVerificationController) PublishedVideo(ctx *gin.Context) {
	publishToCreate := dto.CreatePublishDTO{}

	if bindErr := ctx.BindJSON(&publishToCreate); bindErr != nil {
		helper.RequestBodyEmptyResponse(ctx)
		return
	}

	sv := validator.New()

	if svErr := sv.Struct(&publishToCreate); svErr != nil {
		helper.BuildUnprocessableEntityResponse(ctx, svErr)
		return
	}

	if err := c.verificationService.PublishedVideo(publishToCreate); err != nil {
		helper.BuildUnprocessableEntityResponse(ctx, err)
		return
	}

	err := c.verificationService.CreateVideoProcessHistory(publishToCreate.VideoId)

	if err != nil {
		helper.BuildUnprocessableEntityResponse(ctx, err)
		return
	}

	go func() {
		video := models.Videos{
			ID:            publishToCreate.VideoId,
			IsVideoActive: true,
			IsVerified:    true,
			IsPublished:   true,
		}

		_ = c.videoService.UpdateVideoVerification(video)
	}()

	go func() {
		notification := c.BuildVerificationNotificationData("Video Published", helper.VIDEO_PUBLISHED, publishToCreate.VideoId, "")
		_ = c.verificationService.CreateVerificationNotification(notification)
	}()

	response := helper.BuildSuccessResponse(helper.DATA_INSERTED, helper.VIDEO_VERIFICATION, helper.EmptyObj{})
	ctx.JSON(http.StatusOK, response)
}
func (c *videoVerificationController) GetAllPublishData(ctx *gin.Context) {
	res, err := c.verificationService.GetAllPublishData()

	if err != nil {
		helper.BuildUnprocessableEntityResponse(ctx, err)
		return
	}

	response := helper.BuildSuccessResponse(helper.DATA_INSERTED, helper.VIDEO_VERIFICATION, res)
	ctx.JSON(http.StatusOK, response)
}

func (c *videoVerificationController) CreateVerificationNotification(ctx *gin.Context) {
	notificationToCreate := dto.CreateVerificationNotificationDTO{}
	_ = ctx.BindJSON(&notificationToCreate)

	if (notificationToCreate == dto.CreateVerificationNotificationDTO{}) {
		helper.RequestBodyEmptyResponse(ctx)
		return
	}

	sv := validator.New()

	if svErr := sv.Struct(&notificationToCreate); svErr != nil {
		helper.BuildUnprocessableEntityResponse(ctx, svErr)
		return
	}

	if err := c.verificationService.CreateVerificationNotification(notificationToCreate); err != nil {
		helper.BuildUnprocessableEntityResponse(ctx, err)
		return
	}

	response := helper.BuildSuccessResponse(helper.DATA_INSERTED, helper.VIDEO_VERIFICATION, helper.EmptyObj{})
	ctx.JSON(http.StatusOK, response)
}
func (c *videoVerificationController) GetUserVerificationNotification(ctx *gin.Context) {
	userId := ctx.Request.URL.Query().Get("user_id")

	if userId == "" {
		helper.RequestBodyEmptyResponse(ctx)
		return
	}

	if !primitive.IsValidObjectID(userId) {
		helper.BuildUnprocessableEntityResponse(ctx, errors.New("invalid user id requested"))
		return
	}

	objId, objErr := primitive.ObjectIDFromHex(userId)

	if objErr != nil {
		helper.BuildUnprocessableEntityResponse(ctx, objErr)
		return
	}

	res, err := c.verificationService.GetUserVerificationNotification(objId)

	if err != nil {
		helper.BuildUnprocessableEntityResponse(ctx, err)
		return
	}

	response := helper.BuildSuccessResponse(helper.FETCHED_SUCCESS, helper.VIDEO_VERIFICATION, res)
	ctx.JSON(http.StatusOK, response)
}

func (c *videoVerificationController) BuildVerificationNotificationData(title string, videoStatus string, videoId primitive.ObjectID, reason string) dto.CreateVerificationNotificationDTO {
	notification := &dto.CreateVerificationNotificationDTO{
		Title:       title,
		Description: reason,
		IsApproved:  videoStatus,
		VideoId:     videoId,
	}

	return *notification
}

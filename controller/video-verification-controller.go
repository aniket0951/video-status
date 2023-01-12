package controller

import (
	"errors"
	"github.com/aniket0951/Chatrapati-Maharaj/dto"
	"github.com/aniket0951/Chatrapati-Maharaj/helper"
	"github.com/aniket0951/Chatrapati-Maharaj/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type VideoVerificationController interface {
	CreateVerification(ctx *gin.Context)
	GetAllVideosVerification(ctx *gin.Context)

	CreatePublish(ctx *gin.Context)
	GetAllPublishData(ctx *gin.Context)

	CreateVerificationNotification(ctx *gin.Context)
	GetUserVerificationNotification(ctx *gin.Context)
}

type videoVerificationController struct {
	verificationService services.VideoVerificationService
}

func NewVideoVerificationController(service services.VideoVerificationService) VideoVerificationController {
	return &videoVerificationController{
		verificationService: service,
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

func (c *videoVerificationController) CreatePublish(ctx *gin.Context) {
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

	if err := c.verificationService.CreatePublish(publishToCreate); err != nil {
		helper.BuildUnprocessableEntityResponse(ctx, err)
		return
	}

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

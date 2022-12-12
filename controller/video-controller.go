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
	ctx.JSON(http.StatusOK, response)
}

package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aniket0951/Chatrapati-Maharaj/dto"
	"github.com/aniket0951/Chatrapati-Maharaj/helper"
	notificationmanager "github.com/aniket0951/Chatrapati-Maharaj/notification_manager"
	"github.com/aniket0951/Chatrapati-Maharaj/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserAuthController interface {
	CreateEndUser(ctx *gin.Context)
	CreateAdminUser(ctx *gin.Context)
	AdminUserLogin(ctx *gin.Context)
	GetUserById(ctx *gin.Context)
	DeleteAdminUser(ctx *gin.Context)
	GetAllAdminUser(ctx *gin.Context)
	UpdateAdminUser(ctx *gin.Context)

	AddAdminUserAddress(ctx *gin.Context)
	GetAdminUserAdrress(ctx *gin.Context)
	UpdateAdminAddress(ctx *gin.Context)

	SaveTokens(ctx *gin.Context)
	GetTokens(ctx *gin.Context)
}

var notificationManager = notificationmanager.NewNotificationManagerRepo()

type userauthcontroller struct {
	service    services.UserAuthService
	jwtService services.JWTService
	// notification Manager
	notificationService notificationmanager.NotificationManagerRepo
}

func NewUserAuthController(service services.UserAuthService, jwtService services.JWTService) UserAuthController {
	return &userauthcontroller{
		service:             service,
		jwtService:          jwtService,
		notificationService: notificationManager,
	}
}

func (c *userauthcontroller) CreateEndUser(ctx *gin.Context) {
	var user dto.RegisterEndUserDTO
	ctx.BindJSON(&user)

	if (user == dto.RegisterEndUserDTO{}) {
		helper.RequestBodyEmptyResponse(ctx)
		return
	}

	if !helper.ValidateNumber(user.MobileNumber) || len(user.MobileNumber) < 10 {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, helper.MOBILE_INVALID, helper.DATA)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	newUser, err := c.service.CreateEndUser(user)

	if err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.DATA)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if (newUser == dto.RegisterEndUserDTO{}) {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, helper.DATA_INSERTED_FAILED, helper.DATA)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := helper.BuildSuccessResponse(helper.DATA_INSERTED, helper.DATA, newUser)
	ctx.JSON(http.StatusOK, response)
}

func (c *userauthcontroller) CreateAdminUser(ctx *gin.Context) {
	var user dto.CreateAdminUserDTO
	ctx.BindJSON(&user)

	sv := validator.New()
	structVal := sv.Struct(&user)

	if structVal != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, structVal.Error(), helper.USER_DATA)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	res, err := c.service.CreateAdminUser(user)

	if err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.USER_DATA)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := helper.BuildSuccessResponse(helper.DATA_INSERTED, helper.USER_DATA, res)
	ctx.JSON(http.StatusOK, response)

}

func (c *userauthcontroller) AdminUserLogin(ctx *gin.Context) {
	userCredentials := dto.AdminLoginDTO{}
	ctx.BindJSON(&userCredentials)
	fmt.Println(userCredentials)
	if (userCredentials == dto.AdminLoginDTO{}) {
		helper.RequestBodyEmptyResponse(ctx)
		return
	}

	sv := validator.New()

	if err := sv.Struct(userCredentials); err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.USER_DATA)
		ctx.AbortWithStatusJSON(http.StatusForbidden, response)
		return
	}

	res, err := c.service.ValidateAdminUser(userCredentials.Email, userCredentials.Password)

	if err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.USER_DATA)
		ctx.AbortWithStatusJSON(http.StatusForbidden, response)
		return
	}

	stringObjectID := res.ID.Hex()
	generateToken := c.jwtService.GenerateToken(string(stringObjectID), res.UserType)
	res.Token = generateToken

	response := helper.BuildSuccessResponse("You are login successfully.", helper.USER_DATA, res)
	ctx.JSON(http.StatusOK, response)
}

func (c *userauthcontroller) GetUserById(ctx *gin.Context) {
	userId := helper.USER_ID

	if userId == "" {
		helper.RequestBodyEmptyResponse(ctx)
		return
	}

	if !primitive.IsValidObjectID(userId) {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, helper.INVALID_ID, helper.USER_DATA)
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	objId, objErr := primitive.ObjectIDFromHex(userId)

	if objErr != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, objErr.Error(), helper.USER_DATA)
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	res, err := c.service.GetUserById(objId)

	if err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.USER_DATA)
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse(helper.FETCHED_SUCCESS, helper.USER_DATA, res)
	ctx.JSON(http.StatusOK, response)
}

func (c *userauthcontroller) UpdateAdminUser(ctx *gin.Context) {
	userToUpdate := dto.UpdateAdminUserDTO{}
	ctx.BindJSON(&userToUpdate)

	if (userToUpdate == dto.UpdateAdminUserDTO{}) {
		helper.RequestBodyEmptyResponse(ctx)
		return
	}

	sv := validator.New()

	if svErr := sv.Struct(&userToUpdate); svErr != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, svErr.Error(), helper.USER_DATA)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if upErr := c.service.UpdateAdminUserInfo(userToUpdate); upErr != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, upErr.Error(), helper.USER_DATA)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := helper.BuildSuccessResponse(helper.UPDATE_SUCCESS, helper.USER_DATA, helper.EmptyObj{})
	ctx.JSON(http.StatusOK, response)
}

func (c *userauthcontroller) DeleteAdminUser(ctx *gin.Context) {
	userId := ctx.Request.URL.Query().Get("userId")

	if userId == "" {
		helper.RequestBodyEmptyResponse(ctx)
		return
	}

	if !primitive.IsValidObjectID(userId) {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, helper.INVALID_ID, helper.USER_DATA)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	objId, objErr := primitive.ObjectIDFromHex(userId)

	if objErr != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, objErr.Error(), helper.USER_DATA)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err := c.service.DeleteAdminUser(objId)

	if err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.USER_DATA)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := helper.BuildSuccessResponse(helper.DELETE_SUCCESS, helper.USER_DATA, helper.EmptyObj{})
	ctx.JSON(http.StatusOK, response)
}

func (c *userauthcontroller) GetAllAdminUser(ctx *gin.Context) {
	res, err := c.service.GetAllAdminUsers()

	if err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.USER_DATA)
		ctx.AbortWithStatusJSON(http.StatusForbidden, response)
		return
	}

	response := helper.BuildSuccessResponse(helper.FETCHED_SUCCESS, helper.USER_DATA, res)
	ctx.JSON(http.StatusOK, response)
}

func (c *userauthcontroller) AddAdminUserAddress(ctx *gin.Context) {
	addressToCreate := dto.CreateAdminUserAddress{}
	ctx.BindJSON(&addressToCreate)

	if (addressToCreate == dto.CreateAdminUserAddress{}) {
		helper.RequestBodyEmptyResponse(ctx)
		return
	}

	userId, userIdConErr := primitive.ObjectIDFromHex(helper.USER_ID)

	if userIdConErr != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, userIdConErr.Error(), helper.USER_DATA)
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	addressToCreate.UserID = userId

	sv := validator.New()

	if svErr := sv.Struct(&addressToCreate); svErr != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, svErr.Error(), helper.USER_DATA)
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	err := c.service.AddAdminUserAddress(addressToCreate)

	if err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.USER_DATA)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := helper.BuildSuccessResponse("Address addedd successfully.", helper.USER_DATA, helper.EmptyObj{})
	ctx.JSON(http.StatusOK, response)
}

func (c *userauthcontroller) GetAdminUserAdrress(ctx *gin.Context) {
	userId, err := primitive.ObjectIDFromHex(helper.USER_ID)

	if err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.USER_DATA)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	address, addErr := c.service.GetAdminUserAddress(userId)

	if addErr != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, addErr.Error(), helper.USER_DATA)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := helper.BuildSuccessResponse(helper.DATA_FOUND, helper.USER_DATA, address)
	ctx.JSON(http.StatusOK, response)
}

func (c *userauthcontroller) UpdateAdminAddress(ctx *gin.Context) {
	addressToUpdate := dto.UpdateAdminAddressDTO{}
	ctx.BindJSON(&addressToUpdate)

	if (addressToUpdate == dto.UpdateAdminAddressDTO{}) {
		helper.RequestBodyEmptyResponse(ctx)
		return
	}

	sv := validator.New()

	if svErr := sv.Struct(&addressToUpdate); svErr != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, svErr.Error(), helper.USER_DATA)
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	err := c.service.UpdateAdminAddress(addressToUpdate)

	if err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.USER_DATA)
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse(helper.UPDATE_SUCCESS, helper.USER_DATA, helper.EmptyObj{})
	ctx.JSON(http.StatusOK, response)
}

func (c *userauthcontroller) SaveTokens(ctx *gin.Context) {
	var token_data notificationmanager.TokenRequestDTO

	if err := ctx.ShouldBindJSON(&token_data); err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.USER_DATA)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	tokenData := notificationmanager.TokenData{
		Token:     token_data.Token,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	if err := c.notificationService.AddNewToken(tokenData); err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.USER_DATA)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.BuildSuccessResponse("Token has been saved !", helper.USER_DATA, helper.EmptyObj{})
	ctx.JSON(http.StatusOK, response)

}

func (c *userauthcontroller) GetTokens(ctx *gin.Context) {
	tokens, err := c.notificationService.GetTokens()

	if err != nil {
		response := helper.BuildFailedResponse(helper.FETCHED_FAILED, err.Error(), helper.USER_DATA)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.BuildSuccessResponse(helper.FETCHED_SUCCESS, "tokenData", tokens)
	ctx.JSON(http.StatusOK, response)
}

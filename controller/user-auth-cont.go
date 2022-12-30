package controller

import (
	"fmt"
	"net/http"

	"github.com/aniket0951/Chatrapati-Maharaj/dto"
	"github.com/aniket0951/Chatrapati-Maharaj/helper"
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
}

type userauthcontroller struct {
	service    services.UserAuthService
	jwtService services.JWTService
}

func NewUserAuthController(service services.UserAuthService, jwtService services.JWTService) UserAuthController {
	return &userauthcontroller{
		service:    service,
		jwtService: jwtService,
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
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, helper.MOBILE_INVALID, helper.DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	newUser, err := c.service.CreateEndUser(user)

	if err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if (newUser == dto.RegisterEndUserDTO{}) {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, helper.DATA_INSERTED_FAILED, helper.DATA, helper.EmptyObj{})
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
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, structVal.Error(), helper.USER_DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	res, err := c.service.CreateAdminUser(user)

	if err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.USER_DATA, helper.EmptyObj{})
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
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.USER_DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusForbidden, response)
		return
	}

	res, err := c.service.ValidateAdminUser(userCredentials.Email, userCredentials.Password)

	if err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.USER_DATA, helper.EmptyObj{})
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
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, helper.INVALID_ID, helper.USER_DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	objId, objErr := primitive.ObjectIDFromHex(userId)

	if objErr != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, objErr.Error(), helper.USER_DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	res, err := c.service.GetUserById(objId)

	if err != nil {
		response := helper.BuildFailedResponse(helper.FAILED_PROCESS, err.Error(), helper.USER_DATA, helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.BuildSuccessResponse(helper.FETCHED_SUCCESS, helper.USER_DATA, res)
	ctx.JSON(http.StatusOK, response)
}

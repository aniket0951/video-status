package routers

import (
	"github.com/aniket0951/Chatrapati-Maharaj/controller"
	dbconfig "github.com/aniket0951/Chatrapati-Maharaj/db-config"
	middleware "github.com/aniket0951/Chatrapati-Maharaj/middelware"
	"github.com/aniket0951/Chatrapati-Maharaj/repositories"
	"github.com/aniket0951/Chatrapati-Maharaj/services"
	"github.com/gin-gonic/gin"
)

var userCollection = dbconfig.GetCollection(dbconfig.DB, "users")

var (
	jwtService   = services.NewJWTService()
	userAuthRepo = repositories.NewUserAuthRepository(userCollection)
	userAuthSer  = services.NewUserAuthService(userAuthRepo)
	userAuthCont = controller.NewUserAuthController(userAuthSer, jwtService)
)

func UserAuthRouter(router *gin.Engine) {
	userAuth := router.Group("/api")
	{
		userAuth.POST("/create-user", userAuthCont.CreateEndUser)
		userAuth.POST("/create-admin-user", userAuthCont.CreateAdminUser)
		userAuth.POST("/admin-user-login", userAuthCont.AdminUserLogin)
		userAuth.GET("/get-user-byID", middleware.AuthorizeJWT(jwtService), userAuthCont.GetUserById)
		userAuth.GET("/get-admin-users", middleware.AuthorizeJWT(jwtService), userAuthCont.GetAllAdminUser)
		userAuth.PUT("/update-admin-user", middleware.AuthorizeJWT(jwtService), userAuthCont.UpdateAdminUser)
		userAuth.DELETE("/delete-admin-user", middleware.AuthorizeJWT(jwtService), userAuthCont.DeleteAdminUser)
	}

	adminUserAddress := router.Group("/api/address", middleware.AuthorizeJWT(jwtService))
	{
		adminUserAddress.POST("/add-admin-address", userAuthCont.AddAdminUserAddress)
		adminUserAddress.GET("/get-admin-address", userAuthCont.GetAdminUserAdrress)
		adminUserAddress.PUT("/update-admin-address", userAuthCont.UpdateAdminAddress)
	}

	tokenApp := router.Group("/api/user", middleware.AuthorizeJWT(jwtService))
	{
		tokenApp.POST("/add-fcm-token", userAuthCont.SaveTokens)
		tokenApp.GET("/get-tokens", userAuthCont.GetTokens)
	}
}

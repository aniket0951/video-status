package routers

import (
	"github.com/aniket0951/Chatrapati-Maharaj/controller"
	dbconfig "github.com/aniket0951/Chatrapati-Maharaj/db-config"
	middleware "github.com/aniket0951/Chatrapati-Maharaj/middelware"
	"github.com/aniket0951/Chatrapati-Maharaj/repositories"
	"github.com/aniket0951/Chatrapati-Maharaj/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = dbconfig.GetCollection(dbconfig.DB, "users")

var (
	jwtService   services.JWTService             = services.NewJWTService()
	userAuthRepo repositories.UserAuthRepository = repositories.NewUserAuthRepository(userCollection)
	userAuthSer  services.UserAuthService        = services.NewUserAuthService(userAuthRepo)
	userAuthCont controller.UserAuthController   = controller.NewUserAuthController(userAuthSer, jwtService)
)

func UserAuthRouter(router *gin.Engine) {
	userAuth := router.Group("/api")
	{
		userAuth.POST("/create-user", userAuthCont.CreateEndUser)
		userAuth.POST("/create-admin-user", userAuthCont.CreateAdminUser)
		userAuth.POST("/admin-user-login", userAuthCont.AdminUserLogin)
		userAuth.GET("/get-user-byID", middleware.AuthorizeJWT(jwtService), userAuthCont.GetUserById)
		userAuth.GET("/get-admin-users", middleware.AuthorizeJWT(jwtService), userAuthCont.GetAllAdminUser)
	}

	adminUserAddress := router.Group("/api/address", middleware.AuthorizeJWT(jwtService))
	{
		adminUserAddress.POST("/add-admin-address", userAuthCont.AddAdminUserAddress)
		adminUserAddress.GET("/get-admin-address", userAuthCont.GetAdminUserAdrress)
		adminUserAddress.PUT("/update-admin-address", userAuthCont.UpdateAdminAddress)
	}
}

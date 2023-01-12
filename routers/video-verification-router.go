package routers

import (
	"github.com/aniket0951/Chatrapati-Maharaj/controller"
	middleware "github.com/aniket0951/Chatrapati-Maharaj/middelware"
	"github.com/aniket0951/Chatrapati-Maharaj/repositories"
	"github.com/gin-gonic/gin"
)

var (
	verificationRepository repositories.VideoVerificationRepository = repositories.NewVideoVerificationRepository()
	verificationController controller.VideoVerificationController   = controller.NewVideoVerificationController(verificationService)
)

func VideoVerificationRoute(route *gin.Engine) {
	verification := route.Group("/api/video-verification", middleware.AuthorizeJWT(jwtService))
	{
		verification.POST("/create-verification", verificationController.CreateVerification)
		verification.GET("/get-all-verification", verificationController.GetAllVideosVerification)
		verification.POST("/create-publish", verificationController.CreatePublish)
		verification.GET("/get-all-publish", verificationController.GetAllPublishData)
	}

	notification := route.Group("/api/verification-notification", middleware.AuthorizeJWT(jwtService))
	{
		notification.POST("/create-notification", verificationController.CreateVerificationNotification)
		notification.GET("/user-notification", verificationController.GetUserVerificationNotification)
	}
}

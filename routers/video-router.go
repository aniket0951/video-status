package routers

import (
	"github.com/aniket0951/Chatrapati-Maharaj/controller"
	middleware "github.com/aniket0951/Chatrapati-Maharaj/middelware"
	"github.com/aniket0951/Chatrapati-Maharaj/repositories"
	"github.com/aniket0951/Chatrapati-Maharaj/services"
	"github.com/gin-gonic/gin"
)

var (
	videoRepo           repositories.VideoRepository             = repositories.NewVideoCategoriesRepository()
	userVideoRepo       repositories.UserVideoRepository         = repositories.NewUserVideoRepository()
	verificationRepo    repositories.VideoVerificationRepository = repositories.NewVideoVerificationRepository()
	userVideoService    services.UserVideoService                = services.NewUserVideoService(userVideoRepo)
	videoService        services.VideoService                    = services.NewVideoCategoriesService(videoRepo)
	verificationService services.VideoVerificationService        = services.NewVideoVerificationService(verificationRepository)
	videoController     controller.VideoController               = controller.NewVideoController(videoService, userVideoService, verificationService)
)

func VideoRouter(route *gin.Engine) {

	videoCategory := route.Group("/api/video-category", middleware.AuthorizeJWT(jwtService))
	{
		videoCategory.POST("/create-category", videoController.CreateCategory)
		videoCategory.PUT("/update-category", videoController.UpdateCategory)
		videoCategory.GET("/all-category", videoController.GetAllCategory)
		videoCategory.DELETE("/delete-category", videoController.DeleteCategory)
	}

	videos := route.Group("/api/videos", middleware.AuthorizeJWT(jwtService))
	{
		videos.POST("/add-video", videoController.AddVideo)
		videos.GET("/get-all-videos", videoController.GetAllVideos)
		videos.PUT("/update-video", videoController.UpdateVideo)
		videos.DELETE("/delete-video", videoController.DeleteVideo)
	}

	videoFullDetails := route.Group("/api/video-detail", middleware.AuthorizeJWT(jwtService))
	{
		videoFullDetails.GET("/video-full-details", videoController.VideoFullDetails)
	}
}

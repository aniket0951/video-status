package routers

import (
	"github.com/aniket0951/Chatrapati-Maharaj/controller"
	middleware "github.com/aniket0951/Chatrapati-Maharaj/middelware"
	"github.com/aniket0951/Chatrapati-Maharaj/repositories"
	"github.com/aniket0951/Chatrapati-Maharaj/services"
	"github.com/gin-gonic/gin"
)

var (
<<<<<<< HEAD
	videoRepo       repositories.VideoRepository = repositories.NewVideoCategoriesRepository()
	videoService    services.VideoService        = services.NewVideoCategoriesService(videoRepo, notificationManager)
	videoController controller.VideoController   = controller.NewVideoController(videoService)
=======
	videoRepo           = repositories.NewVideoCategoriesRepository()
	userVideoRepo       = repositories.NewUserVideoRepository()
	userVideoService    = services.NewUserVideoService(userVideoRepo)
	videoService        = services.NewVideoCategoriesService(videoRepo)
	verificationService = services.NewVideoVerificationService(verificationRepository)
	videoController     = controller.NewVideoController(videoService, userVideoService, verificationService)
>>>>>>> 9c19887285b2026e2c65966dca4df5157c7dfcd3
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
		videos.GET("/get-all-videos/:tag", videoController.GetAllVideos)
		videos.PUT("/update-video", videoController.UpdateVideo)
		videos.DELETE("/delete-video", videoController.DeleteVideo)
	}

<<<<<<< HEAD
	inActiveVideo := route.Group("/api")
	{
		inActiveVideo.GET("/inactive-video", videoController.FetchInActiveVideos)
		inActiveVideo.POST("/inactive-video/:videoId/:tag", videoController.ActiveVideo)

		inActiveVideo.GET("/video/:videoId", videoController.GetVideoByID)
	}

	download := route.Group("/api")
	{
		download.POST("/download-increase/:videoId", videoController.IncreaseDownloadCount)

		// share video link
		download.GET("/shared-video/:fileKey", videoController.GenerateSignVideoURL)
=======
	videoFullDetails := route.Group("/api/video-detail", middleware.AuthorizeJWT(jwtService))
	{
		videoFullDetails.GET("/video-full-details", videoController.VideoFullDetails)
>>>>>>> 9c19887285b2026e2c65966dca4df5157c7dfcd3
	}
}

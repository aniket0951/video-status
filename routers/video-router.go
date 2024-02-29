package routers

import (
	"github.com/aniket0951/Chatrapati-Maharaj/controller"
	middleware "github.com/aniket0951/Chatrapati-Maharaj/middelware"
	"github.com/aniket0951/Chatrapati-Maharaj/repositories"
	"github.com/aniket0951/Chatrapati-Maharaj/services"
	"github.com/gin-gonic/gin"
)

var (
	videoRepo       repositories.VideoRepository = repositories.NewVideoCategoriesRepository()
	videoService    services.VideoService        = services.NewVideoCategoriesService(videoRepo)
	videoController controller.VideoController   = controller.NewVideoController(videoService)
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

	inActiveVideo := route.Group("/api")
	{
		inActiveVideo.GET("/inactive-video", videoController.FetchInActiveVideos)
		inActiveVideo.POST("/inactive-video/:videoId", videoController.ActiveVideo)
		inActiveVideo.GET("/get-all-videos", videoController.GetAllVideos)
		inActiveVideo.GET("/video/:videoId", videoController.GetVideoByID)
	}

	download := route.Group("/api")
	{
		download.POST("/download-increase/:videoId", videoController.IncreaseDownloadCount)
	}
}

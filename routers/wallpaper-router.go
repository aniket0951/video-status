package routers

import (
	"github.com/aniket0951/Chatrapati-Maharaj/controller"
	middleware "github.com/aniket0951/Chatrapati-Maharaj/middelware"
	notificationmanager "github.com/aniket0951/Chatrapati-Maharaj/notification_manager"
	"github.com/aniket0951/Chatrapati-Maharaj/repositories"
	"github.com/aniket0951/Chatrapati-Maharaj/services"
	"github.com/gin-gonic/gin"
)

var notificationManager notificationmanager.NotificationManager

func WallPaperRouter(router *gin.Engine) {
	wallPaperRepo := repositories.NewWallPaperRepository()
	notificationManager = notificationmanager.NotificationManager{}
	wallPaperService := services.NewWallPaperService(wallPaperRepo, notificationManager)
	wallPaperController := controller.NewWallPaperController(wallPaperService)

	wallPaper := router.Group("/api", middleware.AuthorizeJWT(jwtService))
	{
		wallPaper.POST("/add-wallpaper", wallPaperController.AddWallPaper)
		wallPaper.GET("/get-wallpapers/:tag", wallPaperController.GetWallPapers)
		wallPaper.POST("/active-wallpaper/:id/:tag", wallPaperController.ActiveInActiveWallPaper)
		wallPaper.POST("/like-wallpapepr/:id", wallPaperController.WallPaperLiked)
	}
}

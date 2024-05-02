package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	// firebase "firebase.google.com/go"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/aniket0951/Chatrapati-Maharaj/helper"
	"github.com/aniket0951/Chatrapati-Maharaj/routers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

func CORSConfig() cors.Config {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers", "Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization")
	corsConfig.AddAllowMethods("GET", "POST", "PUT", "DELETE")
	return corsConfig
}

const (
	ContentTypeBinary = "application/octet-stream"
	ContentTypeForm   = "application/x-www-form-urlencoded"
	ContentTypeJSON   = "application/json"
	ContentTypeHTML   = "text/html; charset=utf-8"
	ContentTypeText   = "text/plain; charset=utf-8"
)

func main() {
	router := gin.New()
	router.Use(cors.New(CORSConfig()))

	router.Static("static", "./static")

	router.GET("/", func(ctx *gin.Context) {
		response := map[string]interface{}{}

		response["message"] = "Program run successfully..."

		ctx.JSON(http.StatusOK, response)

	})

	router.GET("/wallpaper/:dir/:file_name", func(ctx *gin.Context) {
		dir := ctx.Param("dir")
		file_name := ctx.Param("file_name")
		fmt.Printf("DIR %v and file name %v", dir, file_name)
		ctx.JSON(http.StatusOK, gin.H{"test": "aniket"})
	})

	router.GET("/send-notification", func(ctx *gin.Context) {
		// token := ctx.Param("token")
		token := "clHzLONMSFyLDEd4BT5piY:APA91bEHQ6zJajvmSAtuIkKfOZWvod_1iZcfexhMQ6-UzixIv8wDbA51rXWFICk5IQRKZNFuv4be3LKWQeg0RmKP-dobx51T9loiCm1RpNMNCQiftDN6iEIOJ9ma4FKM4La30wSezOze"
		err := generateToken(token)

		response := helper.BuildSuccessResponse(helper.FETCHED_SUCCESS, "TOKEN", err)
		ctx.JSON(http.StatusOK, response)
	})

	routers.UserAuthRouter(router)
	routers.VideoRouter(router)

	routers.WallPaperRouter(router)

	router.Run("0.0.0.0:8080")
}

func generateToken(fcm_token string) error {
	opt := option.WithCredentialsFile("./maharaj-fcm.json")

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
		return err
	}

	// Get the FCM client.
	client, err := app.Messaging(context.Background())
	if err != nil {
		log.Fatalf("error getting FCM client: %v\n", err)
		return err
	}

	// Create a new FCM message with the notification payload.
	message := &messaging.Message{
		Token: fcm_token,
		Notification: &messaging.Notification{
			Title:    "Jay Bhavani !",
			Body:     "New Video Has been uploaded",
			ImageURL: "http://192.168.0.109:5000/static/wallpaper/wallpaper-1019471600.png",
		},
	}

	// Send the message using the FCM client.
	response, err := client.Send(context.Background(), message)
	if err != nil {
		log.Fatalf("error sending message: %v\n", err)
		return err
	}

	routers.VideoVerificationRoute(router)

	// Print the message ID upon successful sending.
	fmt.Printf("Successfully sent message: %s\n", response)
	return nil
}


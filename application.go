package main

import (
	"net/http"

	"github.com/aniket0951/Chatrapati-Maharaj/routers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSConfig() cors.Config {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
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
	router := gin.Default()
	router.Use(cors.New(CORSConfig()))
	gin.SetMode(gin.ReleaseMode)
	router.Static("static", "static")

	router.GET("/", func(ctx *gin.Context) {
		response := map[string]interface{}{}

		response["message"] = "Program run successfully..."

		ctx.JSON(http.StatusOK, response)

	})

	routers.UserAuthRouter(router)
	routers.VideoRouter(router)
	routers.VideoVerificationRoute(router)

	router.Run(":5000")
}

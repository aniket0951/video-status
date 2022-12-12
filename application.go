package main

import (
	"net/http"

	"github.com/aniket0951/Chatrapati-Maharaj/routers"
	"github.com/gin-gonic/gin"
)

const (
	ContentTypeBinary = "application/octet-stream"
	ContentTypeForm   = "application/x-www-form-urlencoded"
	ContentTypeJSON   = "application/json"
	ContentTypeHTML   = "text/html; charset=utf-8"
	ContentTypeText   = "text/plain; charset=utf-8"
)

func main() {
	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		response := map[string]interface{}{}

		response["message"] = "Program run successfully..."

		ctx.JSON(http.StatusOK, response)

	})

	routers.UserAuthRouter(router)
	routers.VideoRouter(router)

	router.Run(":5000")
}

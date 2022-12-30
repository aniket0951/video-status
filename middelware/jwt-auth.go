package middleware

import (
	"fmt"
	"net/http"

	"github.com/aniket0951/Chatrapati-Maharaj/helper"
	"github.com/aniket0951/Chatrapati-Maharaj/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthorizeJWT(jwtService services.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			response := helper.BuildFailedResponse(helper.FAILED_PROCESS, "Token not found", "authentication", helper.EmptyObj{})
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		token, err := jwtService.ValidateToken(authHeader)

		if err != nil {
			response := helper.BuildFailedResponse("Invalid token provided !", err.Error(), "authentication", helper.EmptyObj{})
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			fmt.Println("claim user id is", claims["user_id"])
			str := fmt.Sprintf("%v", claims["user_id"])
			helper.USER_ID = str
		}

	}
}

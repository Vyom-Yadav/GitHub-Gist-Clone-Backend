package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/initializers"
	"github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/models"
	"github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func DeserializeUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var accessToken string
		cookie, err := ctx.Cookie("access_token")

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			accessToken = fields[1]
		} else if err == nil {
			accessToken = cookie
		}

		if accessToken == "" {
			statusCode := http.StatusUnauthorized
			ctx.AbortWithStatusJSON(statusCode, models.ErrorResponseWrapper{
				Error: models.ErrorResponse{
					StatusCode: statusCode,
					Message:    "You are not logged in",
				},
			})
			return
		}

		config, _ := initializers.LoadConfig(os.Getenv("API_ENV_CONFIG_PATH"))
		sub, err := utils.ValidateToken(accessToken, config.AccessTokenPublicKey)
		if err != nil {
			statusCode := http.StatusUnauthorized
			ctx.AbortWithStatusJSON(statusCode, models.ErrorResponseWrapper{
				Error: models.ErrorResponse{
					StatusCode: statusCode,
					Message:    err.Error(),
				},
			})
			return
		}

		var user models.User
		result := initializers.DB.Preload("UserMetadata").Preload("Gists").Preload("Gists.GistContent").First(&user, "username = ?", fmt.Sprint(sub))
		if result.Error != nil {
			statusCode := http.StatusForbidden
			ctx.AbortWithStatusJSON(statusCode, models.ErrorResponseWrapper{
				Error: models.ErrorResponse{
					StatusCode: statusCode,
					Message:    "the user belonging to this token no logger exists",
				},
			})
			zap.L().Error(result.Error.Error())
			return
		}

		ctx.Set("currentUser", user)
		ctx.Next()
	}
}

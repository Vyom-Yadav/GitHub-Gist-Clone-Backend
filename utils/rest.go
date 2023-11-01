package utils

import (
	"net/http"

	"github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/models"
	"github.com/gin-gonic/gin"
)

func SomethingBadHappened(ctx *gin.Context) {
	NewErrorResponse(ctx, http.StatusInternalServerError, "Something bad happened")
}

func NewErrorResponse(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, models.ErrorResponseWrapper{
		Error: models.ErrorResponse{
			Message:    message,
			StatusCode: statusCode,
		},
	})
}

func NewSuccessResponse(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, models.SuccessResponseWrapper{
		Success: models.SuccessResponse{
			Message:    message,
			StatusCode: statusCode,
		},
	})
}

func NewAccessCodeResponse(ctx *gin.Context, statusCode int, accessCode string) {
	ctx.JSON(statusCode, models.AccessCodeResponseWrapper{
		AccessCode: models.AccessCodeResponse{
			AccessCode: accessCode,
		},
	})
}

package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SomethingBadHappened(ctx *gin.Context) {
	ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Something bad happened"})
}

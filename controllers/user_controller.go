package controllers

import (
	"net/http"

	"github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(DB *gorm.DB) UserController {
	return UserController{
		DB: DB,
	}
}

//	@Summary	Get the current logged in user details, panics if not logged in.
//	@Tags		User Operations
//	@Produce	json
//	@Success	200	{object}	map[string]string
//	@Router		/users/me [get]
func (uc *UserController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	userResponse := &models.UserResponse{
		Username:     currentUser.Username,
		FirstName:    currentUser.FirstName,
		LastName:     currentUser.LastName,
		Email:        currentUser.Email,
		Role:         currentUser.Role,
		Provider:     currentUser.Provider,
		UserMetadata: currentUser.UserMetadata,
		CreatedAt:    currentUser.CreatedAt,
		UpdatedAt:    currentUser.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": userResponse}})
}

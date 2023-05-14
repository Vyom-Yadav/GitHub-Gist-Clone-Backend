package controllers

import (
	"net/http"
	"time"

	"github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/models"
	"github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

//
// ------------- READ FUNCTIONS -----------------------
//

//	@Summary	Get the current logged in user details.
//	@Tags		User Operations
//	@Produce	json
//	@Success	200	{object}	map[string]any
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
		Gists:        currentUser.Gists,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": userResponse}})
}

//	@Summary	Get the publicly visible details of a user, does not load gists
//	@Tags		User Operations
//	@Produce	json
//	@Param		username	path		string	true	"The username to get"
//	@Success	200			{object}	map[string]any
//	@Failure	404			{object}	map[string]string
//	@Router		/users/{username} [get]
func (uc *UserController) GetUser(ctx *gin.Context) {
	username := ctx.Params.ByName("username")

	var user models.User
	result := uc.DB.Preload("UserMetadata").First(&user, "username = ?", username)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "user with username: '" + username + "' does not exist"})
		return
	}

	gists := make([]models.Gist, 0)
	for _, gist := range user.Gists {
		if !gist.Private {
			gists = append(gists, gist)
		}
	}

	publicUserProfile := models.PublicUserProfileResponse{
		Username:     user.Username,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		UserMetadata: user.UserMetadata,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": publicUserProfile}})
}

//	@Summary	Get the publicly visible gists of a user, does not load the gist comments
//	@Tags		User Operations
//	@Produce	json
//	@Param		username	path		string	true	"The username to get gists for"
//	@Success	200			{object}	map[string]any
//	@Failure	404			{object}	map[string]string
//	@Router		/users/{username}/gists [get]
func (uc *UserController) GetUserGists(ctx *gin.Context) {
	username := ctx.Params.ByName("username")

	var user models.User
	result := uc.DB.
		Preload("Gists").
		Preload("Gists.GistContent").
		First(&user, "username = ?", username)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "user with username: '" + username + "' does not exist"})
		return
	}

	gists := make([]models.Gist, 0)
	for _, gist := range user.Gists {
		if !gist.Private {
			gists = append(gists, gist)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"username": username, "gists": gists}})
}

//	@Summary	Get the publicly visible gist Ids of a user
//	@Tags		User Operations
//	@Produce	json
//	@Param		username	path		string	true	"The username to get gists for"
//	@Success	200			{object}	map[string]any
//	@Failure	404			{object}	map[string]string
//	@Router		/users/{username}/gistIds [get]
func (uc *UserController) GetUserGistsIds(ctx *gin.Context) {
	username := ctx.Params.ByName("username")

	var user models.User
	result := uc.DB.
		Preload("Gists").
		First(&user, "username = ?", username)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "user with username: '" + username + "' does not exist"})
		return
	}

	gistIds := make([]uuid.UUID, 0)
	for _, gist := range user.Gists {
		if !gist.Private {
			gistIds = append(gistIds, gist.ID)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"username": username, "gistIds": gistIds}})
}

//
// ------------- CREATE FUNCTIONS -----------------------
//

//	@Summary	Create a gist
//	@Tags		User Operations
//	@Accept		json
//	@Produce	json
//	@Param		CreateGistInput	body		models.CreateGistRequest	true	"The Input for creating gist"
//	@Success	201				{object}	map[string]any
//	@Failure	400				{object}	map[string]string
//	@Router		/users/gists [post]
func (uc *UserController) CreateGist(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload *models.CreateGistRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	now := time.Now()

	gists := currentUser.Gists
	for _, gist := range gists {
		if gist.Name == payload.Name {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Gist with name: '" + payload.Name + "' already exists"})
			return
		}
	}

	// TODO: There is some problem with content, check that
	newGist := models.Gist{
		Username:    currentUser.Username,
		Private:     payload.Private,
		GistContent: models.GistContent{
			Content: payload.Content,
		},
		Name:        payload.Name,
		Title:       payload.Title,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	result := uc.DB.Session(&gorm.Session{FullSaveAssociations: true}).Create(&newGist)
	if result.Error != nil {
		// Fuck it, live with this only
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newGist})
}

//	@Summary	Create a comment on a gist
//	@Tags		User Operations
//	@Accept		json
//	@Produce	json
//	@Param		CreateCommentInput	body		models.CommentOnGistRequest	true	"The Input for creating comment"
//	@Success	201					{object}	map[string]any
//	@Failure	400					{object}	map[string]string
//	@Router		/users/comments [post]
func (uc *UserController) CreateCommentOnGist(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload *models.CommentOnGistRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	now := time.Now()

	gistUUID, err := uuid.FromBytes([]byte(payload.GistId))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	newComment := models.Comment{
		GistID:    gistUUID,
		Username:  currentUser.Username,
		Content:   payload.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}
	result := uc.DB.Create(&newComment)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newComment})
}

// TODO: Fork a gist, remember to check if name is already taken

//
// ------------- UPDATE FUNCTIONS -----------------------
//

//	@Summary	Update user metadata
//	@Tags		User Operations
//	@Accept		json
//	@Produce	json
//	@Param		UpdateUserDetailsInput	body		models.UpdateUserDetailsRequest	true	"The Input for updating user metadata"
//	@Success	200						{object}	map[string]any
//	@Failure	400						{object}	map[string]string
//	@Router		/users/details [patch]
func (uc *UserController) UpdateUserDetails(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload *models.UpdateUserDetailsRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	userMetadata := currentUser.UserMetadata

	if payload.ProfilePicture != "" {
		userMetadata.ProfilePicture = payload.ProfilePicture
	}
	if payload.Tagline != "" {
		userMetadata.Tagline = payload.Tagline
	}
	if payload.StatusIcon != "" {
		userMetadata.StatusIcon = payload.StatusIcon
	}
	if payload.Location != "" {
		userMetadata.Location = payload.Location
	}
	if payload.Website != "" {
		userMetadata.Website = payload.Website
	}
	if payload.Twitter != "" {
		userMetadata.Twitter = payload.Twitter
	}

	result := uc.DB.Save(&userMetadata)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": userMetadata})
}

//	@Summary	Update gist data
//	@Tags		User Operations
//	@Accept		json
//	@Produce	json
//	@Param		UpdateGistInput	body		models.UpdateGistRequest	true	"The Input for updating user gist"
//	@Success	200				{object}	map[string]any
//	@Failure	400				{object}	map[string]string
//	@Failure	401				{object}	map[string]string
//	@Failure	404				{object}	map[string]string
//	@Router		/users/gists [patch]
func (uc *UserController) UpdateGist(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload *models.UpdateGistRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	var gist models.Gist
	result := uc.DB.
		Preload("GistContent").
		First(&gist, "id = ?", payload.GistId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "gist does not exist"})
		return
	}

	if gist.Username != currentUser.Username {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "unauthorized"})
		return
	}

	if payload.Name != "" {
		currentUserGists := currentUser.Gists
		for _, currentUserGist := range currentUserGists {
			if currentUserGist.Name == payload.Name {
				ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Gist with name: '" + payload.Name + "' already exists"})
				return
			}
		}
		gist.Name = payload.Name
	}
	if payload.Title != "" {
		gist.Title = payload.Title
	}
	if payload.Content != "" {
		gist.GistContent.Content = payload.Content
	}
	gist.Private = payload.Private
	gist.UpdatedAt = time.Now()

	result = uc.DB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&gist)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gist})
}

//	@Summary	Follow a user
//	@Tags		User Operations
//	@Produce	json
//	@Param		userToFollow	path		string	true	"The ID of the user to follow"
//	@Success	200				{object}	map[string]any
//	@Failure	400				{object}	map[string]string
//	@Failure	404				{object}	map[string]string
//	@Failure	500				{object}	map[string]string
//	@Router		/users/follow/{userToFollow} [patch]
func (uc *UserController) FollowUser(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	userToFollow := ctx.Params.ByName("userToFollow")

	if currentUser.Username == userToFollow {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "You cannot follow yourself"})
		return
	}

	var userToBeFollowed models.User
	result := uc.DB.Preload("UserMetadata").First(&userToBeFollowed, "username = ?", userToFollow)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "user does not exist"})
		return
	}

	// Perform transaction to update both users
	err := uc.DB.Transaction(func(tx *gorm.DB) error {
		// Update current user
		currentUserMetadata := currentUser.UserMetadata
		currentUserMetadata.Following = append(currentUserMetadata.Following, userToBeFollowed.Username)
		result := tx.Save(&currentUserMetadata)
		if result.Error != nil {
			return result.Error
		}

		// Update user to be followed
		userToBeFollowedMetadata := userToBeFollowed.UserMetadata
		userToBeFollowedMetadata.Followers = append(userToBeFollowedMetadata.Followers, currentUser.Username)
		result = tx.Save(&userToBeFollowedMetadata)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "successfully followed user"})
}

//	@Summary	Unfollow a user
//	@Tags		User Operations
//	@Produce	json
//	@Param		userToUnfollow	path		string	true	"The ID of the user to unfollow"
//	@Success	200				{object}	map[string]any
//	@Failure	400				{object}	map[string]string
//	@Failure	404				{object}	map[string]string
//	@Failure	500				{object}	map[string]string
//	@Router		/users/follow/{userToUnfollow} [patch]
func (uc *UserController) UnfollowUser(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	userToUnfollow := ctx.Params.ByName("userToUnfollow")

	if currentUser.Username == userToUnfollow {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "You cannot unfollow yourself"})
		return
	}

	var userToBeUnfollowed models.User
	result := uc.DB.Preload("UserMetadata").First(&userToBeUnfollowed, "username = ?", userToUnfollow)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "user does not exist"})
		return
	}

	// Perform transaction to update both users
	err := uc.DB.Transaction(func(tx *gorm.DB) error {
		// Update current user
		currentUserMetadata := currentUser.UserMetadata
		currentUserMetadata.Following = utils.RemoveStringFromSlice(currentUserMetadata.Following, userToBeUnfollowed.Username)
		result := tx.Save(&currentUserMetadata)
		if result.Error != nil {
			return result.Error
		}

		// Update user to be unfollowed
		userToBeUnfollowedMetadata := userToBeUnfollowed.UserMetadata
		userToBeUnfollowedMetadata.Followers = utils.RemoveStringFromSlice(userToBeUnfollowedMetadata.Followers, currentUser.Username)
		result = tx.Save(&userToBeUnfollowedMetadata)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "successfully unfollowed user"})
}

//	@Summary	Star a gist
//	@Tags		User Operations
//	@Produce	json
//	@Param		gistId	path		string	true	"The ID of the gist to star"
//	@Success	200		{object}	map[string]any
//	@Failure	404		{object}	map[string]string
//	@Failure	500		{object}	map[string]string
//	@Router		/users/gists/{gistId}/star [patch]
func (uc *UserController) StarGist(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	gistId := ctx.Params.ByName("gistId")

	var gist models.Gist
	result := uc.DB.First(&gist, "id = ?", gistId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "gist does not exist"})
		return
	}

	// Perform transaction to update both users
	err := uc.DB.Transaction(func(tx *gorm.DB) error {
		// Update current user
		currentUserMetadata := currentUser.UserMetadata
		currentUserMetadata.StarredGists = append(currentUserMetadata.StarredGists, gist.ID.String())
		result := tx.Save(&currentUserMetadata)
		if result.Error != nil {
			return result.Error
		}

		// Update gist
		gist.Stars = append(gist.Stars, currentUser.Username)
		result = tx.Save(&gist)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "successfully starred gist"})
}

//	@Summary	Un-star a gist
//	@Tags		User Operations
//	@Produce	json
//	@Param		gistId	path		string	true	"The ID of the gist to un-star"
//	@Success	200		{object}	map[string]any
//	@Failure	404		{object}	map[string]string
//	@Failure	500		{object}	map[string]string
//	@Router		/users/gists/{gistId}/unstar [patch]
func (uc *UserController) UnStarGist(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	gistId := ctx.Params.ByName("gistId")

	var gist models.Gist
	result := uc.DB.First(&gist, "id = ?", gistId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "gist does not exist"})
		return
	}

	// Perform transaction to update both users
	err := uc.DB.Transaction(func(tx *gorm.DB) error {
		// Update current user
		currentUserMetadata := currentUser.UserMetadata
		currentUserMetadata.StarredGists = utils.RemoveStringFromSlice(currentUserMetadata.StarredGists, gist.ID.String())
		result := tx.Save(&currentUserMetadata)
		if result.Error != nil {
			return result.Error
		}

		// Update gist
		gist.Stars = utils.RemoveStringFromSlice(gist.Stars, currentUser.Username)
		result = tx.Save(&gist)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "successfully unstarred gist"})
}


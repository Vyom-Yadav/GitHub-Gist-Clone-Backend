package controllers

import (
	"net/http"
	"time"

	"github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
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
//	@Failure	401	{object}	models.ErrorResponseWrapper
//	@Failure	403	{object}	models.ErrorResponseWrapper
//	@Router		/users/me [get]
func (uc *UserController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	userResponse := &models.UserResponse{
		Username:     currentUser.Username,
		FirstName:    currentUser.FirstName,
		Email:        currentUser.Email,
		Role:         currentUser.Role,
		Provider:     currentUser.Provider,
		UserMetadata: currentUser.UserMetadata,
		CreatedAt:    currentUser.CreatedAt,
		UpdatedAt:    currentUser.UpdatedAt,
		Gists:        currentUser.Gists,
	}

	if currentUser.LastName != nil {
		userResponse.LastName = *currentUser.LastName
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
		UserMetadata: user.UserMetadata,
		Verified:     user.Verified,
	}

	if user.LastName != nil {
		publicUserProfile.LastName = *user.LastName
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
//	@Failure	401				{object}	models.ErrorResponseWrapper
//	@Failure	403				{object}	models.ErrorResponseWrapper
//	@Router		/users/gists [post]
func (uc *UserController) CreateGist(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload *models.CreateGistRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	now := time.Now()

	gists := currentUser.Gists
	for _, gist := range gists {
		if gist.Name == payload.Name {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Gist with name: '" + payload.Name + "' already exists"})
			return
		}
	}

	// TODO: There is some problem with content, check that
	newGist := models.Gist{
		Username: currentUser.Username,
		Private:  payload.Private,
		GistContent: models.GistContent{
			Content: payload.Content,
		},
		Name:      payload.Name,
		Title:     payload.Title,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := uc.DB.Session(&gorm.Session{FullSaveAssociations: true}).Create(&newGist)
	if result.Error != nil {
		// Fuck it, live with this only
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": result.Error.Error()})
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
//	@Failure	401					{object}	models.ErrorResponseWrapper
//	@Failure	403					{object}	models.ErrorResponseWrapper
//	@Router		/users/comments [post]
func (uc *UserController) CreateCommentOnGist(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload *models.CommentOnGistRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	now := time.Now()

	gistUUID, err := uuid.Parse(payload.GistId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
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
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": result.Error.Error()})
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
//	@Failure	401						{object}	models.ErrorResponseWrapper
//	@Failure	403						{object}	models.ErrorResponseWrapper
//	@Router		/users/details [patch]
func (uc *UserController) UpdateUserDetails(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload *models.UpdateUserDetailsRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	userMetadata := currentUser.UserMetadata

	if payload.ProfilePicture != "" {
		userMetadata.ProfilePicture = payload.ProfilePicture
	}
	if payload.Tagline != "" {
		userMetadata.Tagline = &payload.Tagline
	}
	if payload.StatusIcon != "" {
		userMetadata.StatusIcon = payload.StatusIcon
	}
	if payload.Location != "" {
		userMetadata.Location = &payload.Location
	}
	if payload.Website != "" {
		userMetadata.Website = &payload.Website
	}
	if payload.Twitter != "" {
		userMetadata.Twitter = &payload.Twitter
	}

	result := uc.DB.Save(&userMetadata)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": result.Error.Error()})
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
//	@Failure	401				{object}	models.ErrorResponseWrapper
//	@Failure	403				{object}	models.ErrorResponseWrapper
//	@Failure	404				{object}	map[string]string
//	@Router		/users/gists [patch]
func (uc *UserController) UpdateGist(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload *models.UpdateGistRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
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
				ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Gist with name: '" + payload.Name + "' already exists"})
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
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gist})
}

//	@Summary	Follow a user
//	@Tags		User Operations
//	@Produce	json
//	@Param		userToFollow	path		string	true	"The username of the user to follow"
//	@Success	200				{object}	map[string]any
//	@Failure	400				{object}	map[string]string
//	@Failure	404				{object}	map[string]string
//	@Failure	401				{object}	models.ErrorResponseWrapper
//	@Failure	403				{object}	models.ErrorResponseWrapper
//	@Failure	500				{object}	map[string]string
//	@Router		/users/follow/{userToFollow} [patch]
func (uc *UserController) FollowUser(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	userToFollow := ctx.Params.ByName("userToFollow")

	if currentUser.Username == userToFollow {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "You cannot follow yourself"})
		return
	}

	var userToBeFollowed models.User
	result := uc.DB.Preload("UserMetadata").First(&userToBeFollowed, "username = ?", userToFollow)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "user does not exist"})
		return
	}

	// Perform transaction to update both users
	err := uc.DB.Transaction(func(tx *gorm.DB) error {
		// Update current user
		currentUserMetadata := currentUser.UserMetadata
		currentUserMetadata.Following += 1
		result := tx.Save(&currentUserMetadata)
		if result.Error != nil {
			zap.L().Error(result.Error.Error())
			return result.Error
		}

		// Update user to be followed
		userToBeFollowedMetadata := userToBeFollowed.UserMetadata
		userToBeFollowedMetadata.Followers += 1
		result = tx.Save(&userToBeFollowedMetadata)
		if result.Error != nil {
			zap.L().Error(result.Error.Error())
			return result.Error
		}

		newFollow := models.Follow{
			Username: userToBeFollowed.Username,
			FollowedBy: currentUser.Username,
		}

		result = tx.Create(&newFollow)
		if result.Error != nil {
			zap.L().Error(result.Error.Error())
			return result.Error
		}

		return nil
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "successfully followed user"})
}

//	@Summary	Unfollow a user
//	@Tags		User Operations
//	@Produce	json
//	@Param		userToUnfollow	path		string	true	"The username of the user to unfollow"
//	@Success	200				{object}	map[string]any
//	@Failure	400				{object}	map[string]string
//	@Failure	404				{object}	map[string]string
//	@Failure	401				{object}	models.ErrorResponseWrapper
//	@Failure	403				{object}	models.ErrorResponseWrapper
//	@Failure	500				{object}	map[string]string
//	@Router		/users/unfollow/{userToUnfollow} [patch]
func (uc *UserController) UnfollowUser(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	userToUnfollow := ctx.Params.ByName("userToUnfollow")

	if currentUser.Username == userToUnfollow {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "You cannot unfollow yourself"})
		return
	}

	var userToBeUnfollowed models.User
	result := uc.DB.Preload("UserMetadata").First(&userToBeUnfollowed, "username = ?", userToUnfollow)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "user does not exist"})
		return
	}

	// Perform transaction to update both users
	err := uc.DB.Transaction(func(tx *gorm.DB) error {
		// Update current user
		currentUserMetadata := currentUser.UserMetadata
		currentUserMetadata.Following -= 1
		if currentUserMetadata.Following < 0 {
			currentUserMetadata.Following = 0
		}
		result := tx.Save(&currentUserMetadata)
		if result.Error != nil {
			zap.L().Error(result.Error.Error())
			return result.Error
		}

		// Update user to be unfollowed
		userToBeUnfollowedMetadata := userToBeUnfollowed.UserMetadata
		userToBeUnfollowedMetadata.Followers -= 1
		if userToBeUnfollowedMetadata.Followers < 0 {
			userToBeUnfollowedMetadata.Followers = 0
		}
		result = tx.Save(&userToBeUnfollowedMetadata)
		if result.Error != nil {
			zap.L().Error(result.Error.Error())
			return result.Error
		}

		followToBeDeleted := models.Follow{
			Username: userToBeUnfollowed.Username,
			FollowedBy: currentUser.Username,
		}

		result = tx.Delete(&followToBeDeleted)
		if result.Error != nil {
			zap.L().Error(result.Error.Error())
			return result.Error
		}

		return nil
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
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
//	@Failure	401		{object}	models.ErrorResponseWrapper
//	@Failure	403		{object}	models.ErrorResponseWrapper
//	@Router		/users/gists/{gistId}/star [patch]
func (uc *UserController) StarGist(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	gistId := ctx.Params.ByName("gistId")

	parsedGistId, err := uuid.Parse(gistId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "invalid gist id"})
		return
	}

	var gist models.Gist
	result := uc.DB.First(&gist, "id = ?", parsedGistId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "gist does not exist"})
		return
	}

	// Perform transaction to update both users
	err = uc.DB.Transaction(func(tx *gorm.DB) error {
		// Update current user
		currentUserMetadata := currentUser.UserMetadata
		currentUserMetadata.StarredGistsCount += 1
		result := tx.Save(&currentUserMetadata)
		if result.Error != nil {
			zap.L().Error(result.Error.Error())
			return result.Error
		}

		// Update gist
		// Read https://www.postgresql.org/docs/9.1/arrays.html#ARRAYS-INPUT, well just 'cause you should know it
		gist.StarCount += 1
		result = tx.Save(&gist)
		if result.Error != nil {
			zap.L().Error(result.Error.Error())
			return result.Error
		}

		newStarredGist := models.Star{
			GistID: gist.ID,
			Username: currentUser.Username,
		}

		result = tx.Create(&newStarredGist)
		if result.Error != nil {
			zap.L().Error(result.Error.Error())
			return result.Error
		}

		return nil
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
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
//	@Failure	401		{object}	models.ErrorResponseWrapper
//	@Failure	403		{object}	models.ErrorResponseWrapper
//	@Router		/users/gists/{gistId}/unstar [patch]
func (uc *UserController) UnstarGist(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	gistId := ctx.Params.ByName("gistId")
	parsedGistId, err := uuid.Parse(gistId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "invalid gist id"})
		return
	}

	var gist models.Gist
	result := uc.DB.First(&gist, "id = ?", parsedGistId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "gist does not exist"})
		return
	}

	// Perform transaction to update both users
	err = uc.DB.Transaction(func(tx *gorm.DB) error {
		// Update current user
		currentUserMetadata := currentUser.UserMetadata
		currentUserMetadata.StarredGistsCount -= 1
		if currentUserMetadata.StarredGistsCount < 0 {
			currentUserMetadata.StarredGistsCount = 0
		}
		result := tx.Save(&currentUserMetadata)
		if result.Error != nil {
			zap.L().Error(result.Error.Error())
			return result.Error
		}

		// Update gist
		gist.StarCount -= 1
		if gist.StarCount < 0 {
			gist.StarCount = 0
		}
		result = tx.Save(&gist)
		if result.Error != nil {
			zap.L().Error(result.Error.Error())
			return result.Error
		}

		starToDelete := models.Star{
			GistID: gist.ID,
			Username: currentUser.Username,
		}

		result = tx.Delete(&starToDelete)
		if result.Error != nil {
			zap.L().Error(result.Error.Error())
			return result.Error
		}

		return nil
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "successfully unstarred gist"})
}

//	@Summary	Get the followers of a user
//	@Tags		User Operations
//	@Produce	json
//	@Param		username	path		string	true	"The username of the user to get the followers of"
//	@Success	200			{object}	map[string]any
//	@Failure	404			{object}	map[string]string
//	@Failure	500			{object}	map[string]string
//	@Router		/users/{username}/followers [get]
func (uc *UserController) GetFollowerList(ctx *gin.Context) {
	username := ctx.Params.ByName("username")

	var user models.User
	result := uc.DB.First(&user, "username = ?", username)
	if result.Error != nil {
		zap.L().Error(result.Error.Error())
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "user does not exist"})
		return
	}

	var followers []models.Follow
	result = uc.DB.Find(&followers, "username = ?", username)
	if result.Error != nil {
		zap.L().Error(result.Error.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}

	var followerUsernames []string
	for _, follower := range followers {
		followerUsernames = append(followerUsernames, follower.FollowedBy)
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "followers": followerUsernames})
}

//	@Summary	Get the list of users a user is following
//	@Tags		User Operations
//	@Produce	json
//	@Param		username	path		string	true	"The username of the user to get the following of"
//	@Success	200			{object}	map[string]any
//	@Failure	404			{object}	map[string]string
//	@Failure	500			{object}	map[string]string
//	@Router		/users/{username}/following [get]
func (uc *UserController) GetFollowingList(ctx *gin.Context) {
	username := ctx.Params.ByName("username")

	var user models.User
	result := uc.DB.First(&user, "username = ?", username)
	if result.Error != nil {
		zap.L().Error(result.Error.Error())
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "user does not exist"})
		return
	}

	var following []models.Follow
	result = uc.DB.Find(&following, "followed_by = ?", username)
	if result.Error != nil {
		zap.L().Error(result.Error.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}

	var followingUsernames []string
	for _, followedUser := range following {
		followingUsernames = append(followingUsernames, followedUser.Username)
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "following": followingUsernames})
}

//	@Summary	Whether a username follows another username
//	@Tags		User Operations
//	@Produce	json
//	@Param		username	path		string	true	"The username of the follower"
//	@Param		otherUser	path		string	true	"The username of the user being followed"
//	@Success	200			{object}	map[string]any
//	@Failure	404			{object}	map[string]any
//	@Router		/users/{username}/follows/{otherUser} [get]
func (uc *UserController) CheckIfUserFollows(ctx *gin.Context) {
	username := ctx.Params.ByName("username")
	otherUser := ctx.Params.ByName("otherUser")

	var follow models.Follow
	result := uc.DB.First(&follow, "username = ? AND followed_by = ?", otherUser, username)
	if result.Error != nil {
		zap.L().Error(result.Error.Error())
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "follows": false})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "follows": true})
}

//	@Summary	Get the list of starred gists of a user
//	@Tags		User Operations
//	@Produce	json
//	@Param		username	path		string	true	"The username of the user to get the starred gists of"
//	@Success	200			{object}	map[string]any
//	@Failure	404			{object}	map[string]string
//	@Failure	500			{object}	map[string]string
//	@Router		/users/{username}/starredGists [get]
func (uc *UserController) GetStarredGists(ctx *gin.Context) {
	username := ctx.Params.ByName("username")

	var user models.User
	result := uc.DB.First(&user, "username = ?", username)
	if result.Error != nil {
		zap.L().Error(result.Error.Error())
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "user does not exist"})
		return
	}

	var stars []models.Star
	result = uc.DB.Find(&stars, "username = ?", username)
	if result.Error != nil {
		zap.L().Error(result.Error.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}

	var starredGistIds []uuid.UUID
	for _, star := range stars {
		starredGistIds = append(starredGistIds, star.GistID)
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "starred_gists": starredGistIds})
}

//	@Summary	Whether a user has starred a gist
//	@Tags		User Operations
//	@Produce	json
//	@Param		username	path		string	true	"The username of the user to check the starred gist of"
//	@Param		gistId		path		string	true	"The ID of the gist to check if it is starred"
//	@Success	200			{object}	map[string]any
//	@Failure	400			{object}	map[string]string
//	@Failure	404			{object}	map[string]any
//	@Router		/users/{username}/starredGist/{gistId} [get]
func (uc *UserController) CheckIfGistStarred(ctx *gin.Context) {
	username := ctx.Params.ByName("username")
	gistId := ctx.Params.ByName("gistId")

	parsedGistId, err := uuid.Parse(gistId)
	if err != nil {
		zap.L().Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "invalid gist id"})
		return
	}

	var star models.Star
	result := uc.DB.First(&star, "username = ? AND gist_id = ?", username, parsedGistId)
	if result.Error != nil {
		zap.L().Error(result.Error.Error())
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "starred": false})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "starred": true})
}


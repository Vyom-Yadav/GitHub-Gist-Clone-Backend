package routes

import (
	"github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/controllers"
	"github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/middleware"
	"github.com/gin-gonic/gin"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewUserRouteController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup) {
	router := rg.Group("users")

	router.GET("/me", middleware.DeserializeUser(), uc.userController.GetMe)
	router.GET("/:username", uc.userController.GetUser)
	router.GET("/:username/gists", uc.userController.GetUserGists)
	router.GET("/:username/gistIds", uc.userController.GetUserGistsIds)

	router.POST("/gists", middleware.DeserializeUser(), uc.userController.CreateGist)
	router.POST("/comments", middleware.DeserializeUser(), uc.userController.CreateCommentOnGist)

	router.PATCH("/details", middleware.DeserializeUser(), uc.userController.UpdateUserDetails)
	router.PATCH("/gists", middleware.DeserializeUser(), uc.userController.UpdateGist)
	router.PATCH("follow/:userToFollow", middleware.DeserializeUser(), uc.userController.FollowUser)
	router.PATCH("unfollow/:userToUnfollow", middleware.DeserializeUser(), uc.userController.UnfollowUser)
	router.PATCH("gists/:gistId/star", middleware.DeserializeUser(), uc.userController.StarGist)
	router.PATCH("gists/:gistId/unstar", middleware.DeserializeUser(), uc.userController.UnstarGist)

	router.GET("/:username/followers", uc.userController.GetFollowerList)
	router.GET("/:username/following", uc.userController.GetFollowingList)
	router.GET("/:username/starredGists", uc.userController.GetStarredGists)
	router.GET("/:username/follows/:otherUser", uc.userController.CheckIfUserFollows)
	router.GET("/:username/starredGist/:gistId", uc.userController.CheckIfGistStarred)
}

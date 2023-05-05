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
}

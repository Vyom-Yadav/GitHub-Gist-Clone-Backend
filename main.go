package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/controllers"
	"github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/docs"
	"github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/initializers"
	"github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/models"
	"github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	server              *gin.Engine
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	GistController      controllers.GistController
	GistRouteController routes.GistRouteController
)

func init() {
	config, err := initializers.LoadConfig("/app/env")
	if err != nil {
		log.Fatal("Could not load environment variables ", err)
	}

	initializers.ConnectDB(&config)

	err = initializers.DB.AutoMigrate(
		&models.User{},
		&models.UserMetadata{},
		&models.Gist{},
		&models.Comment{},
		&models.GistContent{},
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Migration complete")

	AuthController = controllers.NetAuthController(initializers.DB)
	UserController = controllers.NewUserController(initializers.DB)
	GistController = controllers.NewGistController(initializers.DB)

	AuthRouteController = routes.NewAuthRouteController(AuthController)
	UserRouteController = routes.NewUserRouteController(UserController)
	GistRouteController = routes.NewGistRouteController(GistController)

	server = gin.Default()
}

//	@title			GitHub Gist Backend REST API
//	@version		1.0-alpha
//	@description	The REST API for GitHub Gist Backend

//	@BasePath	/api/
func main() {
	config, err := initializers.LoadConfig("/app/env")
	if err != nil {
		log.Fatal("Could not load environment variables ", err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:" + config.ServerPort, config.ClientOrigin}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/health", healthHandler)

	docs.SwaggerInfo.Host = "localhost:" + config.ServerPort
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	AuthRouteController.AuthRoute(router)
	UserRouteController.UserRoute(router)
	GistRouteController.GistRoute(router)
	log.Fatal(server.Run(":" + config.ServerPort))
}

//	@Summary	Check the basic health of api
//	@Tags		Health
//	@Produce	json
//	@Success	200	{object}	map[string]string
//	@Router		/health [get]
func healthHandler(ctx *gin.Context) {
	message := "Fuck Off, I am working!"
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
}

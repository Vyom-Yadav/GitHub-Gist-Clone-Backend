package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/controllers"
	"github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/docs"
	"github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/initializers"
	"github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/models"
	"github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
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
	// Volume mapping in docker container ./app.env:/app/env/app.env
	config, err := initializers.LoadConfig(os.Getenv("API_ENV_CONFIG_PATH"))
	if err != nil {
		log.Fatal("Could not load environment variables ", err)
	}

	var logger *zap.Logger
	if config.AppEnv == "production" {
		logger = zap.Must(zap.NewProduction())
	} else {
		logger = zap.Must(zap.NewDevelopment())
	}

	zap.ReplaceGlobals(logger)

	initializers.ConnectDB(&config)

	err = initializers.DB.AutoMigrate(
		&models.User{},
		&models.UserMetadata{},
		&models.Gist{},
		&models.Comment{},
		&models.GistContent{},
		&models.Follow{},
		&models.Star{},
	)
	if err != nil {
		zap.L().Error(err.Error())
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
	// TODO: Add logging to the application
	// Volume mapping in docker container ./app.env:/app/env/app.env
	config, err := initializers.LoadConfig(os.Getenv("API_ENV_CONFIG_PATH"))
	if err != nil {
		zap.L().Fatal("Could not load environment variables", zap.Error(err))
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
	zap.L().Fatal("running server on port: " + config.ServerPort,
		zap.Error(server.Run(":" + config.ServerPort)))
}

//	@Summary	Check the basic health of api
//	@Tags		Health
//	@Produce	json
//	@Success	200	{object}	map[string]string
//	@Router		/health [get]
func healthHandler(ctx *gin.Context) {
	message := "I am working! :)"
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
}

package controllers

import (
	"net/http"

	"github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type GistController struct {
	DB *gorm.DB
}

func NewGistController(DB *gorm.DB) GistController {
	return GistController{
		DB: DB,
	}
}


//	@Summary	Get the gist by gist id, does not load gist comments
//	@Tags		Gist Operations
//	@Produce	json
//	@Param		gistId	path		string	true	"The ID of the gist"
//	@Success	200		{object}	map[string]any
//	@Failure	404		{object}	map[string]string
//	@Router		/gists/{gistId} [get]
func (gc *GistController) GetGistById(ctx *gin.Context) {
	gistId := ctx.Params.ByName("gistId")

	var gist models.Gist
	result := gc.DB.
		Preload("GistContent").
		First(&gist, "id = ?", gistId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "gist does not exist"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"gist": gist}})
}

//	@Summary	Get the comments of a gist
//	@Tags		Gist Operations
//	@Produce	json
//	@Param		gistId	path		string	true	"The ID of the gist"
//	@Success	200		{object}	map[string]any
//	@Failure	404		{object}	map[string]string
//	@Router		/gists/{gistId}/comments [get]
func (gc *GistController) GetGistComments(ctx *gin.Context) {
	gistId := ctx.Params.ByName("gistId")

	var comments []models.Comment
	result := gc.DB.Find(&comments, "gist_id = ?", gistId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "gist does not exist"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"comments": comments}})
}

//	@Summary	Get the stargazers of a gist
//	@Tags		Gist Operations
//	@Produce	json
//	@Param		gistId	path		string	true	"The ID of the gist"
//	@Success	200		{object}	map[string]any
//	@Failure	400		{object}	map[string]any
//	@Failure	404		{object}	map[string]any
//	@Router		/gists/{gistId}/stargazers [get]
func (gc *GistController) GetGistStargazers(ctx *gin.Context) {
	gistId := ctx.Params.ByName("gistId")

	parsedGistId, err := uuid.Parse(gistId)
	if err != nil {
		zap.L().Error(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "invalid gist id"})
		return
	}

	var stars []models.Star
	result := gc.DB.Find(&stars, "gist_id = ?", parsedGistId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "gist does not exist"})
		return
	}

	var stargazers []string
	for _, star := range stars {
		stargazers = append(stargazers, star.Username)
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "stargazers": stargazers})
}




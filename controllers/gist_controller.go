package controllers

import (
	"net/http"

	"github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/models"
	"github.com/Vyom-Yadav/GitHub-Gist-Clone-Backend/utils"
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
//	@Success	200		{object}	models.GistWithoutCommentsWrapper
//	@Failure	404		{object}	models.ErrorResponseWrapper
//	@Failure	400		{object}	models.ErrorResponseWrapper
//	@Router		/gists/{gistId} [get]
func (gc *GistController) GetGistById(ctx *gin.Context) {
	gistId := ctx.Params.ByName("gistId")

	gistIdParsed, err := uuid.Parse(gistId)
	if err != nil {
		utils.NewErrorResponse(ctx, http.StatusBadRequest, "invalid gist id")
		zap.L().Error(err.Error())
		return
	}

	var gist models.Gist
	result := gc.DB.
		Preload("GistContent").
		First(&gist, "id = ?", gistIdParsed)
	if result.Error != nil {
		utils.NewErrorResponse(ctx, http.StatusNotFound, "gist does not exist")
		return
	}

	ctx.JSON(http.StatusOK, models.GistWithoutCommentsWrapper{
		Gist: models.GistWithoutComments{
			Username:    gist.Username,
			StarCount:   gist.StarCount,
			ID:          gist.ID,
			Private:     gist.Private,
			GistContent: gist.GistContent,
			Name:        gist.Name,
			Title:       gist.Title,
			CreatedAt:   gist.CreatedAt,
			UpdatedAt:   gist.UpdatedAt,
		},
	})
}

//	@Summary	Get the comments of a gist
//	@Tags		Gist Operations
//	@Produce	json
//	@Param		gistId	path		string	true	"The ID of the gist"
//	@Success	200		{object}	models.CommentArrayWrapper
//	@Failure	404		{object}	models.ErrorResponseWrapper
//	@Failure	400		{object}	models.ErrorResponseWrapper
//	@Router		/gists/{gistId}/comments [get]
func (gc *GistController) GetGistComments(ctx *gin.Context) {
	gistId := ctx.Params.ByName("gistId")

	gistIdParsed, err := uuid.Parse(gistId)
	if err != nil {
		utils.NewErrorResponse(ctx, http.StatusBadRequest, "invalid gist id")
		zap.L().Error(err.Error())
		return
	}

	var comments []models.Comment
	result := gc.DB.Find(&comments, "gist_id = ?", gistIdParsed)
	if result.Error != nil {
		utils.NewErrorResponse(ctx, http.StatusNotFound, "gist does not exist")
		return
	}

	ctx.JSON(http.StatusOK, models.CommentArrayWrapper{Comments: comments})
}

//	@Summary	Get the stargazers of a gist
//	@Tags		Gist Operations
//	@Produce	json
//	@Param		gistId	path		string	true	"The ID of the gist"
//	@Success	200		{object}	models.StringArrayWrapper
//	@Failure	400		{object}	models.ErrorResponseWrapper
//	@Failure	404		{object}	models.ErrorResponseWrapper
//	@Router		/gists/{gistId}/stargazers [get]
func (gc *GistController) GetGistStargazers(ctx *gin.Context) {
	gistId := ctx.Params.ByName("gistId")

	parsedGistId, err := uuid.Parse(gistId)
	if err != nil {
		zap.L().Error(err.Error())
		utils.NewErrorResponse(ctx, http.StatusBadRequest, "invalid gist id")
		return
	}

	var stars []models.Star
	result := gc.DB.Find(&stars, "gist_id = ?", parsedGistId)
	if result.Error != nil {
		utils.NewErrorResponse(ctx, http.StatusNotFound, "gist does not exist")
		return
	}

	var stargazers []string
	for _, star := range stars {
		stargazers = append(stargazers, star.Username)
	}

	ctx.JSON(http.StatusOK, models.StringArrayWrapper{StringArray: stargazers})
}

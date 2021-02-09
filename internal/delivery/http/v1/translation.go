package v1

import (
	"net/http"

	"github.com/evrone/go-service-template/internal/domain"

	_ "github.com/evrone/go-service-template/docs"
	"github.com/gin-gonic/gin"
)

func (r *router) translationRouts(api *gin.RouterGroup) {
	translation := api.Group("/translation")
	{
		translation.GET("/history", r.history)
		translation.POST("/do-translate", r.doTranslate)
	}
}

// TODO

// @Summary     Show history
// @Description Show all translation history
// @ID          history
// @Tags  	    translation
// @Accept      json
// @Produce     json
// @Success     200 {array} domain.Translation
// @Failure     400 {object} errorResponse
// @Router      /translation/history [get]
func (r *router) history(c *gin.Context) {
	entities, err := r.translationService.History()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"history": entities})
}

type translationInput struct {
	Source      string `json:"source"       binding:"required"  example:"auto"`
	Destination string `json:"destination"  binding:"required"  example:"en"`
	Original    string `json:"original"     binding:"required"  example:"текст для перевода"`
}

// @Summary     Translate
// @Description Translate a text
// @ID          do-translate
// @Tags  	    translation
// @Accept      json
// @Produce     json
// @Param       input body translationInput true "Set up translation"
// @Success     200 {object} domain.Translation
// @Failure     400 {object} errorResponse
// @Router      /translation/do-translate [post]
func (r *router) doTranslate(c *gin.Context) {
	var input translationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	translationResponse, err := r.translationService.DoTranslate(domain.Translation{
		Source:      input.Source,
		Destination: input.Destination,
		Original:    input.Original,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, translationResponse)
}

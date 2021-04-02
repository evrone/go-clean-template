package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/evrone/go-service-template/internal/domain"
)

func (r *router) translationRouts(api *gin.RouterGroup) {
	translation := api.Group("/translation")
	{
		translation.GET("/history", r.history)
		translation.POST("/do-translate", r.doTranslate)
	}
}

type historyResponse struct {
	History []domain.Translation `json:"history"`
}

// @Summary     Show history
// @Description Show all translation history
// @ID          history
// @Tags  	    translation
// @Accept      json
// @Produce     json
// @Success     200 {object} historyResponse
// @Failure     400 {object} response
// @Router      /translation/history [get].
func (r *router) history(c *gin.Context) {
	translations, err := r.translationService.History()
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err, "database problems")

		return
	}

	c.JSON(http.StatusOK, historyResponse{translations})
}

type doTranslateRequest struct {
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
// @Param       request body doTranslateRequest true "Set up translation"
// @Success     200 {object} domain.Translation
// @Failure     400 {object} response
// @Router      /translation/do-translate [post].
func (r *router) doTranslate(c *gin.Context) {
	var request doTranslateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse(c, http.StatusBadRequest, err, "invalid request body")

		return
	}

	translation, err := r.translationService.Translate(domain.Translation{
		Source:      request.Source,
		Destination: request.Destination,
		Original:    request.Original,
	})
	if err != nil {
		errorResponse(c, http.StatusBadRequest, err, "translation service problems")

		return
	}

	c.JSON(http.StatusOK, translation)
}

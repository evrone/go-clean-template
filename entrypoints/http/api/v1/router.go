package v1

import (
	"net/http"

	"github.com/evrone/go-service-template/business-logic/domain"

	"github.com/gin-gonic/gin"
)

type apiRouter struct {
	*gin.Engine
	useCase domain.EntityUseCase
}

func NewApiRouter(uc domain.EntityUseCase) http.Handler {
	router := gin.Default()
	api := &apiRouter{router, uc}

	v1 := router.Group("/api/v1")
	{
		v1.GET("/history", api.history)
		v1.POST("/do-translate", api.doTranslate)
	}

	return api
}

func (a *apiRouter) history(c *gin.Context) {
	entities, err := a.useCase.History()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entities)
}

func (a *apiRouter) doTranslate(c *gin.Context) {
	var entity domain.Entity
	if err := c.ShouldBindJSON(&entity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entity, err := a.useCase.DoTranslate(entity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entity)
}

package v1

import (
	"alukart32.com/bank/internal/usecase"
	"alukart32.com/bank/pkg/zerologx"
	"github.com/gin-gonic/gin"
)

type entryRoutes struct {
	service usecase.EntryService
	logger  zerologx.Logger
}

func newEntriesRoutes(handler *gin.RouterGroup, s usecase.EntryService, l zerologx.Logger) {
	r := entryRoutes{
		service: s,
		logger:  l,
	}

	h := handler.Group("/entries")
	{
		h.GET("/:id", r.getById)
		h.GET("/", r.list)
		h.POST("/", r.create)
		h.PUT("/:id/amount", r.updateAmount)
		h.DELETE("/:id", r.delete)
	}
}

func (r *entryRoutes) getById(c *gin.Context) {

}

func (r *entryRoutes) create(c *gin.Context) {

}

func (r *entryRoutes) updateAmount(c *gin.Context) {

}

func (r *entryRoutes) list(c *gin.Context) {

}

func (r *entryRoutes) delete(c *gin.Context) {

}

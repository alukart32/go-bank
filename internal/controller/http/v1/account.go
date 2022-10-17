package v1

import (
	"alukart32.com/bank/internal/usecase"
	"github.com/gin-gonic/gin"
)

type accountsRoutes struct {
	s usecase.AccountService
}

func newAccountsRoutes(handler *gin.RouterGroup, s usecase.AccountService) {
	r := &accountsRoutes{s}

	h := handler.Group("/accounts")
	{
		h.GET("/:id", r.get)
		h.DELETE("/:id", r.delete)
	}
}

func (r *accountsRoutes) get(c *gin.Context) {

}

func (r *accountsRoutes) delete(c *gin.Context) {

}

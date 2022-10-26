package v1

import (
	"net/http"

	"alukart32.com/bank/entity"
	"alukart32.com/bank/internal/usecase"
	"alukart32.com/bank/pkg/zerologx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type accountRoutes struct {
	service usecase.AccountService
	logger  zerologx.Logger
}

func newAccountsRoutes(handler *gin.RouterGroup, s usecase.AccountService, l zerologx.Logger) {
	r := &accountRoutes{
		service: s,
		logger:  l,
	}

	h := handler.Group("/accounts")
	{
		h.GET("/:id", r.getById)
		h.GET("/:id/entries", r.listEntries)
		h.GET("/:id/transfers", r.listTransfers)
		h.POST("/", r.create)
		h.POST("/add", r.addBalance)
		h.PATCH("/:id/owner", r.updateOwner)
		h.DELETE("/:id", r.delete)
	}
}

func (r *accountRoutes) getById(c *gin.Context) {
	id := c.Param("id")
	account, err := r.service.Get(
		c.Request.Context(),
		uuid.MustParse(id),
	)
	if err != nil {
		r.logger.Error(err, "http - v1 - account - getByID")
		errorResponse(c, http.StatusInternalServerError, "account service problems")

		return
	}

	c.JSON(http.StatusOK, account)
}

type createAccountReq struct {
}

func (r *accountRoutes) create(c *gin.Context) {
	var accountToCreate entity.Account

	// TODO: Add business validation
	if err := c.BindJSON(&accountToCreate); err != nil {
		r.logger.Error(err, "http - v1 - account")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	id, err := r.service.Create(c.Request.Context(), accountToCreate)
	if err != nil {
		r.logger.Error(err, "http - v1 - account - create")
		errorResponse(c, http.StatusInternalServerError, "account service problems")

		return
	}

	c.JSON(http.StatusOK, id)
}

func (r *accountRoutes) updateOwner(c *gin.Context) {

}

func (r *accountRoutes) addBalance(c *gin.Context) {

}

func (r *accountRoutes) listEntries(c *gin.Context) {

}

func (r *accountRoutes) listTransfers(c *gin.Context) {

}

func (r *accountRoutes) delete(c *gin.Context) {

}

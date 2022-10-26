package v1

import (
	"net/http"

	"alukart32.com/bank/entity"
	"alukart32.com/bank/internal/usecase"
	"alukart32.com/bank/pkg/zerologx"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type transferRoutes struct {
	service usecase.TransferService
	logger  zerologx.Logger
}

func newTransfersRoutes(handler *gin.RouterGroup, s usecase.TransferService, l zerologx.Logger) {
	r := transferRoutes{
		service: s,
		logger:  l,
	}

	h := handler.Group("/transfers")
	{
		h.GET("/:id", r.getById)
		h.GET("/", r.list)
		h.POST("/transfer", r.transfer)
		h.DELETE("/rollback/:id", r.rollback)
	}
}

func (r *transferRoutes) getById(c *gin.Context) {

}

func (r *transferRoutes) list(c *gin.Context) {

}

type doTransferRequest struct {
	FromAccountID uuid.UUID `json:"fromAccountID" binding:"required"`
	ToAccountID   uuid.UUID `json:"toAccountID"  binding:"required"`
	Amount        int64     `json:"amount"     binding:"required"`
	Currency      string    `json:"currency"  binding:"required" validate:"required,oneof=rub,usd"`
}

func (r *transferRoutes) transfer(c *gin.Context) {
	var request doTransferRequest
	// TODO: Add business validation
	if err := c.BindJSON(&request); err != nil {
		r.logger.Error(err, "http - v1 - transfer")
		errorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	translation, err := r.service.Transfer(
		c.Request.Context(),
		entity.Transfer{
			FromAccountID: request.FromAccountID,
			ToAccountID:   request.ToAccountID,
			Amount:        request.Amount,
		},
	)
	if err != nil {
		r.logger.Error(err, "http - v1 - doTranslate")
		errorResponse(c, http.StatusInternalServerError, "translation service problems")

		return
	}

	c.JSON(http.StatusOK, translation)
}

func (r *transferRoutes) rollback(c *gin.Context) {

}

package v1

import (
	"net/http"

	"alukart32.com/bank/internal/usecase"
	"alukart32.com/bank/pkg/middleware"
	"alukart32.com/bank/pkg/zerologx"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, l zerologx.Logger, as usecase.AccountService,
	es usecase.EntryService, ts usecase.TransferService) http.Handler {
	// Routes
	h := handler.Group("/v1")
	h.Use(middleware.AuthJWT())
	{
		newAccountsRoutes(h, as, l)
		newEntriesRoutes(h, es, l)
		newTransfersRoutes(h, ts, l)
	}

	return handler
}

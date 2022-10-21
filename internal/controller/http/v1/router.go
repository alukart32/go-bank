package v1

import (
	"net/http"

	"alukart32.com/bank/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(engine *gin.Engine, as usecase.AccountService,
	es usecase.EntryService, ts usecase.TransferService) http.Handler {
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	engine.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	engine.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	h := engine.Group("/v1")
	{
		newAccountsRoutes(h, as)
		newEntriesRoutes(h, es)
		newTransfersRoutes(h, ts)
	}

	return engine
}

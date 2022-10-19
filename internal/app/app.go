package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"alukart32.com/bank/config"
	v1 "alukart32.com/bank/internal/controller/http/v1"
	"alukart32.com/bank/internal/usecase"
	"alukart32.com/bank/internal/usecase/repo"
	"alukart32.com/bank/pkg/httpserver"
	"alukart32.com/bank/pkg/postgres"
	"github.com/gin-gonic/gin"
)

func Run(cfg config.Config) {
	fail := func(err error) {
		log.Fatal(err)
	}

	// Prepare tools
	db, err := postgres.New(&cfg.DB)
	if err != nil {
		fail(fmt.Errorf("app - init db instance error: " + err.Error()))
	}

	accountService := usecase.NewAccountService(repo.NewAccountSQLRepo(db))
	entryService := usecase.NewEntryService(repo.NewEntrySQLRepo(db))
	transferService := usecase.NewTransferService(repo.NewTransferSQLRepo(db))

	gin := gin.New()
	v1.NewRouter(gin, accountService, entryService, transferService)
	httpServer := httpserver.New(gin, cfg.Http)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Println(s.String())
	case err = <-httpServer.Notify():
		log.Print(fmt.Errorf("app - Run - httpServer.Notify: %v", err))
	}

	// Shutdown
	if err = httpServer.Shutdown(); err != nil {
		log.Print(fmt.Errorf("app - Run - httpServer.Shutdown: %v", err))
	}
}

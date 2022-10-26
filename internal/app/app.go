package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"alukart32.com/bank/config"
	v1 "alukart32.com/bank/internal/controller/http/v1"
	"alukart32.com/bank/internal/usecase"
	"alukart32.com/bank/internal/usecase/repo"
	"alukart32.com/bank/pkg/ginx"
	"alukart32.com/bank/pkg/httpserver"
	"alukart32.com/bank/pkg/postgres"
	"alukart32.com/bank/pkg/zerologx"
)

func Run(cfg config.Config) {
	logger := zerologx.New(cfg.Logger.Level, nil)

	fail := func(err error) {
		logger.Fatal(err)
	}

	// Prepare tools
	db, err := postgres.New(cfg.DB)
	if err != nil {
		fail(fmt.Errorf("app - init db instance error: " + err.Error()))
	}

	accountService := usecase.NewAccountService(repo.NewAccountSQLRepo(db), &logger)
	entryService := usecase.NewEntryService(repo.NewEntrySQLRepo(db), &logger)
	transferService := usecase.NewTransferService(repo.NewTransferSQLRepo(db), &logger)

	handler := v1.NewRouter(ginx.NewGinEngine(), &logger, accountService, entryService, transferService)
	httpServer := httpserver.New(handler, cfg.Http)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info(s.String())
	case err = <-httpServer.Notify():
		logger.Error(fmt.Errorf("app - Run - httpServer.Notify: %v", err))
	}

	// Shutdown
	if err = httpServer.Shutdown(); err != nil {
		logger.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %v", err))
	}

	if err = postgres.Close(); err != nil {
		logger.Error(fmt.Errorf("app - Run - postgres.Close: %v", err))
	}
}

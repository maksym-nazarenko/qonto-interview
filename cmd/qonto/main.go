package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/maxim-nazarenko/qonto-interview/internal/qonto"
	"github.com/maxim-nazarenko/qonto-interview/internal/qonto/api"
	"github.com/maxim-nazarenko/qonto-interview/internal/qonto/app"
	"github.com/maxim-nazarenko/qonto-interview/internal/qonto/core"
	"github.com/maxim-nazarenko/qonto-interview/internal/qonto/storage"
	"github.com/maxim-nazarenko/qonto-interview/internal/qonto/utils"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run(args []string) error {
	appCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	appLogger := qonto.NewInstanceLogger(os.Stdout, "Qonto")
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGQUIT, syscall.SIGHUP, syscall.SIGTERM)
	go func() {
		<-signalChan
		appLogger.Info("received interruption request, closing the app")
		cancel()
	}()

	config, err := app.ConfigurationFromEnv(os.Getenv)
	if err != nil {
		return err
	}

	mysqlConfig := storage.NewMysqlConfig()
	mysqlConfig.User = config.DB.User
	mysqlConfig.Passwd = config.DB.Password
	mysqlConfig.DBName = config.DB.Name
	mysqlConfig.Net = "tcp"
	mysqlConfig.Addr = config.DB.Address

	mysqlStorage, err := storage.NewMysqlStorage(mysqlConfig)
	if err != nil {
		return err
	}
	defer func() {
		if err := mysqlStorage.Close(); err != nil {
			appLogger.Error("could not close database connection: %v", err)
		}
	}()
	if err := dbConnect(appCtx, mysqlStorage, appLogger); err != nil {
		return err
	}
	projectRoot := utils.ProjectRootDir()
	if err := storage.Migrate("file://"+projectRoot+"/migrations/", mysqlStorage.DB()); err != nil {
		return fmt.Errorf("migrations failed: %v", err)
	}
	appLogger.Info("migration completed")

	qontoAPI := api.NewAPI(core.NewQontoTransferManager(mysqlStorage))
	router := chi.NewRouter()
	router.Post("/v1/transfers", qontoAPI.HandleTransfers)

	server := http.Server{
		Addr:         "127.0.0.1:8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      router,
	}

	wg := sync.WaitGroup{}

	// starting HTTP server in background
	wg.Add(1)
	go func(srv *http.Server, logger qonto.Logger) {
		defer wg.Done()
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(err.Error())
		}

	}(&server, appLogger.SubLogger("http server"))

	// app shutdown watcher that will stop HTTP server
	wg.Add(1)
	go func(ctx context.Context, srv *http.Server, logger qonto.Logger) {
		defer wg.Done()
		<-ctx.Done()

		timeout := 10 * time.Second
		shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		logger.Info("trying gracefully shutdown in %v", timeout)
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Error(err.Error())
		}
		logger.Info("done")
	}(appCtx, &server, appLogger.SubLogger("http server"))

	wg.Wait()

	appLogger.Info("all tasks stopped, exiting application")

	return nil
}

func dbConnect(ctx context.Context, mysqlStorage storage.Storage, appLogger qonto.Logger) error {
	dbPingCtx, dbPingCancel := context.WithTimeout(ctx, 30*time.Second)
	defer dbPingCancel()

	dbUpWaitFunc := func(db *sql.DB) (bool, error) {
		for {
			select {
			case <-dbPingCtx.Done():
				return false, dbPingCtx.Err()
			case <-time.After(1 * time.Second):
				if err := db.PingContext(dbPingCtx); err != nil {
					appLogger.Info("db ping failed: %v", err)
					return true, err
				}
				return false, nil
			}
		}
	}
	if err := mysqlStorage.Wait(dbUpWaitFunc); err != nil {
		return err
	}

	return nil
}

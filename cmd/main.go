package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/inspectorvitya/shortenerURL/internal/app"
	httpserver "github.com/inspectorvitya/shortenerURL/internal/server"
	"github.com/inspectorvitya/shortenerURL/internal/storage"
	"github.com/inspectorvitya/shortenerURL/internal/storage/memory"
	"github.com/inspectorvitya/shortenerURL/internal/storage/sql"
	"go.uber.org/zap"
)

type Config struct {
	Port  string `env:"PORT" env-default:"8080"`
	DBURL string `env:"DBURL"`
}

func main() {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalln(err)
	}

	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{"stdout"}
	zapLogger, err := config.Build()
	if err != nil {
		log.Fatal("failed to init logger: ", err)
	}
	zap.ReplaceGlobals(zapLogger)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	var db storage.Storage
	if cfg.DBURL == "" {
		db = memory.New()
	} else {
		fmt.Println(cfg.DBURL)
		sqlDB, err := sql.New(cfg.DBURL)
		if err != nil {
			zap.L().Fatal("failed connect db: ", zap.Error(err))
		}
		defer sqlDB.Close()
		db = sqlDB
	}

	shorter := app.New(db)

	server := httpserver.New(cfg.Port, shorter)

	go func() {
		zap.L().Info("http server start...")
		if err := server.Start(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				zap.L().Info("http server stopped....")
			} else {
				zap.L().Fatal("failed to start http server: ", zap.Error(err))
			}
		}
	}()

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = server.Stop(ctx)
	if err != nil {
		zap.L().Error("failed stop http server: ", zap.Error(err))
	}
}

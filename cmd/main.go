package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"os"
	"time"

	config "go-gin-hex-arch/internal/adapter/config"
	"go-gin-hex-arch/internal/adapter/http"
	myLogger "go-gin-hex-arch/internal/adapter/logger"
	"go-gin-hex-arch/internal/adapter/middleware"
	"go-gin-hex-arch/internal/core/service"
	"go-gin-hex-arch/internal/storage/mongodb"
	"go-gin-hex-arch/internal/storage/mysql"
	"go-gin-hex-arch/internal/storage/mysql/repository"
	"go-gin-hex-arch/internal/util"
)

func main() {
	if err := util.InitTimeZone(); err != nil {
		panic(err.Error())
	}

	time.Local = util.Loc

	cfg, err := config.New()
	if err != nil {
		slog.Error("Error loading environment", "error", err)
		os.Exit(1)
	}

	myLogger.Set()

	slog.Info("Starting the app", "app", cfg.App.Name, "env", cfg.App.Env)

	router := gin.Default()

	mySqlDB, err := mysql.ConnectMySQL(cfg.DB)
	if err != nil {
		panic(err)
	}

	mongoDB, err := mongodb.NewDB(cfg.MONGO)
	if err != nil {
		panic(err)
	}

	productRepo := repository.NewProductRepositoryDB(mySqlDB)
	productService := service.NewProductService(productRepo)

	router.Use(middleware.MonitoringFuncPerformance(mongoDB, cfg))

	http.SetupRoutes(router, productService, mongoDB, cfg.HTTP)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, "Hello World")
	})

	listenPort := fmt.Sprintf(":%s", cfg.HTTP.Port)
	if err := router.Run(listenPort); err != nil {
		slog.Error(fmt.Sprintf("Error listening on port %s", cfg.HTTP.Port), "error", err)
		os.Exit(1)
	}
}

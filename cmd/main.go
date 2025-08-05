package main

import (
	"log"
	"os"
	"wangfeng/transaction-system/internal/config"
	"wangfeng/transaction-system/internal/controller"
	"wangfeng/transaction-system/internal/db"
	"wangfeng/transaction-system/internal/repository"
	"wangfeng/transaction-system/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/yaml.v3"
)

var configFile = "conf/config.yaml"

func main() {
	appConfig, err := readConfig()
	if err != nil {
		log.Fatalf("Could not read configuration: %v", err)
	}
	if err := db.InitDB(appConfig.DatabaseURL); err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	accountRepository, err := repository.NewAccountRepository()
	if err != nil {
		log.Fatalf("Could not create account repository: %v", err)
	}

	accountService, err := service.NewAccountService(accountRepository)
	if err != nil {
		log.Fatalf("Could not create account service: %v", err)
	}
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	_, err = controller.NewAccountController(e, accountService)
	if err != nil {
		log.Fatalf("Could not create account controller: %v", err)
	}
	log.Println("Starting server on :8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}

func readConfig() (*config.AppConfig, error) {
	content, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	// unmarshal
	var appConfig config.AppConfig
	if err = yaml.Unmarshal(content, &appConfig); err != nil {
		return nil, err
	}
	return &appConfig, nil
}

package main

import (
	"log"

	"billing-app/config"
	"billing-app/models"
	"billing-app/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	db, err := models.ConnectDatabase(cfg.MySQLDSN)
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	routes.Register(router, db)

	log.Printf("billing API listening on :%s", cfg.Port)
	return router.Run(":" + cfg.Port)
}

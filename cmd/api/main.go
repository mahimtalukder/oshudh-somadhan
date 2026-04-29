package main

import (
	"log"

	"github.com/mahimtalukder/oshudh-somadhan/internal/config"
	"github.com/mahimtalukder/oshudh-somadhan/internal/database"
	"github.com/mahimtalukder/oshudh-somadhan/internal/server"
)

func main() {
	//load env
	cfg, err := config.Load()
	if err != nil{
		log.Fatalf("error on loading env. err: %v", err)
	}

	//setup db
	db, err := database.Connect(cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("failed to connect to db. err: %v", err)
	}
	// close db when work done
	defer db.Close()

	//setup router
	router := server.NewRouter(db)

	log.Printf("%s server running on port %s", cfg.AppName, cfg.AppPort)
	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("failed to start server. err: %v", err)
	}

}
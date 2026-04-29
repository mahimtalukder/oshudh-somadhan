package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mahimtalukder/oshudh-somadhan/internal/handler"
)

func NewRouter(db *pgxpool.Pool) *gin.Engine{
	router := gin.Default()

	//health handlers
	healthHandlers := handler.NewHealthHandler(db)

	router.GET("/health", healthHandlers.Health)
	router.GET("/debug/db-ping", healthHandlers.CheckDbConnection)

	return router
}
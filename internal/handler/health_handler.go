package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	DB *pgxpool.Pool
}

func NewHealthHandler(db *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{
		DB: db,
	}
}

func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "Ok",
		"message": "Oshudh Somadhan API is running",
	})
}

func (h *HealthHandler) CheckDbConnection(c *gin.Context){
	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()

	if err := h.DB.Ping(ctx); err != nil {
		message := fmt.Sprintf("Database connection failed. err: %v", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "ERROR",
			"message": message,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "Ok",
		"message": "Database connected successfully",
	})
}

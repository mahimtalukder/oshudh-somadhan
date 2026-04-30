package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mahimtalukder/oshudh-somadhan/internal/response"
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
	response.Success(c, http.StatusOK, "Oshudh Somadhan API is running", nil)
}

func (h *HealthHandler) CheckDbConnection(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := h.DB.Ping(ctx); err != nil {
		message := fmt.Sprintf("Database connection failed. err: %v", err)
		response.Error(c, http.StatusServiceUnavailable, message, err)
	}

	response.Success(c, http.StatusOK, "Database connected successfully", nil)
}

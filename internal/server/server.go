package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mahimtalukder/oshudh-somadhan/internal/handler"
	"github.com/mahimtalukder/oshudh-somadhan/internal/repository"
)

func NewRouter(db *pgxpool.Pool) *gin.Engine {
	router := gin.Default()

	//health handlers
	healthHandlers := handler.NewHealthHandler(db)

	//medicine handler set
	medicineRepository := repository.NewMedicineRepository(db)
	medicineHandler := handler.NewMedicineHandler(medicineRepository)

	//router group
	v1 := router.Group("/api/v1")

	//health check routes
	router.GET("/health", healthHandlers.Health)
	router.GET("/debug/db-ping", healthHandlers.CheckDbConnection)

	//medicine routes

	medicineRoute := v1.Group("/medicines")
	medicineRoute.GET("/search", medicineHandler.SearchMedicines)
	v1.GET("/medicines", medicineHandler.ListMedicines)
	medicineRoute.GET("/:id", medicineHandler.GetMedicineByID)
	medicineRoute.GET("/:id/alternatives", medicineHandler.GetAlternatives)
	v1.GET("/generic/:generic_id/medicines", medicineHandler.GetMedicinesByGenericID)

	return router
}

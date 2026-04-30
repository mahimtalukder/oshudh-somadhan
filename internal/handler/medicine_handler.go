package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/mahimtalukder/oshudh-somadhan/internal/helper"
	"github.com/mahimtalukder/oshudh-somadhan/internal/repository"
	"github.com/mahimtalukder/oshudh-somadhan/internal/response"
)

type MedicineHandler struct {
	repopository *repository.MedicineRepository
}

func NewMedicineHandler(repo *repository.MedicineRepository) *MedicineHandler {
	return &MedicineHandler{
		repopository: repo,
	}
}

func (h *MedicineHandler) ListMedicines(c *gin.Context) {
	page, limit, ok := helper.ParsePageLimit(c)
	if !ok {
		return
	}

	genericID, ok := helper.ParseOptionalIntQuery(c, "generic_id")
	if !ok {
		return
	}

	companyID, ok := helper.ParseOptionalIntQuery(c, "company_id")
	if !ok {
		return
	}

	dosageFormID, ok := helper.ParseOptionalIntQuery(c, "dosage_form_id")
	if !ok {
		return
	}

	medicines, err := h.repopository.ListMedicines(c, page, limit, genericID, companyID, dosageFormID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch medicines", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Medicines fetched successfully", medicines)
}

func (h *MedicineHandler) GetMedicineByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id < 1 {
		response.Error(c, http.StatusBadRequest, "id must be a positive integer", nil)
		return
	}

	medicine, err := h.repopository.GetMedicineByID(c, id)
	if errors.Is(err, pgx.ErrNoRows) {
		response.Error(c, http.StatusNotFound, fmt.Sprintf("medicine not found with this id: %d", id), err.Error())
		return
	}
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch medicine", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "medicine successfully found", medicine)
}

func (h *MedicineHandler) SearchMedicines(c *gin.Context) {
	q := strings.TrimSpace(c.Query("q"))
	if q == "" {
		response.Error(c, http.StatusBadRequest, "Search query is required", "q is required")
		return
	}

	page, limit, ok := helper.ParsePageLimit(c)
	if !ok {
		return
	}

	medicines, err := h.repopository.SearchMedicines(c, q, limit, page)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch medicines", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Medicines successfully fetched", medicines)
}

func (h *MedicineHandler) GetMedicinesByGenericID(c *gin.Context) {
	page, limit, ok := helper.ParsePageLimit(c)
	if !ok {
		return
	}

	genericIdParam := c.Param("generic_id")
	genericId, err := strconv.Atoi(genericIdParam)
	if err != nil || genericId < 1 {
		response.Error(c, http.StatusBadRequest, "Invalid generic ID", "id must be a positive number")
		return
	}

	medicines, err := h.repopository.GetMedicinesByGenericID(c, genericId, limit, page)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch medicine", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Medicine successfully fetched", medicines)
}

func (h *MedicineHandler) GetAlternatives(c *gin.Context) {
	medicineId, err := strconv.Atoi(c.Param("id"))
	if err != nil || medicineId < 1 {
		response.Error(c, http.StatusBadRequest, "Invalid medicine id", "Medicine id must be positive number")
		return
	}

	page, limit, ok := helper.ParsePageLimit(c)
	if !ok {
		return
	}

	alternatives, err := h.repopository.GetAlternatives(c, medicineId, limit, page)
	if errors.Is(err, pgx.ErrNoRows) {
		response.Error(c, http.StatusNotFound, "No medicine found", "Invalid medicine Id")
		return
	}
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch medicines", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Alternatives successfully fetched", alternatives)
}

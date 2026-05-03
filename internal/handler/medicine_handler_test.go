package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mahimtalukder/oshudh-somadhan/internal/repository"
	"github.com/mahimtalukder/oshudh-somadhan/internal/testutil"
)

func setupMedicineTestRouter(t *testing.T) (*gin.Engine, testutil.SeededCatalog) {
	t.Helper()

	gin.SetMode(gin.TestMode)

	db := testutil.TestDB(t)
	seed := testutil.SeedMedicineCatalog(t, db)

	medicineRepo := repository.NewMedicineRepository(db)
	medicineHandler := NewMedicineHandler(medicineRepo)

	router := gin.New()

	v1 := router.Group("/api/v1")
	{
		medicineRoute := v1.Group("/medicines")
		{
			medicineRoute.GET("/", medicineHandler.ListMedicines)
			medicineRoute.GET("/search", medicineHandler.SearchMedicines)
			medicineRoute.GET("/:id", medicineHandler.GetMedicineByID)
			medicineRoute.GET("/:id/alternatives", medicineHandler.GetAlternatives)
		}

		genericRoute := v1.Group("/generics")
		{
			genericRoute.GET("/:id/medicines", medicineHandler.GetMedicinesByGenericID)
		}
	}

	return router, seed
}

func performRequest(router http.Handler, method string, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	return recorder
}

func TestMedicineHandler_ListMedicines_Success(t *testing.T) {
	router, _ := setupMedicineTestRouter(t)

	recorder := performRequest(router, http.MethodGet, "/api/v1/medicines/?page=1&limit=10")

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d. body: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
}

func TestMedicineHandler_GetMedicineByID_Success(t *testing.T) {
	router, seed := setupMedicineTestRouter(t)

	path := fmt.Sprintf("/api/v1/medicines/%d", seed.MedicineID)
	recorder := performRequest(router, http.MethodGet, path)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d. body: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
}

func TestMedicineHandler_GetMedicineByID_NotFound(t *testing.T) {
	router, _ := setupMedicineTestRouter(t)

	recorder := performRequest(router, http.MethodGet, "/api/v1/medicines/999999999")

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d. body: %s", http.StatusNotFound, recorder.Code, recorder.Body.String())
	}
}

func TestMedicineHandler_GetMedicineByID_InvalidID(t *testing.T) {
	router, _ := setupMedicineTestRouter(t)

	recorder := performRequest(router, http.MethodGet, "/api/v1/medicines/abc")

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d. body: %s", http.StatusBadRequest, recorder.Code, recorder.Body.String())
	}
}

func TestMedicineHandler_SearchMedicines_Success(t *testing.T) {
	router, _ := setupMedicineTestRouter(t)

	recorder := performRequest(router, http.MethodGet, "/api/v1/medicines/search?q=Test%20Napa&page=1&limit=10")

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d. body: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
}

func TestMedicineHandler_SearchMedicines_MissingQuery(t *testing.T) {
	router, _ := setupMedicineTestRouter(t)

	recorder := performRequest(router, http.MethodGet, "/api/v1/medicines/search")

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d. body: %s", http.StatusBadRequest, recorder.Code, recorder.Body.String())
	}
}

func TestMedicineHandler_ListMedicines_InvalidPage(t *testing.T) {
	router, _ := setupMedicineTestRouter(t)

	recorder := performRequest(router, http.MethodGet, "/api/v1/medicines/?page=abc&limit=10")

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d. body: %s", http.StatusBadRequest, recorder.Code, recorder.Body.String())
	}
}

func TestMedicineHandler_ListMedicines_InvalidLimit(t *testing.T) {
	router, _ := setupMedicineTestRouter(t)

	recorder := performRequest(router, http.MethodGet, "/api/v1/medicines/?page=1&limit=abc")

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d. body: %s", http.StatusBadRequest, recorder.Code, recorder.Body.String())
	}
}

func TestMedicineHandler_PaginationResponse(t *testing.T) {
	router, _ := setupMedicineTestRouter(t)

	recorder := performRequest(router, http.MethodGet, "/api/v1/medicines/?page=1&limit=2")

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d. body: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}

	var body map[string]any

	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse json response: %v", err)
	}

	data, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data object, got %+v", body["data"])
	}

	limit, ok := data["limit"].(float64)
	if !ok {
		t.Fatalf("expected limit in response, got %+v", data["limit"])
	}

	if int(limit) != 2 {
		t.Fatalf("expected limit 2, got %v", limit)
	}
}

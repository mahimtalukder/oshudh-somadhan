package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/mahimtalukder/oshudh-somadhan/internal/testutil"
)

func TestMedicineRepository_ListMedicines(t *testing.T) {
	db := testutil.TestDB(t)
	seed := testutil.SeedMedicineCatalog(t, db)

	repo := NewMedicineRepository(db)

	result, err := repo.ListMedicines(
		context.Background(),
		1,
		10,
		&seed.GenericID,
		&seed.CompanyID,
		&seed.DosageFormID,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
	}

	if result.Page != 1 {
		t.Errorf("expected page 1, got %d", result.Page)
	}

	if result.Limit != 10 {
		t.Errorf("expected limit 10, got %d", result.Limit)
	}

	if result.Total < 3 {
		t.Errorf("expected at least 3 medicines, got %d", result.Total)
	}

	if len(result.Data) == 0 {
		t.Fatal("expected medicine data, got empty list")
	}
}

func TestMedicineRepository_GetMedicineByID_Success(t *testing.T) {
	db := testutil.TestDB(t)
	seed := testutil.SeedMedicineCatalog(t, db)

	repo := NewMedicineRepository(db)

	medicine, err := repo.GetMedicineByID(context.Background(), seed.MedicineID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if medicine == nil {
		t.Fatal("expected medicine, got nil")
	}

	if medicine.ID != seed.MedicineID {
		t.Errorf("expected medicine id %d, got %d", seed.MedicineID, medicine.ID)
	}

	if medicine.GenericID != seed.GenericID {
		t.Errorf("expected generic id %d, got %d", seed.GenericID, medicine.GenericID)
	}

	if medicine.CompanyID != seed.CompanyID {
		t.Errorf("expected company id %d, got %d", seed.CompanyID, medicine.CompanyID)
	}

	if medicine.DosageFormID != seed.DosageFormID {
		t.Errorf("expected dosage form id %d, got %d", seed.DosageFormID, medicine.DosageFormID)
	}
}

func TestMedicineRepository_GetMedicineByID_NotFound(t *testing.T) {
	db := testutil.TestDB(t)
	repo := NewMedicineRepository(db)

	medicine, err := repo.GetMedicineByID(context.Background(), 999999999)

	if !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("expected pgx.ErrNoRows, got %v", err)
	}

	if medicine != nil {
		t.Fatalf("expected medicine nil, got %+v", medicine)
	}
}

func TestMedicineRepository_SearchMedicines(t *testing.T) {
	db := testutil.TestDB(t)
	testutil.SeedMedicineCatalog(t, db)

	repo := NewMedicineRepository(db)

	result, err := repo.SearchMedicines(context.Background(), "Test Napa", 10, 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
	}

	if len(result.Data) == 0 {
		t.Fatal("expected search result, got empty list")
	}
}

func TestMedicineRepository_Pagination(t *testing.T) {
	db := testutil.TestDB(t)
	seed := testutil.SeedMedicineCatalog(t, db)

	repo := NewMedicineRepository(db)

	result, err := repo.ListMedicines(
		context.Background(),
		1,
		2,
		&seed.GenericID,
		&seed.CompanyID,
		&seed.DosageFormID,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Page != 1 {
		t.Errorf("expected page 1, got %d", result.Page)
	}

	if result.Limit != 2 {
		t.Errorf("expected limit 2, got %d", result.Limit)
	}

	if len(result.Data) > 2 {
		t.Errorf("expected max 2 medicines, got %d", len(result.Data))
	}

	if result.TotalPages < 2 {
		t.Errorf("expected at least 2 total pages, got %d", result.TotalPages)
	}
}

func TestMedicineRepository_GetAlternatives(t *testing.T) {
	db := testutil.TestDB(t)
	seed := testutil.SeedMedicineCatalog(t, db)

	repo := NewMedicineRepository(db)

	result, err := repo.GetAlternatives(context.Background(), seed.MedicineID, 10, 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
	}

	if len(result.Data) == 0 {
		t.Fatal("expected alternatives, got empty list")
	}

	for _, item := range result.Data {
		if item.ID == seed.MedicineID {
			t.Fatalf("alternative list should not contain current medicine id %d", seed.MedicineID)
		}
	}
}

package testutil

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SeededCatalog struct {
	GenericID     int
	CompanyID     int
	DosageFormID  int
	MedicineID    int
	AlternativeID int
	OtherMedID     int
}

//set db $env:TEST_DATABASE_URL="postgres://postgres:postgres@localhost:5433/oshudh_somadhan_test?sslmode=disable"
func TestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("failed to connect test database: %v", err)
	}

	if err := db.Ping(ctx); err != nil {
		t.Fatalf("failed to ping test database: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func SeedMedicineCatalog(t *testing.T, db *pgxpool.Pool) SeededCatalog {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	unique := time.Now().UnixNano()

	genericName := fmt.Sprintf("Test Generic %d", unique)
	companyName := fmt.Sprintf("Test Company %d", unique)
	dosageName := fmt.Sprintf("Test Tablet %d", unique)

	var seed SeededCatalog

	err := db.QueryRow(ctx, `
		INSERT INTO generics (name, siteid, slug, type, url)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, genericName, int(unique%100000), fmt.Sprintf("test-generic-%d", unique), "ALLOPATHIC", "https://example.com/generic").Scan(&seed.GenericID)

	if err != nil {
		t.Fatalf("failed to insert generic: %v", err)
	}

	err = db.QueryRow(ctx, `
		INSERT INTO companie (name, siteid, url)
		VALUES ($1, $2, $3)
		RETURNING id
	`, companyName, int(unique%100000), "https://example.com/company").Scan(&seed.CompanyID)

	if err != nil {
		t.Fatalf("failed to insert company: %v", err)
	}

	err = db.QueryRow(ctx, `
		INSERT INTO dosageform (name, siteid, imglink)
		VALUES ($1, $2, $3)
		RETURNING id
	`, dosageName, int(unique%100000), "https://example.com/tablet.png").Scan(&seed.DosageFormID)

	if err != nil {
		t.Fatalf("failed to insert dosage form: %v", err)
	}

	err = db.QueryRow(ctx, `
		INSERT INTO medicine 
			(name, packsizeandprice, siteid, strength, url, companie_id, dosageform_id, generics_id)
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`,
		fmt.Sprintf("Test Napa %d", unique),
		"10 tablets: 20 tk",
		int(unique%100000),
		"500 mg",
		"https://example.com/medicine-1",
		seed.CompanyID,
		seed.DosageFormID,
		seed.GenericID,
	).Scan(&seed.MedicineID)

	if err != nil {
		t.Fatalf("failed to insert medicine: %v", err)
	}

	err = db.QueryRow(ctx, `
		INSERT INTO medicine 
			(name, packsizeandprice, siteid, strength, url, companie_id, dosageform_id, generics_id)
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`,
		fmt.Sprintf("Test Alternative %d", unique),
		"10 tablets: 15 tk",
		int(unique%100000)+1,
		"500 mg",
		"https://example.com/medicine-2",
		seed.CompanyID,
		seed.DosageFormID,
		seed.GenericID,
	).Scan(&seed.AlternativeID)

	if err != nil {
		t.Fatalf("failed to insert alternative medicine: %v", err)
	}

	err = db.QueryRow(ctx, `
		INSERT INTO medicine 
			(name, packsizeandprice, siteid, strength, url, companie_id, dosageform_id, generics_id)
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`,
		fmt.Sprintf("Test Other Medicine %d", unique),
		"5 tablets: 10 tk",
		int(unique%100000)+2,
		"250 mg",
		"https://example.com/medicine-3",
		seed.CompanyID,
		seed.DosageFormID,
		seed.GenericID,
	).Scan(&seed.OtherMedID)

	if err != nil {
		t.Fatalf("failed to insert other medicine: %v", err)
	}

	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cleanupCancel()

		_, _ = db.Exec(cleanupCtx, `DELETE FROM medicine WHERE id IN ($1, $2, $3)`, seed.MedicineID, seed.AlternativeID, seed.OtherMedID)
		_, _ = db.Exec(cleanupCtx, `DELETE FROM dosageform WHERE id = $1`, seed.DosageFormID)
		_, _ = db.Exec(cleanupCtx, `DELETE FROM companie WHERE id = $1`, seed.CompanyID)
		_, _ = db.Exec(cleanupCtx, `DELETE FROM generics WHERE id = $1`, seed.GenericID)
	})

	return seed
}
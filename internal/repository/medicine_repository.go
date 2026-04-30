package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mahimtalukder/oshudh-somadhan/internal/model"
)

type MedicineRepository struct {
	DB *pgxpool.Pool
}

func NewMedicineRepository(db *pgxpool.Pool) *MedicineRepository {
	return &MedicineRepository{
		DB: db,
	}
}

func (r *MedicineRepository) ListMedicines(
	ctx context.Context,
	page int,
	limit int,
	genericID *int,
	companyID *int,
	dosageFormID *int,
) (*model.PeginatedMedicineResponse, error) {
	offset := (page - 1) * limit

	//set qeury whare
	whereParts := []string{"1=1"}
	args := []any{}
	argPos := 1

	if genericID != nil {
		whereParts = append(whereParts, fmt.Sprintf("m.generics_id = $%d", argPos))
		args = append(args, *genericID)
		argPos++
	}

	if companyID != nil {
		whereParts = append(whereParts, fmt.Sprintf("m.companie_id = $%d", argPos))
		args = append(args, *companyID)
		argPos++
	}

	if dosageFormID != nil {
		whereParts = append(whereParts, fmt.Sprintf("m.dosageform_id = $%d", argPos))
		args = append(args, *dosageFormID)
		argPos++
	}

	whareQuery := strings.Join(whereParts, " AND ")

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM medicine m
		WHERE %s`,
		whareQuery,
	)

	var total int
	if err := r.DB.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
		SELECT 
			m.id,
			COALESCE(m.name, ''),
			COALESCE(m.strength, ''),
			COALESCE(m.packsizeandprice, ''),
			COALESCE(g.name, ''),
			COALESCE(c.name, ''),
			COALESCE(d.name, '')
		FROM medicine m
		LEFT JOIN generics g ON g.id = m.generics_id
		LEFT JOIN companie c ON c.id = m.companie_id
		LEFT JOIN dosageform d ON d.id = m.dosageform_id
		WHERE %s
		ORDER BY m.id DESC
		LIMIT $%d OFFSET $%d
	`, whareQuery, argPos, argPos+1)

	args = append(args, limit, offset)

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	medicines := make([]model.MedicineListItem, 0)

	for rows.Next() {
		var item model.MedicineListItem
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Strength,
			&item.PackSizeAndPrice,
			&item.GenericName,
			&item.CompanyName,
			&item.DosageFormName,
		)

		if err != nil {
			return nil, err
		}

		medicines = append(medicines, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	totalPage := 0
	if total > 0 {
		totalPage = (total + limit - 1) / limit
	}

	return &model.PeginatedMedicineResponse{
		Data:       medicines,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPage,
	}, nil
}

func (r *MedicineRepository) GetMedicineByID(ctx context.Context, id int) (*model.MedicineDetails, error) {
	query := `
			SELECT 
			m.id,
			COALESCE(m.name, ''),
			COALESCE(m.strength, ''),
			COALESCE(m.packsizeandprice, ''),

			COALESCE(g.id, 0),
			COALESCE(g.name, ''),
			COALESCE(g.slug, ''),
			COALESCE(g.type, ''),

			COALESCE(c.id, 0),
			COALESCE(c.name, ''),

			COALESCE(d.id, 0),
			COALESCE(d.name, '')
		FROM medicine m
		LEFT JOIN generics g ON g.id = m.generics_id
		LEFT JOIN companie c ON c.id = m.companie_id
		LEFT JOIN dosageform d ON d.id = m.dosageform_id
		WHERE m.id = $1
	`

	var medicine model.MedicineDetails

	err := r.DB.QueryRow(ctx, query, id).Scan(
		&medicine.ID,
		&medicine.Name,
		&medicine.Strength,
		&medicine.PackSizeAndPrice,
		&medicine.GenericID,
		&medicine.GenericName,
		&medicine.GenericSlug,
		&medicine.GenericType,
		&medicine.CompanyID,
		&medicine.CompanyName,
		&medicine.DosageFormID,
		&medicine.DosageFormName,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, pgx.ErrNoRows
	}

	if err != nil {
		return nil, err
	}

	return &medicine, nil
}

func (r *MedicineRepository) SearchMedicines(ctx context.Context, q string, limit int, page int) (*model.PeginatedMedicineResponse, error) {
	offset := (page - 1) * limit

	countQuery := `
		SELECT 
		COUNT(*)
		FROM medicine m
		LEFT JOIN generics g ON g.id = m.generics_id
		LEFT JOIN companie c ON c.id = m.companie_id
		LEFT JOIN dosageform d ON d.id = m.dosageform_id
		WHERE LOWER(m.name) LIKE LOWER($1)
	`

	var total int
	if err := r.DB.QueryRow(ctx, countQuery, "%"+q+"%").Scan(&total); err != nil {
		return nil, err
	}

	totalPage := (total + limit - 1) / limit

	query := `
		SELECT 
			m.id,
			COALESCE(m.name, ''),
			COALESCE(m.strength, ''),
			COALESCE(m.packsizeandprice, ''),
			COALESCE(g.name, ''),
			COALESCE(c.name, ''),
			COALESCE(d.name, '')
		FROM medicine m
		LEFT JOIN generics g ON g.id = m.generics_id
		LEFT JOIN companie c ON c.id = m.companie_id
		LEFT JOIN dosageform d ON d.id = m.dosageform_id
		WHERE LOWER(m.name) LIKE LOWER($1)
		ORDER BY m.name ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.DB.Query(ctx, query, "%"+q+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	medicines := make([]model.MedicineListItem, 0)
	for rows.Next() {
		var medicine model.MedicineListItem
		err = rows.Scan(
			&medicine.ID,
			&medicine.Name,
			&medicine.Strength,
			&medicine.PackSizeAndPrice,
			&medicine.GenericName,
			&medicine.CompanyName,
			&medicine.DosageFormName,
		)
		if err != nil {
			return nil, err
		}
		medicines = append(medicines, medicine)
	}

	return &model.PeginatedMedicineResponse{
		Data:       medicines,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPage,
	}, rows.Err()
}

func (r *MedicineRepository) GetMedicinesByGenericID(ctx context.Context, genericId, limit, page int) (*model.PeginatedMedicineResponse, error) {
	return r.ListMedicines(ctx, page, limit, &genericId, nil, nil)
}

func (r *MedicineRepository) GetAlternatives(ctx context.Context, medicineId, limit, page int) (*model.PeginatedMedicineResponse, error) {
	offset := (page - 1) * limit

	var genericId int

	genericIdGetQuery := `
		SELECT 
			COALESCE(generics_id, 0)
		FROM medicine
		WHERE id = $1
	`
	if err := r.DB.QueryRow(ctx, genericIdGetQuery, medicineId).Scan(&genericId); err != nil {
		return nil, err
	}

	countQuery := `
		SELECT 
			COUNT(*)
		FROM medicine m
		LEFT JOIN generics g ON g.id = m.generics_id
		LEFT JOIN companie c ON c.id = m.companie_id
		LEFT JOIN dosageform d ON d.id = m.dosageform_id
		WHERE m.generics_id = $1 AND m.id != $2
	`

	var total int
	if err := r.DB.QueryRow(ctx, countQuery, genericId, medicineId).Scan(&total); err != nil {
		return nil, err
	}
	totalPages := (total + limit - 1) / limit

	query := `
		SELECT 
			m.id,
			COALESCE(m.name, ''),
			COALESCE(m.strength, ''),
			COALESCE(m.packsizeandprice, ''),
			COALESCE(g.name, ''),
			COALESCE(c.name, ''),
			COALESCE(d.name, '')
		FROM medicine m
		LEFT JOIN generics g ON g.id = m.generics_id
		LEFT JOIN companie c ON c.id = m.companie_id
		LEFT JOIN dosageform d ON d.id = m.dosageform_id
		WHERE m.generics_id = $1 AND m.id != $2
		ORDER BY m.name ASC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.DB.Query(ctx, query, genericId, medicineId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	alternatives := make([]model.MedicineListItem, 0)

	for rows.Next() {
		var item model.MedicineListItem
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Strength,
			&item.PackSizeAndPrice,
			&item.GenericName,
			&item.CompanyName,
			&item.DosageFormName,
		)
		if err != nil {
			return nil, err
		}
		alternatives = append(alternatives, item)
	}

	return &model.PeginatedMedicineResponse{
		Data:       alternatives,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, rows.Err()
}

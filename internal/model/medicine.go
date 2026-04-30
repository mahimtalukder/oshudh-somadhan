package model

type MedicineListItem struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Strength         string `json:"strength"`
	PackSizeAndPrice string `json:"pack_size_and_price"`
	GenericName      string `json:"generic_name"`
	CompanyName      string `json:"company_name"`
	DosageFormName   string `json:"dosage_form_name"`
}

type MedicineDetails struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Strength         string `json:"strength"`
	PackSizeAndPrice string `json:"pack_size_and_price"`

	GenericID   int    `json:"generic_id"`
	GenericName string `json:"generic_name"`
	GenericSlug string `json:"generic_slug"`
	GenericType string `json:"generic_type"`

	CompanyID   int    `json:"company_id"`
	CompanyName string `json:"company_name"`

	DosageFormID   int    `json:"dosage_form_id"`
	DosageFormName string `json:"dosage_form_name"`
}

type PeginatedMedicineResponse struct {
	Data       []MedicineListItem `json:"data"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	Total      int                `json:"total"`
	TotalPages int                `json:"total_pages"`
}

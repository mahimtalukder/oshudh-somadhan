CREATE INDEX IF NOT EXISTS idx_medicine_name ON medicine(name);
CREATE INDEX IF NOT EXISTS idx_medicine_generics_id ON medicine(generics_id);
CREATE INDEX IF NOT EXISTS idx_medicine_companie_id ON medicine(companie_id);
CREATE INDEX IF NOT EXISTS idx_medicine_dosageform_id ON medicine(dosageform_id);
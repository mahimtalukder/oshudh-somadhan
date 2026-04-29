package config

import "testing"

func TestConfigDatabaseUrl(t *testing.T){
	cfg := &Config{
		DbHost:     "localhost",
		DbPort:     "5433",
		DbUser:     "postgres",
		DbPassword: "postgres",
		DbName:     "oshudh_somadhan",
		DbSSLMode:  "disable",
	}

	got := cfg.DatabaseURL()

		want := "postgres://postgres:postgres@localhost:5433/oshudh_somadhan?sslmode=disable"

	if got != want {
		t.Errorf("expected %s, got %s", want, got)
	}
}
package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName string
	AppEnv  string
	AppPort string

	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
	DbSSLMode  string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	appName := 	os.Getenv("APP_NAME")
	appEnv 	:= 	os.Getenv("APP_ENV")
	appPort	:=	os.Getenv("APP_PORT")

	if appEnv == "" || appName == "" || appPort == "" {
		return nil, errors.New("'APP_ENV', 'APP_NAME', 'APP_PORT' can't be null")
	}

	dbHost  	:=	os.Getenv("DB_HOST")
	dbPort		:=  os.Getenv("DB_PORT")
	dbUser		:=	os.Getenv("DB_USER")
	dbPassword	:= 	os.Getenv("DB_PASSWORD")
	dbName		:=	os.Getenv("DB_NAME")
	dbSSLMode	:=  os.Getenv("DB_SSL_MODE")
	
	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbSSLMode == "" {
		return nil, errors.New("'DB_HOST', 'DB_PORT', 'DB_USER', 'DB_PASSWORD', 'DB_NAME', 'DB_SSL_MODE' can't be null")
	}
 
	return &Config{
		AppName: appName,
		AppEnv:  appEnv,
		AppPort: appPort,

		DbHost:     dbHost,
		DbPort:     dbPort,
		DbUser:     dbUser,
		DbPassword: dbPassword,
		DbName:     dbName,
		DbSSLMode:  dbSSLMode,
	}, nil
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DbUser,
		c.DbPassword,
		c.DbHost,
		c.DbPort,
		c.DbName,
		c.DbSSLMode,
	)
}

package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ConfigDatabase struct {
	Host     string `env:"DB_HOST" env-required:"true"`
	Port     string `env:"DB_PORT" env-required:"true"`
	Name     string `env:"DB_NAME" env-required:"true"`
	User     string `env:"DB_USER" env-required:"true"`
	Password string `env:"DB_PASSWORD" env-required:"true"`
}

func ReadConfigDB() (config ConfigDatabase, err error) {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found, loading from system environment variables.")
	}

	config = ConfigDatabase{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Name:     os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	}
	return config, nil
}

func BuildDB_URL(configDB ConfigDatabase) string {
	dbURL := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		configDB.User,
		configDB.Password,
		configDB.Host,
		configDB.Port,
		configDB.Name,
	)
	return dbURL
}

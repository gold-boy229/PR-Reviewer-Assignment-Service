package config

import (
	"fmt"
	"pr-reviewer-assignment-service/internal/consts"

	"github.com/ilyakaznacheev/cleanenv"
)

type ConfigDatabase struct {
	Host     string `env:"DB_HOST" env-required:"true"`
	Port     string `env:"DB_PORT" env-required:"true"`
	Name     string `env:"DB_NAME" env-required:"true"`
	User     string `env:"DB_USER" env-required:"true"`
	Password string `env:"DB_PASSWORD" env-required:"true"`
}

func ReadConfigDB() (config ConfigDatabase, err error) {
	pathToEnvFile := fmt.Sprintf("%v/.env", consts.PathToEnvFile)
	if err = cleanenv.ReadConfig(pathToEnvFile, &config); err != nil {
		return config, err
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

package boostrap

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvVar int

const (
	EnvDBHost EnvVar = iota
	EnvDBPort
	EnvDBName
	EnvDBUser
	EnvDBPassword
	EnvServerPort
)

var envVarNames = map[EnvVar]string{
	EnvDBHost:     "POSTGRES_HOST",
	EnvDBPort:     "POSTGRES_PORT",
	EnvDBName:     "POSTGRES_DB",
	EnvDBUser:     "POSTGRES_USER_NAME",
	EnvDBPassword: "POSTGRES_PASSWORD",
	EnvServerPort: "PORT",
}

func LoadEnvFile() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func GetEnv(envVar EnvVar) string {
	envVarName, ok := envVarNames[envVar]
	if !ok {
		log.Fatal("Can not get env")
		return ""
	}
	value := os.Getenv(envVarName)

	if value == "" {
		log.Fatal("Can not get env")
		return ""
	}

	return value

}

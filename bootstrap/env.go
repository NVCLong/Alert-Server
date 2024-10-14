package boostrap

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
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
	EnvAllowPort
	EnvRedisHost
	EnvRedisAccessKey
	EnvRedisTTL
)

var envVarNames = map[EnvVar]string{
	EnvDBHost:         "POSTGRES_HOST",
	EnvDBPort:         "POSTGRES_PORT",
	EnvDBName:         "POSTGRES_DB",
	EnvDBUser:         "POSTGRES_USER_NAME",
	EnvDBPassword:     "POSTGRES_PASSWORD",
	EnvServerPort:     "PORT",
	EnvAllowPort:      "ALLOW_PORT",
	EnvRedisHost:      "REDIS_HOST",
	EnvRedisAccessKey: "REDIS_ACCESS_KEY",
	EnvRedisTTL:       "REDIS_TTL",
}

func LoadEnvFile() {
	if gin.Mode() != gin.ReleaseMode {
		err := godotenv.Load()
		if err != nil {
			log.Printf("Error loading .env file: %v", err)
		}
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

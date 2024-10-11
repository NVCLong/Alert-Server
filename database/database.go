package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	bootstrap "github.com/NVCLong/Alert-Server/bootstrap"
)

var Db *gorm.DB

func ConnectDatabase() {
	bootstrap.LoadEnvFile()

	host := bootstrap.GetEnv(bootstrap.EnvDBHost)
	port := bootstrap.GetEnv(bootstrap.EnvDBPort)
	databaseName := bootstrap.GetEnv(bootstrap.EnvDBName)
	userName := bootstrap.GetEnv(bootstrap.EnvDBUser)
	password := bootstrap.GetEnv(bootstrap.EnvDBPassword)

	// Set up PostgreSQL connection string
	psqlSetup := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host, port, userName, databaseName, password,
	)

	// Connect to the PostgreSQL database
	db, errSql := gorm.Open(postgres.Open(psqlSetup), &gorm.Config{})
	if errSql != nil {
		log.Println("There is an error while connecting to the database:", errSql)
		log.Fatalf("Database connection failed: %v", errSql)
	} else {
		Db = db
		log.Println("Successfully connected to the database!")
	}
}

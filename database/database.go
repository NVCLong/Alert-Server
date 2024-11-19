package database

import (
	"fmt"
	"gorm.io/gorm/logger"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	bootstrap "github.com/NVCLong/Alert-Server/bootstrap"
	workflow "github.com/NVCLong/Alert-Server/models/workflow"
)

var Db *gorm.DB

func ConnectDatabase() *gorm.DB {
	bootstrap.LoadEnvFile()

	host := bootstrap.GetEnv(bootstrap.EnvDBHost)
	port := bootstrap.GetEnv(bootstrap.EnvDBPort)
	databaseName := bootstrap.GetEnv(bootstrap.EnvDBName)
	userName := bootstrap.GetEnv(bootstrap.EnvDBUser)
	password := bootstrap.GetEnv(bootstrap.EnvDBPassword)

	// Set up PostgreSQL connection string
	psqlSetup := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable prefer_simple_protocol=true",
		host, port, userName, databaseName, password,
	)

	// Connect to the PostgreSQL database
	db, err := gorm.Open(postgres.Open(psqlSetup), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Info),
		PrepareStmt: true,
	})
	if err != nil {
		log.Println("Error connecting to the database:", err)
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Println("Error getting database instance:", err)
		return nil
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	db.Exec("DEALLOCATE ALL")
	AutoMigrate(db)
	log.Println("Successfully connected to the database!")
	return db
}

func AutoMigrate(db *gorm.DB) {
	workflow.Migrate(db)
}

// repository interface
type AbstractRepository[T any] interface {
	Create(entity T) error
	FindByCondition(id uint) (T, error)
	Update(entity T) error
	Delete(entity T) error
}

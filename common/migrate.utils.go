package common

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB, tables []Model) {
	for _, table := range tables {
		tableName := table.TableName()

		var exists bool
		query := fmt.Sprintf(`SELECT EXISTS (
				SELECT 1 FROM information_schema.tables 
				WHERE table_schema = 'public' 
				AND table_name = '%s'
			);`, tableName)

		if err := db.Raw(query).Scan(&exists).Error; err != nil {
			log.Printf("Error checking existence of table %s: %v", tableName, err)
			continue
		}

		if !exists {
			log.Printf("Table for model %T does not exist. Creating...", table)
			if err := db.AutoMigrate(table); err != nil {
				log.Printf("Error migrating table for model %T: %v", table, err)
			} else {
				log.Printf("Table for model %T created successfully.", table)
			}
		} else {
			log.Printf("Table for model %T already exists. Attempting to update schema...", table)
			if err := db.AutoMigrate(table); err != nil {
				log.Printf("Error updating table for model %T: %v", table, err)
			} else {
				log.Printf("Table for model %T updated successfully.", table)
			}
		}
	}
}

type Model interface {
	TableName() string
}

package workflow

import (
	"github.com/NVCLong/Alert-Server/common"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	tables := []common.Model{
		&ConditionBatch{},
		&WorkFlow{},
	}

	common.Migrate(db, tables)
}

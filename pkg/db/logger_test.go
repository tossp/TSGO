package db_test

import (
	"github.com/tossp/tsgo/pkg/db"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Example() {
	logger := db.NewLogger(zap.L())
	logger.SetAsDefault() // optional: configure gorm to use this zapgorm.Logger for callbacks
	db, _ := gorm.Open(nil, &gorm.Config{Logger: logger})

	// do stuff normally
	var _ = db // avoid "unused variable" warn
}

package db

import (
	"fmt"
	"sync"
	"time"

	"github.com/tossp/tsgo/pkg/log"
	"github.com/tossp/tsgo/pkg/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"
)

var (
	g             *gorm.DB
	gormTableLock = new(sync.Mutex)
	dbGormTables  []interface{}
)

func StartGorm() (err error) {
	log.Info("开始初始化 gorm")
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return fmt.Sprintf("%s_%s", setting.DbPrefix(), defaultTableName)
	}
	db, err := gorm.Open("postgres", fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s&fallback_application_name=%s&Schema",
		setting.DbUser(), setting.DbPassword(), setting.DbHost(), setting.DbPort(), setting.DbName(), setting.DbMode(), setting.AppName,
	))
	if err != nil {
		panic("尝试连接数据库失败：" + fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=%s&fallback_application_name=%s&Schema",
			setting.DbUser(), "********", setting.DbHost(), setting.DbPort(), setting.DbName(), setting.DbMode(), setting.AppName,
		))
	}
	db.LogMode(true)
	db.DB().SetMaxIdleConns(setting.DbMaxIdleConns())
	db.DB().SetMaxOpenConns(setting.DbMaxOpenConns())
	db.DB().SetConnMaxLifetime(time.Hour)

	db.SetLogger(log.GetLogger())
	//db.LogMode(false)
	g = db

	go gPing()
	return
}

func G() *gorm.DB {
	return g
}

func autoMigrate() {
	for _, err := range g.Exec(`CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS postgis_topology;
CREATE EXTENSION IF NOT EXISTS fuzzystrmatch;
CREATE EXTENSION IF NOT EXISTS postgis_tiger_geocoder;
CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS hstore;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).GetErrors() {
		log.Warn("Gorm数据库插件集成失败", err)
	}
	//g.AutoMigrate(&Organisation{})
}

func gPing() {
	if err := g.DB().Ping(); err != nil {
		log.Error(errors.Wrap(err, "连接数据库测试失败"))
		return
	}
	autoMigrate()
	if err := gsync(); err != nil {
		log.Error(errors.Wrap(err, "gorm数据库同步失败"))
		return
	}
}

func AddGormTables(t ...interface{}) {
	if len(t) == 0 {
		return
	}
	gormTableLock.Lock()
	defer gormTableLock.Unlock()
	dbGormTables = append(dbGormTables, t...)
}

func gsync() (err error) {
	gormTableLock.Lock()
	defer gormTableLock.Unlock()

	if err = g.AutoMigrate(dbGormTables...).Error; err != nil {
		err = errors.Wrap(err, "同步数据库实体错误")
		return
	}
	return
}

//func ValidateUniq(fl validator.FieldLevel) bool {
//	var result struct{ Count int }
//	currentField, _, _ := fl.GetStructFieldOK()
//	table := modelTableNameMap[currentField.Type().Name()] // table name
//	value := fl.Field().String()                           // value
//	column := fl.FieldName()                               // column name
//	sql := fmt.Sprintf("select count(*) from %s where %s='%s'", table, column, value)
//	db.PG.Raw(sql).Scan(&result)
//	dup := result.Count > 0
//	return !dup
//}

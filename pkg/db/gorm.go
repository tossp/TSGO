package db

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/tossp/tsgo/pkg/log"
	"github.com/tossp/tsgo/pkg/setting"
	"go.uber.org/zap"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"
)

var (
	g             *gorm.DB
	gormTableLock = new(sync.RWMutex)
	dbGormTables  []interface{}
)

func StartGorm() (err error) {
	log.Info("初始化数据模型")
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return TableName(defaultTableName)
	}
	dialect := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s&fallback_application_name=%s&Schema",
		setting.DbUser(), setting.DbPassword(), setting.DbHost(), setting.DbPort(), setting.DbName(), setting.DbMode(), setting.AppName,
	)
	db, err := gorm.Open("postgres", dialect)
	if err != nil {
		panic("尝试连接数据库失败：" + strings.Replace(dialect, setting.DbPassword(), "******", -1))
	}
	db.DB().SetMaxIdleConns(setting.DbMaxIdleConns())
	db.DB().SetMaxOpenConns(setting.DbMaxOpenConns())
	db.DB().SetConnMaxLifetime(time.Minute * 15)
	//db.LogMode(true)
	db.SetLogger(newLog(log.Desugar().Named("db").WithOptions(zap.AddCallerSkip(6))))
	//db.LogMode(false)
	g = db
	go gPing()
	return
}

func G() *gorm.DB {
	return g
}
func TableName(defaultTableName string) string {
	return fmt.Sprintf("%s_%s", setting.DbPrefix(), defaultTableName)
}

func autoMigrate() {
	//CREATE EXTENSION IF NOT EXISTS postgis;
	//CREATE EXTENSION IF NOT EXISTS postgis_topology;
	//CREATE EXTENSION IF NOT EXISTS postgis_tiger_geocoder;
	//CREATE EXTENSION IF NOT EXISTS fuzzystrmatch;
	//CREATE EXTENSION IF NOT EXISTS citext;
	for _, err := range g.Exec(`CREATE EXTENSION IF NOT EXISTS hstore;CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).GetErrors() {
		log.Warnw("数据库插件集成失败", "err", err)
	}
	//g.AutoMigrate(&Organisation{})
}

func gPing() {
	if err := g.DB().Ping(); err != nil {
		log.Errorw("连接数据库测试失败", "err", err)
		return
	}
	autoMigrate()
	if err := gsync(); err != nil {
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
	gormTableLock.RLock()
	defer gormTableLock.RUnlock()

	if err = g.Debug().AutoMigrate(dbGormTables...).Error; err != nil {
		log.Errorw("同步数据库实体错误", "err", err)
		err = errors.Wrap(err, "同步数据库实体错误")
		return
	}
	return
}

type Logger struct {
	zap *zap.Logger
}

func newLog(logger *zap.Logger) Logger {
	return Logger{zap: logger}
}

func (l Logger) Print(values ...interface{}) {
	if len(values) < 2 {
		log.Warn("遗漏来源", values)
		return
	}

	switch values[0] {
	case "sql":
		l.zap.Debug("sql",
			zap.String("query", values[3].(string)),
			zap.Any("values", values[4]),
			zap.Duration("duration", values[2].(time.Duration)),
			zap.Int64("affected-rows", values[5].(int64)),
			zap.String("source", values[1].(string)), // if AddCallerSkip(6) is well defined, we can safely remove this field
		)
	default:
		l.zap.Debug("other",
			zap.Any("values", values[2:]),
			zap.String("source", values[1].(string)), // if AddCallerSkip(6) is well defined, we can safely remove this field
		)
	}
}

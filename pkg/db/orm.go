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
	"github.com/jinzhu/inflection"
)

var (
	g             *gorm.DB
	gormTableLock = new(sync.RWMutex)
	dbGormTables  []interface{}
)

func init() {
	setting.SetDefault("db.User", "ts")
	setting.SetDefault("db.Password", "123456")
	setting.SetDefault("db.Prefix", "ts")
	setting.SetDefault("db.Host", "127.0.0.1")
	setting.SetDefault("db.Port", 5432)
	setting.SetDefault("db.Name", "ts")
	setting.SetDefault("db.Ssl_mode", "disable")
	setting.SetDefault("db.Max_Idle_Conns", 10)
	setting.SetDefault("db.Max_Open_Conns", 20)
}

//StartGorm 启动GORM
func StartGorm() (err error) {
	log.Info("初始化数据模型")
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return TableName(defaultTableName)
	}
	dialect := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s&fallback_application_name=%s&Schema",
		setting.GetString("db.User"), setting.GetString("db.Password"), setting.GetString("db.Host"),
		setting.GetInt64("db.Port"), setting.GetString("db.Name"), setting.GetString("db.Ssl_mode"),
		setting.AppName+setting.GetString("name"),
	)
	db, err := gorm.Open("postgres", dialect)
	if err != nil {
		panic("尝试连接数据库失败：" + strings.Replace(dialect, setting.GetString("db.Password"), "******", -1) + err.Error())
	}
	db.DB().SetMaxIdleConns(setting.GetInt("db.Max_Idle_Conns"))
	db.DB().SetMaxOpenConns(setting.GetInt("db.Max_Open_Conns"))
	db.DB().SetConnMaxLifetime(time.Minute * 15)
	//db.LogMode(true)
	db.SetLogger(newLog(log.Desugar().Named("db").WithOptions(zap.AddCallerSkip(6))))
	//db.LogMode(false)
	g = db
	go gPing()
	return
}

//G GORM实例
func G() *gorm.DB {
	return g
}

//IsRecordNotFoundError GORM实例
func IsRecordNotFoundError(err error) bool {
	return gorm.IsRecordNotFoundError(err)
}

//TableName 计算表名
func TableName(name string) string {
	return fmt.Sprintf("tsl_%s_%s", strings.ToLower(setting.GetString("db.Prefix")), inflection.Plural(gorm.ToTableName(name)))
}

//TableName 计算表名
func ColumnName(name string) string {
	return gorm.ToColumnName(name)
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

//AddGormTables 添加同步表
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
	defer log.Info("同步数据库实体过程结束")
	if err = g.Debug().AutoMigrate(dbGormTables...).Error; err != nil {
		log.Errorw("同步数据库实体错误", "err", err)
		return
	}
	return
}

//Logger 日志实体
type Logger struct {
	zap *zap.Logger
}

func newLog(logger *zap.Logger) Logger {
	return Logger{zap: logger}
}

//Print 打印信息
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

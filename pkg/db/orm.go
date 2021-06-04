package db

import (
	"fmt"
	"strings"
	"sync"

	"github.com/jackc/pgtype"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/tossp/tsgo/pkg/errors"
	"github.com/tossp/tsgo/pkg/log"
	"github.com/tossp/tsgo/pkg/setting"
)

var (
	g             *gorm.DB
	gormTableLock = new(sync.RWMutex)
	Expr          = gorm.Expr
	dbGormTables  []interface{}
	naming        schema.NamingStrategy
	confStr       = ""
)

func makeConfStr() string {
	return setting.GetString("db.User") +
		setting.GetString("db.password") +
		setting.GetString("db.Prefix") +
		setting.GetString("db.Name") +
		setting.GetString("db.Ssl_mode") +
		setting.GetString("db.TimeZone") +
		setting.GetString("db.Max_Idle_Conns") +
		setting.GetString("db.Max_Open_Conns") +
		setting.GetString("db.Host") +
		setting.GetString("db.Port") +
		setting.GetString("name")
}

func autoMakeDB() (err error) {
	if confStr == makeConfStr() {
		return
	}
	return makeDB()
}
func init() {
	pgtype.IgnoreUndefined()
	setting.SetDefault("db.User", "ts")
	setting.SetDefault("db.Password", "123456")
	setting.SetDefault("db.Prefix", "ts")
	setting.SetDefault("db.Host", "127.0.0.1")
	setting.SetDefault("db.Port", 5432)
	setting.SetDefault("db.Name", "ts")
	setting.SetDefault("db.Ssl_mode", "disable")
	setting.SetDefault("db.TimeZone", "Asia/Shanghai")
	setting.SetDefault("db.Max_Idle_Conns", 10)
	setting.SetDefault("db.Max_Open_Conns", 20)
	_ = setting.Subscribe(autoMakeDB)
}

//Start 启动GORM
func Start() (err error) {
	err = makeDB()
	gPing()
	return
}
func makeDB() (err error) {
	log.Info("初始化数据模型")
	confStr = makeConfStr()
	naming = schema.NamingStrategy{TablePrefix: fmt.Sprintf("tsl_%s_", strings.ToLower(setting.GetString("db.Prefix")))}
	dialect := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s&application_name=%s&TimeZone=%s",
		setting.GetString("db.User"), setting.GetString("db.Password"), setting.GetString("db.Host"),
		setting.GetInt64("db.Port"), setting.GetString("db.Name"), setting.GetString("db.Ssl_mode"),
		setting.AppName()+setting.GetString("name"), setting.GetString("db.TimeZone"),
	)
	logger := NewLogger(log.Desugar())
	logger.SetAsDefault()
	db, err := gorm.Open(postgres.Open(dialect), &gorm.Config{
		NamingStrategy: naming,
		Logger:         logger,
	})
	if err != nil {
		panic("尝试连接数据库失败：" + strings.Replace(dialect, setting.GetString("db.Password"), "******", -1) + err.Error())
	}
	g = db
	return
}

//G GORM实例
func G() *gorm.DB {
	return g
}

//IsRecordNotFoundError GORM实例
func IsRecordNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

//TableName 计算表名
func TableName(name string) string {
	return naming.TableName(name)
}

//TableName 计算表名
func ColumnName(name string) string {
	return naming.ColumnName("", name)
}

func autoMigrate() {
	//CREATE EXTENSION IF NOT EXISTS postgis;
	//CREATE EXTENSION IF NOT EXISTS postgis_topology;
	//CREATE EXTENSION IF NOT EXISTS postgis_tiger_geocoder;
	//CREATE EXTENSION IF NOT EXISTS fuzzystrmatch;
	//CREATE EXTENSION IF NOT EXISTS citext;
	sqlDB, err := g.DB()
	if err != nil {
		log.Errorw("获取数据库失败", "err", err)
		return
	}
	if _, err := sqlDB.Exec(`CREATE EXTENSION IF NOT EXISTS hstore;CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`); err != nil {
		log.Warnw("数据库插件集成失败", "err", err)
	}
	//g.AutoMigrate(&Organisation{})
}

func gPing() {
	sqlDB, err := g.DB()
	if err != nil {
		log.Errorw("获取数据库失败", "err", err)
		return
	}
	if err = sqlDB.Ping(); err != nil {
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
	if err = g.AutoMigrate(dbGormTables...); err != nil {
		log.Errorw("同步数据库实体错误", "err", err)
		return
	}
	return
}

func DelAll() (err error) {
	gormTableLock.RLock()
	defer gormTableLock.RUnlock()
	defer log.Info("清除数据库实体过程结束")
	if err = g.Debug().Migrator().DropTable(dbGormTables...); err != nil {
		log.Errorw("清除数据库实体错误", "err", err)
		return
	}
	return
}

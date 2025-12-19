package dao

import (
	"go-mall/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _DbMaster *gorm.DB
var _DbSlave *gorm.DB

// DB 返回只读实例
func DB() *gorm.DB {
	return _DbSlave
}

// DBMaster 返回主库实例
func DBMaster() *gorm.DB {
	return _DbMaster
}

func init() {
	_DbMaster = initDB(config.Database.Master)
	_DbSlave = initDB(config.Database.Slave)
}

func initDB(options config.DBConnectOption) *gorm.DB {
	db, err := gorm.Open(mysql.Open(options.DSN), &gorm.Config{
		Logger: NewGormLogger(),
	})
	if err != nil {
		panic(err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(options.MaxIdleConns)
	sqlDB.SetMaxOpenConns(options.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(options.MaxLifeTime)
	if err = sqlDB.Ping(); err != nil {
		panic(err)
	}

	return db
}

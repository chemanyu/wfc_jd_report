package mysqldb

import (
	"dmp_distribution/core"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Config 结构体定义了需要读取的配置项
var adnDb *gorm.DB

// LoadConfig 从配置文件中读取配置项
func InitAdnMysql() {
	dsn := core.GetConfig().MYSQL_ADN_DB + "?charset=utf8&parseTime=True&loc=Local"
	var err error
	adnDb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// 获取底层 *sql.DB 并设置连接池
	sqlDB, err := adnDb.DB()
	if err != nil {
		panic("failed to get underlying *sql.DB: " + err.Error())
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(100)                 // 最大连接数
	sqlDB.SetMaxIdleConns(50)                  // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(10 * time.Second) // 连接最大存活时间

	// 添加连接池状态监控
	go monitorConnectionPool(sqlDB)
}

func GetAdnConnected() *gorm.DB {
	if adnDb == nil {
		log.Panic("Database connection is not initialized. Call InitAdnMysql() first")
	}
	return adnDb
}

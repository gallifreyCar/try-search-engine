package database

import (
	mysqlCfg "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DbOne *gorm.DB

// ConnectDB 连接数据库
func ConnectDB() (db *gorm.DB, err error) {
	// 连接配置
	cfg := mysqlCfg.Config{
		User:                 os.Getenv("MYSQL_USER"),
		Passwd:               os.Getenv("MYSQL_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("MYSQL_ADDR"),
		DBName:               os.Getenv("MYSQL_DBNAME"),
		AllowNativePasswords: true,
		ParseTime:            true,
		Loc:                  time.Local,
	}

	// 设置日志
	file, err := os.Create("sql.log")
	if err != nil {
		return nil, err
	}
	dbLogger := logger.New(log.New(file, "[dbOne]", log.LstdFlags), logger.Config{
		SlowThreshold:             time.Second * 4, // 慢 SQL 阈值
		Colorful:                  false,           // 禁用彩色打印
		IgnoreRecordNotFoundError: true,            // 忽略ErrRecordNotFound（记录未找到）错误
		LogLevel:                  logger.Info,     // 日志级别
	})
	// 连接数据库
	db, err = gorm.Open(mysql.Open(cfg.FormatDSN()), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InitDB() error {
	db, err := ConnectDB()
	if err != nil {
		return err
	}
	DbOne = db
	return nil
}

package main

import (
	"fmt"
	mysqlCfg "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"time"
)

func main() {
	//初始化配置
	initEnv()
	//开始爬取数据
	nextStep(time.Now())
	//阻塞
	select {}
}

// 循环爬取数据
func nextStep(startTime time.Time) {
	//获取数据库url
	var status []Status
	db, err := connectDB()
	if err != nil {
		return
	}
	db.Table("pages").Where("craw_done", 0).Order("id asc").Limit(100).Find(&status)
	fmt.Println("开始爬取", len(status), "条数据")
	fmt.Println("跑完一轮", time.Now().Unix()-startTime.Unix(), "秒")
	nextStep(time.Now())
	//todo 爬取数据
}

// 初始化变量
func initEnv() {

}

// Status 表结构
type Status struct {
	ID       uint      `gorm:"primaryKey"`
	Url      string    `gorm:"default:null"`
	Host     string    `gorm:"default:null"`
	CrawDone int       `gorm:"type:tinyint(1);default:0"`
	CrawTime time.Time `gorm:"default:'2001-01-01 00:00:01'"`
}

// 连接数据库
func connectDB() (db *gorm.DB, err error) {
	cfg := mysqlCfg.Config{
		User:                 os.Getenv("MYSQL_USER"),
		Passwd:               os.Getenv("MYSQL_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("MYSQL_ADDR"),
		DBName:               os.Getenv("MYSQL_DBNAME"),
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	db, err = gorm.Open(mysql.Open(cfg.FormatDSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

package mysql

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DbInstance *gorm.DB

// Init 初始化数据库
func InitMysql() {
	startingTime := time.Now().UTC()
	source := "%s:%s@tcp(%s)/%s?readTimeout=15000ms&writeTimeout=15000ms&charset=utf8&loc=Local&parseTime=true"
	user := os.Getenv("MYSQL_USERNAME")
	pwd := os.Getenv("MYSQL_PASSWORD")
	addr := os.Getenv("MYSQL_ADDRESS")
	dataBase := os.Getenv("MYSQL_DATABASE")

	if dataBase == "" {
		dataBase = "hackathon"
	}

	if user == "" {
		user = "liucx_211459"
	}

	if pwd == "" {
		pwd = "Liuchenxing123"
	}
	if addr == "" {
		addr = "111.62.122.91"
	}
	source = fmt.Sprintf(source, user, pwd, addr, dataBase)
	fmt.Println("start init mysql with ", source)

	db, err := gorm.Open(mysql.Open(source), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		}})
	if err != nil {
		fmt.Println("DB Open error,err=", err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("DB Init error,err=", err.Error())
	}

	// 用于设置连接池中空闲连接的最大数量
	sqlDB.SetMaxIdleConns(100)
	// 设置打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(200)
	// 设置了连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Hour)

	DbInstance = db

	endingTime := time.Now().UTC()
	fmt.Println("finish init mysql with ", source, "duration is: ", endingTime.Sub(startingTime))
}

// Get dbInstance
func GetMysql() *gorm.DB {
	return DbInstance
}

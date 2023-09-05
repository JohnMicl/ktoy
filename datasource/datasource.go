package datasource

import (
	"ktoy/config"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}

func Init() error {
	var err error
	path := strings.Join([]string{config.Config.DbInfo.UserName, ":", config.Config.DbInfo.Password, "@(",
		config.Config.DbInfo.Ip, ":", strconv.Itoa(config.Config.DbInfo.Port), ")/", config.Config.DbInfo.DBName, "?charset=utf8&parseTime=true&loc=Local"}, "")
	db, err = gorm.Open("mysql", path)
	if err != nil {
		return err
	}
	db.DB().SetConnMaxLifetime(5 * time.Minute)
	//最大打开的连接数
	db.DB().SetMaxIdleConns(20)
	//设置最大闲置个数
	db.DB().SetMaxOpenConns(2000)
	//表生成结尾不带s
	db.SingularTable(true)
	// 启用Logger，显示详细日志
	db.LogMode(true)

	err = createTableIfNotExist()
	return err
}

func createTableIfNotExist() error {
	db := GetDB().AutoMigrate()
	return db.Error
}

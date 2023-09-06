package datasource

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"

	"ktoy/config"
	"ktoy/models"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}

func Init() error {
	var err error
	masterDSN := GenerageDSN(config.Config.DbInfo)
	db, err = gorm.Open(mysql.Open(masterDSN), &gorm.Config{})
	if err != nil {
		return err
	}

	// 启动读写分离Separation=true
	if config.Config.Separation {
		err = registrySlaveDB(db, masterDSN)
		if err != nil {
			return err
		}
	}

	err = createTableIfNotExist()
	return err
}

func createTableIfNotExist() error {
	db := GetDB().AutoMigrate(
		&models.User{},
	)
	return db
}

func registrySlaveDB(db *gorm.DB, masterDSN string) error {
	if config.Config.Slave == nil || len(config.Config.Slave) == 0 {
		return errors.New("slave config not specify")
	}
	replicas := []gorm.Dialector{}
	for _, v := range config.Config.Slave {
		cfg := mysql.Config{
			DSN: GenerageDSN(v),
		}
		replicas = append(replicas, mysql.New(cfg))
	}
	db.Use(dbresolver.Register(dbresolver.Config{
		Sources: []gorm.Dialector{mysql.New(mysql.Config{
			DSN: masterDSN,
		})},
		Replicas: replicas,
		Policy:   dbresolver.RandomPolicy{},
	}).
		SetMaxIdleConns(10).
		SetConnMaxLifetime(time.Hour).
		SetMaxOpenConns(200))
	return nil
}

func GenerageDSN(conf config.DbConfigs) string {
	return strings.Join([]string{conf.UserName, ":", conf.Password, "@tcp(",
		config.Config.DbInfo.Ip, ":", strconv.Itoa(config.Config.DbInfo.Port), ")/", config.Config.DbInfo.DBName, "?charset=utf8&parseTime=true&loc=Local"}, "")
}

package main

import (
	"flag"
	"ktoy/config"
	"ktoy/datasource"
	"ktoy/logger"
	"ktoy/utils"

	log "github.com/sirupsen/logrus"
)

var confYaml = flag.String("c", "conf.yaml", "config file name")

func main() {
	logger.Loginit()
	log.Info("welcome to ktoy :)")

	// 读取配置文件
	flag.Parse()
	utils.ParseYamlFile(*confYaml, &config.Config)
	log.WithFields(log.Fields{
		"config": config.Config,
	}).Info("ktoy get config file")

	// 连接ktoy数据库
	err := datasource.Init()
	if err != nil {
		log.Panic(err)
	}

	log.Info("ktoy go :)")
}

package config

var (
	Config Configs
)

type Configs struct {
	DbInfo     DbConfigs   `yaml:"dbInfo"`
	Separation bool        `yaml:"separation"`
	Slave      []DbConfigs `yaml:"slavedb"`
}

type DbConfigs struct {
	Ip       string `yaml:"ip"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"dbName"`
	UserName string `yaml:"dbUserName"`
	Password string `yaml:"dbPassword"`
}

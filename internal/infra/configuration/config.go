package configuration

import "github.com/spf13/viper"

var config = viper.New()

func Init() *viper.Viper {
	config.SetDefault("server.host", "0.0.0.0:8080")
	config.SetDefault("database.url", "postgres://admin:admin@localhost:5432/rinha")
	return config
}

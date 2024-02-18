package configuration

import "github.com/spf13/viper"

var config = viper.New()

func Init() *viper.Viper {
	//config.SetDefault("app.name", "rinha-de-backend")
	config.SetDefault("server.host", "0.0.0.0:8080")
	return config
}

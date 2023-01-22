package config

import "github.com/spf13/viper"

type Config struct {
	DBLogin  string `mapstructure:"DB_LOGIN"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBPath string `mapstructure:"DB_PATH"`
	ApiKey string `mapstructure:"API_KEY"`
	// DB int `mapstructure:"POSTGRES_PORT"`
}

func LoadConfig(prod bool) (c Config, err error) {
	viper.AddConfigPath(".")
	if prod {
		viper.SetConfigName("prod")
	} else {
		viper.SetConfigName("dev")
	}
	
	viper.SetConfigType("env")

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&c)

	return
}
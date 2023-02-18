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
	v := viper.New()
	v.AddConfigPath(".")

	if prod {
		v.SetConfigName("prod")
	} else {
		v.SetConfigName("dev")
	}
	
	v.SetConfigType("env")

	err = v.ReadInConfig()

	if err != nil {
		return
	}

	err = v.Unmarshal(&c)

	return
}
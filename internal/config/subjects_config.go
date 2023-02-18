package config

import "github.com/spf13/viper"

type SubjectsConfig struct {
	Subjects []struct {
		Id      int      `yaml:"id"`
		Name    string   `yaml:"name"`
		Alias   string   `yaml:"alias"`
		Classes []string `yaml:"classes"`
	} `yaml:"subjects"`
}

func LoadSubjectsConfig() (c SubjectsConfig, err error) {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigName("subjects")
	v.SetConfigType("yaml")

	err = v.ReadInConfig()

	if err != nil {
		return
	}

	err = v.Unmarshal(&c)

	return
}

package config

import "github.com/spf13/viper"

type Config struct {
	Port      string `mapstructure:"PORT"`
	DbURL     string `mapstructure:"DB_URL"`
	DbMigrate bool   `mapstructure:"DB_AUTO_MIGRATE"`
	JwtSecret string `mapstructure:"JWT_SECRET"`
}

func LoadConfig() (config Config, err error) {

	viper.AddConfigPath(".") // Look for config in the root directory
	viper.SetConfigFile(".env")
	viper.AutomaticEnv() // Read from system ENV if available

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

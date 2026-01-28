package config

import "github.com/spf13/viper"

type Config struct {
	Port      string `mapstructure:"PORT"`
	JwtSecret string `mapstructure:"JWT_SECRET"`

	DbURL     string `mapstructure:"DB_URL"`
	DbMigrate bool   `mapstructure:"DB_AUTO_MIGRATE"`

	KafkaBrokers []string `mapstructure:"KAFKA_BROKERS"`
	KafkaRequiredAcks string `mapstructure:"KAFKA_REQUIRED_ACKS"`
	KafkaAsync bool `mapstructure:"KAFKA_ASYNC"`
	KafkaTimeout int `mapstructure:"KAFKA_TIMEOUT"`
}

func LoadConfig() (config Config, err error) {
	// Set default
	viper.SetDefault("PORT", "8080")
    viper.SetDefault("KAFKA_REQUIRED_ACKS", "all")
    viper.SetDefault("KAFKA_ASYNC", false)
    viper.SetDefault("KAFKA_TIMEOUT", 10)
    viper.SetDefault("DB_AUTO_MIGRATE", true)

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

package config

import (
	

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config Структура конфигурации;
// Содержит все конфигурационные данные о сервисе;
// автоподгружается при изменении исходного файла
type Config struct {
	ServerAddress string `mapstructure:"SERVER_ADDRESS"` // Новый параметр для адреса и порта
}

// NewConfig Создаёт новый объект конфигурации, загружая данные из файла конфигурации и переменных окружения
func NewConfig() (*Config, error) {
	var err error

	// Загружаем переменные окружения из .env
	_ = godotenv.Load() // Загружаем переменные окружения из файла .env

	viper.AutomaticEnv() // Автоматически загружаем переменные окружения

	// Читаем значение SERVER_ADDRESS из .env
	viper.BindEnv("SERVER_ADDRESS", "SERVER_ADDRESS")

	// Устанавливаем значение по умолчанию для SERVER_ADDRESS
	viper.SetDefault("SERVER_ADDRESS", "0.0.0.0:8080")

	cfg := &Config{}
	err = viper.Unmarshal(cfg) // Распаковка значений в структуру Config
	if err != nil {
		return nil, err
	}

	log.Info("config parsed")
	return cfg, nil
}
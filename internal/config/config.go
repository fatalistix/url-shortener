package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
	Service     `yaml:"service"`
}

type HTTPServer struct {
	Address                 string        `yaml:"address" env-default:"localhost:8080"`
	Timeout                 time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout             time.Duration `yaml:"idle_timeout" env-default:"60s"`
	GracefulShutdownTimeout time.Duration `yaml:"graceful_shutdown_timeout" env-default:"10s"`
	User                    string        `yaml:"user" env-required:"true"`
	Password                string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

type Service struct {
	AliasLength int `yaml:"alias_length" env-required:"true"`
}

// Приставка Must по соглашению должна не возвращать ошибку в случае проблем, а паниковать. Так делать по-хорошему надо очень редко. Также некоторые игнорируют семантическое значение приставки Must и просто говорят, что она выдает ошибку. Здесь это оправданно, так как что ещё делать, если не падать, ведь даже конфиг не загрузился
func MustLoad(configPath string) *Config {
	if configPath == "" {
		// Если возникла проблема (то есть переменная окружения просто не была найдена), то просто падаем через `log.Fatal`, так как основной логгер `slog` просто еще не запущен
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exists: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}

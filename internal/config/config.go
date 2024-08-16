package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Config - структура для хранения конфигурации
type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		Sslmode  string `yaml:"sslmode"`
	} `yaml:"database"`
	Logger struct {
		Output   string `yaml:"output"`
		FilePath string `yaml:"filepath"`
		Level    string `yaml:"level"`
	} `yaml:"logger"`
	Jwt struct {
		JwtKey string `yaml:"jwtkey"`
	} `yaml:"jwt"`
	Migration struct {
		MigrationPath string `yaml:"migrationpath"`
	} `yaml:"migration"`
}

// LoadConfig - функция для загрузки конфигурации из YAML файла
func LoadConfig(filename string) *Config {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("cannot read config file: %v", err)
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("cannot unmarshal config data: %v", err)
	}
	return &cfg
}

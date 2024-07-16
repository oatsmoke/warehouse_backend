package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Logger   string  `json:"logger"`
	Port     string  `json:"port"`
	TokenTTL string  `json:"token_ttl"`
	DB       *DB     `json:"db"`
	Client   *Client `json:"client"`
}

type DB struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	DBName   string `json:"db_name"`
}

type Client struct {
	Protocol string `json:"protocol"`
	Ip       string `json:"ip"`
	Port     string `json:"port"`
}

// Init is initializing the configuration file
func Init() *Config {
	configData, err := os.ReadFile("cmd/config.json")
	if err != nil {
		log.Panic(err)
	}

	var config *Config
	err = json.Unmarshal(configData, &config)
	if err != nil {
		log.Panic(err)
	}

	return config
}

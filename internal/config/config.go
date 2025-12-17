package config

import (
	"errors"
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env  string     `yaml:"env"`
	Bot  BotConfig  `yaml:"bot"`
	Mono MonoConfig `yaml:"mono"`
	DB   DBConfig   `yaml:"db"`
}

type BotConfig struct {
	Token      string `yaml:"token" env:"BOT_TOKEN"`
	LongPoller int    `yaml:"long_poller"`
	Password   string `yaml:"password" env:"BOT_PASSWORD"`
}

type MonoConfig struct {
	EncryptKey string `yaml:"encrypt_key" env:"MONO_ENCRYPT_KEY"`
	ApiURL     string `yaml:"api_url"`
}

type DBConfig struct {
	User string `yaml:"user"`
	Pass string `yaml:"pass" env:"DB_PASS"`
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Name string `yaml:"name"`
}

func MustLoad() Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed to read config path: " + err.Error())
	}

	return cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

func ConvertTokenKeyToBytes(k string) ([]byte, error) {
	kbytes := []byte(k)

	if len(kbytes) != 32 {
		return nil, errors.New("encryption key must be exactly 32 bytes for AES-256")
	}

	return kbytes, nil
}

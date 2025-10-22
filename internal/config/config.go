package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env  string      `yaml:"env" env-default:"local"`
	Bot  *BotConfig  `yaml:"bot"`
	Mono *MonoConfig `yaml:"mono"`
	DB   *DBConfig   `yaml:"db"`
}

type BotConfig struct {
	Token      string `yaml:"token"`
	LongPoller int    `yaml:"long_poller"`
	Password   string `yaml:"password"`
}

type MonoConfig struct {
	EncryptKey string `yaml:"encrypt_key"`
	ApiURL     string `yaml:"api_url"`
}

type DBConfig struct {
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Name string `yaml:"name"`
}

func MustLoad() *Config {
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

	if token := os.Getenv("BOT_TOKEN"); token != "" {
		cfg.Bot.Token = token
	}
	if pass := os.Getenv("BOT_PASSWORD"); pass != "" {
		cfg.Bot.Password = pass
	}
	if dbPass := os.Getenv("DB_PASS"); dbPass != "" {
		cfg.DB.Pass = dbPass
	}
	if monoKey := os.Getenv("MONO_ENCRYPT_KEY"); monoKey != "" {
		cfg.Mono.EncryptKey = monoKey
	}

	return &cfg
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

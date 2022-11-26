package config

import (
	"time"

	"github.com/spf13/viper"
)

type App struct {
	Env         string   `mapstructure:"ENV_ID"`
	Server      Server   `mapstructure:"Server"`
	Database    Database `mapstructure:"DATABASE"`
	Environment string   `mapstructure:"ENVIRONMENT"`
	Log         Log      `mapstructure:"LOG"`
}

type Server struct {
	Http Http `mapstructure:"Http"`
}

type Http struct {
	Host            string        `mapstructure:"HOST"`
	Port            string        `mapstructure:"PORT"`
	ShutdownTimeout time.Duration `mapstructure:"SHUTDOWN_TIMEOUT"`
}

type Database struct {
	Name     string `mapstructure:"NAME"`
	User     string `mapstructure:"USER"`
	Password string `mapstructure:"PASSWORD" json:"-"`
	Host     string `mapstructure:"HOST"`
	Port     string `mapstructure:"PORT"`
}

type Log struct {
	Level     string `mapstructure:"LEVEL"`
	Formatter string `mapstructure:"FORMATTER"`
}

func Load(path string) (*App, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil {
		return nil, err
	}

	app := App{}

	return &app, viper.Unmarshal(&app)
}

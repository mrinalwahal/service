package main

import (
	"fmt"

	"github.com/spf13/viper"
)

// Base configuration.
type config struct {
	Environment    *environment    `mapstructure:"environment"`
	Database       *database       `mapstructure:"database"`
	Authentication *authentication `mapstructure:"authentication"`
}

// Environment configuration.
type environment struct {
	Environment string `mapstructure:"environment"`
	Debug       bool   `mapstructure:"debug"`
}

// Database configuration.
type database struct {
	Engine string `mapstructure:"engine"`
	DSN    string `mapstructure:"dsn"` // Data Source Name
}

// Authentication configuration.
type authentication struct {
	Method string `mapstructure:"method"`
	Key    struct {
		Algorithm string `mapstructure:"algorithm"`
		Key       string `mapstructure:"key"`
	} `mapstructure:"key"`
}

var c config

func Get() *config {
	return &c
}

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("unable to read config file, %v", err))
	}
	err := viper.Unmarshal(&c)
	if err != nil {
		panic(fmt.Sprintf("unable to decode into struct, %v", err))
	}
}

func main() {
	fmt.Println(Get().Database)
}

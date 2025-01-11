package config

import (
	"errors"
	"flag"
	"time"

	"github.com/caarlos0/env"
)

type Config struct {
	RunAddr                 string        `env:"RUN_ADDRESS"`
	AccrualSysAddr          string        `env:"ACCRUAL_SYSTEM_ADDRESS"`
	DatabaseURI             string        `env:"DATABASE_URI"`
	LogLvl                  string        `env:"LOGLVL"`
	TokenSecretKey          string        `env:"TOKEN_KEY"`
	TokenLifetime           time.Duration `env:"TOKEN_LIFETIME"`
	AccrualRateLimit        int           `env:"RATE_LIMIT"`
	OrderInfoUpdateInterval time.Duration `env:"ORDER_UPDATE_INTERVAL"`
}

func Parse() (*Config, error) {
	defaultValues := GetDefault()

	conf := Config{}
	flag.StringVar(&conf.RunAddr, "a", defaultValues.RunAddr, "address and port to run service")
	flag.StringVar(&conf.AccrualSysAddr, "r", defaultValues.AccrualSysAddr, "accrual system address")
	flag.StringVar(&conf.DatabaseURI, "d", defaultValues.DatabaseURI, "database connection string")
	flag.StringVar(&conf.LogLvl, "l", defaultValues.LogLvl, "application logging level")
	flag.StringVar(&conf.TokenSecretKey, "k", defaultValues.TokenSecretKey, "secret key to handle authentication")
	flag.DurationVar(&conf.TokenLifetime, "t", defaultValues.TokenLifetime, "authentication token lifetime")
	flag.IntVar(&conf.AccrualRateLimit, "rl", defaultValues.AccrualRateLimit, "accrual requests rate limit")
	flag.DurationVar(&conf.OrderInfoUpdateInterval, "o", defaultValues.OrderInfoUpdateInterval,
		"order info update interval")
	flag.Parse()

	env.Parse(&conf)

	if conf.AccrualSysAddr == "" {
		return nil, errors.New("empty value for accrual system address")
	}

	return &conf, nil
}

func GetDefault() (conf *Config) {
	return &Config{
		RunAddr:                 ":8080",
		AccrualSysAddr:          "localhost:8081",
		DatabaseURI:             "",
		LogLvl:                  "Info",
		TokenSecretKey:          "SECRET_KEY",
		TokenLifetime:           8 * time.Hour,
		AccrualRateLimit:        2,
		OrderInfoUpdateInterval: 30 * time.Second,
	}
}

func (c *Config) NormilizedAccrualSysAddr() string {
	return "http://" + c.AccrualSysAddr
}

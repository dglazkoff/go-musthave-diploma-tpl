package config

import (
	"flag"
	"os"
)

type Config struct {
	RunAddr              string
	AccrualSystemAddress string
	DatabaseURI          string
}

func ParseConfig() Config {
	cfg := Config{}

	flag.StringVar(&cfg.RunAddr, "a", ":8080", "address of the server")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", "", "address of the accrual system")
	flag.StringVar(&cfg.DatabaseURI, "d", "", "database uri string")
	flag.Parse()

	if runAddr := os.Getenv("RUN_ADDRESS"); runAddr != "" {
		cfg.RunAddr = runAddr
	}

	if accrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); accrualSystemAddress != "" {
		cfg.AccrualSystemAddress = accrualSystemAddress
	}

	if databaseURI := os.Getenv("DATABASE_URI"); databaseURI != "" {
		cfg.DatabaseURI = databaseURI
	}

	return cfg
}

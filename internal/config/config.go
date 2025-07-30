package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabasePath      string `env:"DATABASE_PATH"`
	PharmaciesCSVPath string `env:"PHARMACIES_CSV_PATH"`
	ClaimsDataPath    string `env:"CLAIMS_DATA_PATH"`
	RevertsDataPath   string `env:"REVERTS_DATA_PATH"`
	AuthToken         string `env:"AUTH_TOKEN"`
	Port              string `env:"PORT"`
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file. Using existing environment variables: %v", err)
	}

	cfg := &Config{
		DatabasePath:      os.Getenv("DATABASE_PATH"),
		PharmaciesCSVPath: os.Getenv("PHARMACIES_CSV_PATH"),
		ClaimsDataPath:    os.Getenv("CLAIMS_DATA_PATH"),
		RevertsDataPath:   os.Getenv("REVERTS_DATA_PATH"),
		AuthToken:         os.Getenv("AUTH_TOKEN"),
		Port:              os.Getenv("PORT"),
	}

	if cfg.DatabasePath == "" {
		cfg.DatabasePath = "./data/pharmacy.db"
		log.Printf("DATABASE_PATH not defined, using default: %s", cfg.DatabasePath)
	}
	if cfg.PharmaciesCSVPath == "" {
		cfg.PharmaciesCSVPath = "pharmacies.csv"
		log.Printf("PHARMACIES_CSV_PATH not defined, using default: %s", cfg.PharmaciesCSVPath)
	}
	if cfg.ClaimsDataPath == "" {
		cfg.ClaimsDataPath = "./data/claims"
		log.Printf("CLAIMS_DATA_PATH not defined, using default: %s", cfg.ClaimsDataPath)
	}
	if cfg.RevertsDataPath == "" {
		cfg.RevertsDataPath = "./data/reverts"
		log.Printf("REVERTS_DATA_PATH not defined, using default: %s", cfg.RevertsDataPath)
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
		log.Printf("PORT not defined, using default: %s", cfg.Port)
	}
	if cfg.AuthToken == "" {
		log.Println("Warning: AUTH_TOKEN not defined. Authentication might not work correctly.")
	}

	return cfg, nil
}

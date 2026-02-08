package config

import (
	"database/sql"
	"errors"
	"gopkg.in/yaml.v3"
	_ "modernc.org/sqlite"
	"os"
	"strings"
)

const defaultConfigPath = "config/config.yaml"

type Config struct {
	TG *Tg       `yaml:"tg"`
	DB *DbConfig `yaml:"db"`
}

type Tg struct {
	Token string `yaml:"token"`
}

type DbConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func MustLoadConfig(path string) *Config {
	if path == "" {
		path = defaultConfigPath
	}

	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var cfg Config
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}

	err = validateConfig(cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}

func validateConfig(cfg Config) error {
	if cfg.TG != nil {
		if cfg.TG.Token == "" {
			return errors.New("токен телеграм не передан")
		}
	}
	if cfg.DB != nil {
		dbErr := make([]string, 0)
		if cfg.DB.Driver == "" {
			dbErr = append(dbErr, "драйвер")
		}
		if cfg.DB.Host == "" {
			dbErr = append(dbErr, "хост")
		}
		if cfg.DB.Port == "" {
			dbErr = append(dbErr, "порт")
		}
		if cfg.DB.Username == "" {
			dbErr = append(dbErr, "пользователь")
		}
		if cfg.DB.Password == "" {
			dbErr = append(dbErr, "пароль")
		}
		if cfg.DB.Database == "" {
			dbErr = append(dbErr, "имя БД")
		}
		if len(dbErr) > 0 {
			return errors.New("Не указаны обязательные параметры БД: " + strings.Join(dbErr, ","))
		}
	}
	return nil

}

func MigrateDb(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS subscribers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT NOT NULL,
		address TEXT NOT NULL,
		channel_id TEXT,
		sign TEXT NOT NULL,
		created_at DATETIME DEFAULT  CURRENT_TIMESTAMP,
		status INTEGER,
        UNIQUE(type, address)
	);
	`

	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
}

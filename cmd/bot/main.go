package main

import (
	"database/sql"
	"fmt"
	"horoscope/internal/config"
	"horoscope/internal/repositories/sqlrepo"
	"horoscope/internal/telegram"
	"log"
)

func main() {
	cfg := config.MustLoadConfig("internal/config/config.yaml")
	service := telegram.NewService(cfg.TG.Token, sqlrepo.NewSubscriberRepo(NewDb(cfg.DB)))

	err := service.ProcessUpdates()
	if err != nil {
		log.Fatal(err)
	}
}

func NewDb(dbConfig *config.DbConfig) *sql.DB {
	db, err := sql.Open(dbConfig.Driver, fmt.Sprintf("file:%s.db?_auth&_auth_user=%s&_auth_pass=%s",
		dbConfig.Database, dbConfig.Username, dbConfig.Password))
	if err != nil {
		panic(err)
	}
	config.MigrateDb(db)
	return db
}

package main

import (
	"horoscope/internal/repositories/sqlrepo"
	"horoscope/internal/telegram"
	"log"
	"os"
)

func main() {
	token := os.Args[1]
	if token == "" {
		log.Fatalln("Токен телеграм-бота не передан!")
	}
	service := telegram.NewService(token, sqlrepo.NewSubscriberRepo(nil))

	err := service.ProcessUpdates()
	if err != nil {
		log.Fatal(err)
	}
}

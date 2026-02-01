package main

import (
	"fmt"
	"horoscope/internal/telegram"
	"log"
	"os"
)

func main() {
	token := os.Args[1]
	if token == "" {
		log.Fatalln("Токен телеграм-бота не передан!")
	}
	service := telegram.NewService(token)

	updates, err := service.GetUpdates()
	if err != nil {
		return
	}
	for _, update := range *updates {
		fmt.Println(update, " ")
	}
	fmt.Println()
}

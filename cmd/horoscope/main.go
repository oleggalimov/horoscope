package main

import (
	"horoscope/internal/horoscope"
	"horoscope/internal/providers/radioc"
)

func main() {
	providers := []horoscope.Provider{radioc.NewRadioCiProvider()}
	var senders []horoscope.Sender
	var subsDao horoscope.SubscribersDao
	service := horoscope.NewService(subsDao, senders, providers)

	service.SendDailyForecast()

}

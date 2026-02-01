package main

import (
	"fmt"
	"horoscope/internal/horoscope"
	"horoscope/internal/horoscope/model"
	"horoscope/internal/providers/radioc"
)

func main() {
	providers := []horoscope.Provider{radioc.NewRadioCiProvider()}
	forecasts := make([]model.DailyHoroscope, 0)

	for _, provider := range providers {
		f, e := provider.GetDailyHoroscope()
		if e != nil {
			fmt.Println("Не удалось получить гороскоп: %w", e)
			continue
		}
		forecasts = append(forecasts, f)
	}

	fmt.Println("Полученные горскопы: ")
	fmt.Println(forecasts)

}

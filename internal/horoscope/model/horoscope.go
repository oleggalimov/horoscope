package model

import "time"

type DailyHoroscope struct {
	Source    string
	Date      time.Time
	Forecasts []ZodiacForecast
}

type ZodiacForecast struct {
	Zodiac ZodiacSign
	Text   string
}

type ZodiacSign string

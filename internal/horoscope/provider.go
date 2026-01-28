package horoscope

import "horoscope/internal/horoscope/model"

type Provider interface {
	GetDailyHoroscope() (model.DailyHoroscope, error)
}

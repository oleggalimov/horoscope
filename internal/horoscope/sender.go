package horoscope

import "horoscope/internal/horoscope/model"

type Sender interface {
	TargetChannel() model.Channel
	Send([]model.DailyHoroscope)
}

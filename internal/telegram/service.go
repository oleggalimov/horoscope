package telegram

import "horoscope/internal/telegram/model"

type Service struct {
	bot *bot
}

func NewService(token string) *Service {
	return &Service{newBot(token)}
}

func (s *Service) GetUpdates() (*[]model.Update, error) {
	return s.bot.getUpdates()

}

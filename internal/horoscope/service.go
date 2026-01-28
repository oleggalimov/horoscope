package horoscope

import (
	"context"
	"fmt"
	"horoscope/internal/horoscope/model"
	"log"
)

type Service struct {
	providers  []Provider
	subsDao    SubscribersDao
	sendersMap map[model.Channel]Sender
}

func NewService(subsDao SubscribersDao, senders []Sender, providers []Provider) *Service {
	sendersMap := make(map[model.Channel]Sender)
	for _, sender := range senders {
		sendersMap[sender.TargetChannel()] = sender
	}
	return &Service{
		providers:  providers,
		subsDao:    subsDao,
		sendersMap: sendersMap,
	}
}

func (s *Service) SendDailyForecast() {
	//собираем гороскоп
	horoscopes := make([]model.DailyHoroscope, 0)
	for _, p := range s.providers {
		h, err := p.GetDailyHoroscope()
		if err != nil {
			continue
		}
		horoscopes = append(horoscopes, h)
	}
	//собираем подписчиков
	subs, err := s.subsDao.FindAll(context.Background())
	if err != nil {
		log.Println("Ошибка при получении списка подписчиков", err)
		return
	}
	//отправляем
	for _, subscriber := range subs {
		if sender, ok := s.sendersMap[subscriber.Type]; ok {
			sender.Send(horoscopes)
		} else {
			fmt.Printf("Для подписчика %d не определен канал! \n", subscriber.ID)
		}
	}
}

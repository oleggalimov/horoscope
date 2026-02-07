package telegram

import (
	"fmt"
	_ "go/parser"
	"horoscope/internal/horoscope"
	"horoscope/internal/horoscope/model"
	"horoscope/internal/telegram/commands"
	tgmodel "horoscope/internal/telegram/model"
	"math"
)

type Service struct {
	bot *bot
	dao horoscope.SubscribersDao
}

func NewService(token string, dao horoscope.SubscribersDao) *Service {
	return &Service{newBot(token), dao}
}

func (s *Service) ProcessUpdates() error {
	updates, err := s.bot.getUpdates()
	if err != nil {
		return err
	}

	statuses := make(map[model.Subscriber]model.Status)

	for _, u := range updates {
		subs, st, isOk := mapToSubscriberWithStatus(u, s.bot.Id)
		if isOk {
			statuses[subs] = st
		}
	}

	subsByStatus := make(map[model.Status][]model.Subscriber)

	for u, s := range statuses {
		if subsByStatus[s] == nil {
			subsByStatus[s] = make([]model.Subscriber, 0)
		}
		subsByStatus[s] = append(subsByStatus[s], u)
	}

	fmt.Println("Пользователи по статусам: ", subsByStatus)

	return nil

}

func mapToSubscriberWithStatus(u tgmodel.Update, botId int) (model.Subscriber, model.Status, bool) {
	var status model.Status = math.MinInt32
	switch commands.Parse(u.Message.Text) {
	case commands.START:
		status = model.NEW
		break
	case commands.SUBSCRIBE:
		status = model.ACTIVE
		break
	case commands.UNSUBSCRIBE:
		status = model.INACTIVE
		break
	case commands.UNDEFINED:
		if isBotKicked(u, botId) {
			return model.Subscriber{
				ID:        int64(u.MyChatMember.Chat.Id),
				Type:      model.Telegram,
				Address:   "@" + u.MyChatMember.Chat.Username,
				CreatedAt: nil,
				Status:    model.INACTIVE,
			}, model.INACTIVE, true
		}
		//fmt.Println("Не удалось определить команду: ", u)
	}
	if status == math.MinInt32 {
		return model.Subscriber{}, model.INACTIVE, false
	}
	return mapToSubscriber(u), status, true
}

func mapToSubscriber(u tgmodel.Update) model.Subscriber {
	return model.Subscriber{
		ID:        int64(u.Message.Chat.Id),
		Type:      model.Telegram,
		Address:   "@" + u.Message.From.Username,
		CreatedAt: nil,
		Status:    model.INACTIVE,
	}
}

func isBotKicked(u tgmodel.Update, botId int) bool {
	if u.MyChatMember.NewChatMember.User.Id == botId {
		return u.MyChatMember.NewChatMember.Status == "kicked"
	}
	return false
}

package telegram

import (
	"context"
	"fmt"
	_ "go/parser"
	"horoscope/internal/horoscope"
	"horoscope/internal/horoscope/model"
	"horoscope/internal/telegram/commands"
	tgmodel "horoscope/internal/telegram/model"
	"math"
)

type Service struct {
	bot    *bot
	dao    horoscope.SubscribersDao
	offset int64
}

func NewService(token string, dao horoscope.SubscribersDao) *Service {
	return &Service{newBot(token), dao, 0}
}

func (s *Service) ProcessUpdates() error {
	updates, err := s.bot.getUpdates()
	if err != nil {
		return err
	}

	statuses := make(map[model.Subscriber]model.Status)

	for _, u := range updates {
		s.offset = u.UpdateId
		subs, st, isOk := mapToSubscriberWithStatus(u, s.bot.Id)
		if isOk {
			statuses[subs] = st
		}
	}

	subsByStatus := make(map[model.Status][]model.Subscriber)

	for u, s := range statuses {
		u.Status = s
		if subsByStatus[s] == nil {
			subsByStatus[s] = make([]model.Subscriber, 0)
		}
		subsByStatus[s] = append(subsByStatus[s], u)
	}

	for k, v := range subsByStatus {
		if len(v) == 0 {
			continue
		}
		fmt.Printf("Обновляем %d пользователей на статус %d\n", len(v), k)

		c, e := s.dao.BatchUpdateAll(context.Background(), v)
		if e != nil {
			fmt.Println("Не удалось провести обновление, %w", e)
		} else {
			fmt.Printf("Обновлено пользователей: %d\n", c)
		}
		if k == model.ACTIVE {
			//отправить приветствие
		}
	}

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
			return mapToSubscriberFromChat(u.MyChatMember), model.BANNED, true
		}
	}
	if status == math.MinInt32 {
		return model.Subscriber{}, model.INACTIVE, false
	}
	return mapToSubscriber(u.Message), status, true
}

func isBotKicked(u tgmodel.Update, botId int) bool {
	if u.MyChatMember.NewChatMember.User.Id == botId {
		return u.MyChatMember.NewChatMember.Status == "kicked"
	}
	return false
}

func mapToSubscriber(updateMsg tgmodel.Message) model.Subscriber {
	return model.Subscriber{
		ID:        int64(updateMsg.Chat.Id),
		Type:      model.Telegram,
		Address:   "@" + updateMsg.From.Username,
		CreatedAt: nil,
		Status:    model.INACTIVE,
	}
}

func mapToSubscriberFromChat(chatMemberUpdated tgmodel.ChatMemberUpdated) model.Subscriber {
	return model.Subscriber{
		ID:        int64(chatMemberUpdated.Chat.Id),
		Type:      model.Telegram,
		Address:   "@" + chatMemberUpdated.From.Username,
		CreatedAt: nil,
		Status:    model.INACTIVE,
	}
}

package commands

import "horoscope/internal/telegram/model"

type Type string

const (
	START       Type = "/start"
	SUBSCRIBE   Type = "/subscribe"
	UNSUBSCRIBE Type = "/unsubscribe"
	UNDEFINED   Type = ""
)

type Command struct {
	Type Type
	User model.User
}

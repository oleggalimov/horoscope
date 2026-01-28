package model

import "time"

type Channel string

const (
	Telegram Channel = "telegram"
	Email    Channel = "email"
	Push     Channel = "push"
)

type Subscriber struct {
	ID        int64
	Type      Channel
	Address   string
	Sign      ZodiacSign
	CreatedAt time.Time
}

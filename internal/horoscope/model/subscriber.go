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
	ChannelId string
	Sign      ZodiacSign
	CreatedAt *time.Time
	Status    Status
}

type Status int

const (
	NEW      Status = -1
	ACTIVE          = 1
	INACTIVE        = 0
	BANNED          = -2
)

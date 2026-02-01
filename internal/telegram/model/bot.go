package model

type UpdatesResponse struct {
	IsOk   bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	UpdateId     int               `json:"update_id"`
	MyChatMember ChatMemberUpdated `json:"my_chat_member"`
	Message      Message           `json:"message"`
}

type ChatMemberUpdated struct {
	Chat          Chat       `json:"chat"`
	From          User       `json:"from"`
	Date          int        `json:"date"`
	OldChatMember ChatMember `json:"old_chat_member"`
	NewChatMember ChatMember `json:"new_chat_member"`
}

type Chat struct {
	Id        int    `json:"id"`
	ChatType  string `json:"type"`
	Title     string `json:"title"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type User struct {
	Id           int    `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type ChatMember struct {
	user   User
	status string
}

type Message struct {
	MessageId int64           `json:"message_id"`
	From      User            `json:"from"`
	Chat      Chat            `json:"chat"`
	Date      int64           `json:"date"`
	Text      string          `json:"text"`
	Entities  []MessageEntity `json:"entities"`
}

type MessageEntity struct {
	MessageType string `json:"type"`
	Offset      int    `json:"offset"`
	Length      int    `json:"length"`
	Url         string `json:"url"`
	User        User   `json:"user"`
	Language    string `json:"language"`
}

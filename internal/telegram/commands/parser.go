package commands

func Parse(text string) Type {
	switch text {
	case "/start":
		return START
	case "/subscribe":
		return SUBSCRIBE
	case "/unsubscribe":
		return UNSUBSCRIBE
	default:
		return UNDEFINED
	}
}

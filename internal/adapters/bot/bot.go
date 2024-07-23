package bot

type IBot interface {
	Start()
	SendSlices(to *MessageRecipient, slides []StorySlide) error
	SendMsg(to *MessageRecipient, msg string) error
}

type ITelegramBotRegistry interface {
	HandleCommand(command string, fn func(data *TelegramCommandData)) error
}

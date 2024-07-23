package bot

type IBot interface {
	SendStory(to *MessageRecipient, slides []StorySlide) error
	SendMsg(to *MessageRecipient, msg string) (mid int, err error)
	EditMsg(to *MessageRecipient, msg string) (mid int, err error)
	Start()
	Stop()
}

type IBotRegistry interface {
	HandleCommand(command string, fn func(data *CommandData) error)
}

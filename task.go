package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (task UpdateTask) execute(chat *Chat) (result TaskChan) {
	screen, new, chatTable := chat.Screen.Handle(task.update)
	if chatTable != nil {
		_, fail := BotAPI.Request(chatTable)
		if fail != nil {
			panic("UpdateTask.execute-1: " + fail.Error())
		}
	}
	if screen != nil {
		result = screen.Init()
		chat.Screen.Close()
		chat.out(screen, new)
	}
	return
}

type (
	TaskInterface interface {
		execute(chat *Chat) TaskChan
	}

	UpdateTask struct {
		update tgbotapi.Update
	}
)

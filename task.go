package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (task UpdateTask) execute(chat *Chat) (result TaskChan) {
	screen, new, chattable := chat.Screen.Handle(task.update)
	if chattable != nil {
		_, fail := BotAPI.Request(chattable)
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

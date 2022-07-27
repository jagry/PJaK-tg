package main

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (task UpdateTask) execute(chat *Chat) (result TaskChan) {
	screen, new, chatTable := chat.Screen.Handle(task.update)
	if chatTable != nil {
		_, fail := BotAPI.Request(chatTable)
		if fail != nil {
			var error *tgbotapi.Error
			if errors.As(fail, &error) {
				panic("UpdateTask.execute-0: " + fail.Error())
			} else {
				panic("UpdateTask.execute-1: " + fail.Error())
			}
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

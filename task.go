package main

import (
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func (task UpdateTask) execute(chat *Chat) (result TaskChan) {
	screen, reboot, chatTable := chat.Screen.Handle(chat.Message, task.update)
	if chatTable != nil {
		_, fail := BotAPI.Request(chatTable)
		if fail != nil {
			var tgError *tgbotapi.Error
			if errors.As(fail, &tgError) {
				log.Println("UpdateTask.execute-0: " + fail.Error())
			} else {
				log.Println("UpdateTask.execute-1: " + fail.Error())
			}
		}
	}
	if screen != nil {
		result = screen.Init()
		chat.Screen.Close()
		chat.out(screen, reboot)
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

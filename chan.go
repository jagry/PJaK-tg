package main

import "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type (
	TaskChan chan bool

	UpdateChan chan tgbotapi.Update
)

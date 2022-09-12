package screens

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type (
	EventInterface interface {
		Execute()
	}

	Interface interface {
		GetBase() Base
		Close()
		Init() chan bool
		Handle(int, telegram.Update) (Interface, bool, telegram.Chattable)
		Hook(Interface) Interface
		Out() *InterfaceOut
	}

	InterfaceOut struct {
		Keyboard [][]telegram.InlineKeyboardButton
		Text     string
	}
)

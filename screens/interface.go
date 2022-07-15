package screens

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type (
	EventInterface interface {
		Execute()
	}

	Interface interface {
		Channel() chan Event
		Close()
		Init() chan bool
		Handle(telegram.Update) (Interface, bool, telegram.Chattable)
		Hook(Interface) Interface
		Out() *InterfaceOut
		User() int8
	}

	InterfaceOut struct {
		Keyboard [][]telegram.InlineKeyboardButton
		Text     string
	}
)

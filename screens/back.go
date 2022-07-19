package screens

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	backText = "⭕ Назад"
	backId   = "back"
)

func NewBack(section Section, caller, loader Interface) Back {
	return Back{Section: section, buttons: backRow, caller: caller, loader: loader}
}

func (back Back) Handle(update telegram.Update) (Interface, bool, telegram.Chattable) {
	if update.CallbackQuery != nil && update.CallbackQuery.Data == backId {
		return back.caller, false, telegram.NewCallback(update.CallbackQuery.ID, backText)
	}
	return back.Base.Handle(update)
}

type (
	Back struct {
		Section
		buttons        []telegram.InlineKeyboardButton
		caller, loader Interface
	}
)

var (
	backButton   = telegram.NewInlineKeyboardButtonData(backText, backId)
	backKeyboard = [][]telegram.InlineKeyboardButton{backRow}
	backRow      = []telegram.InlineKeyboardButton{backButton, baseCloseButton}
)

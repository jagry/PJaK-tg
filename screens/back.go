package screens

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	backText = "⭕ Назад"
	backId   = "back"
)

func NewBack(base Base, caller, loader Interface) Back {
	return Back{Base: base, buttons: backRow, caller: caller, loader: loader}
}

func (back Back) Handle(id int, update telegram.Update) (Interface, bool, telegram.Chattable) {
	if update.CallbackQuery != nil && update.CallbackQuery.Data == backId {
		if update.CallbackQuery.Message != nil && update.CallbackQuery.Message.MessageID == id {
			return back.caller, false, telegram.NewCallback(update.CallbackQuery.ID, backText)
		}
	}
	return back.Base.Handle(id, update)
}

type Back struct {
	Base
	buttons        []telegram.InlineKeyboardButton
	caller, loader Interface
}

var (
	backButton   = telegram.NewInlineKeyboardButtonData(backText, backId)
	backKeyboard = [][]telegram.InlineKeyboardButton{backRow}
	backRow      = []telegram.InlineKeyboardButton{backButton, baseCloseButton}
)

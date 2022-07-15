package screens

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
//errorId   = "error"
//errorText = closeTextConst
)

func NewError(base Base, caller Interface, caption, text string) Error {
	return Error{Base: base, caller: caller, view: NewView(caption, text)}
}

func (error Error) Handle(update telegram.Update) (Interface, bool, telegram.Chattable) {
	if update.CallbackQuery != nil && update.CallbackQuery.Data == "error" {
		return error.caller, false, nil
	}
	return error.Base.Handle(update)
}

func (error Error) Hook(argument Interface) Interface { return nil }

func (error Error) Out() *InterfaceOut {
	return &InterfaceOut{
		Keyboard: [][]telegram.InlineKeyboardButton{errorKeyRow},
		Text:     error.view.Text()}
}

type Error struct {
	Base
	caller Interface
	view   View
}

var errorKeyRow = []telegram.InlineKeyboardButton{telegram.NewInlineKeyboardButtonData(baseCloseText, baseCloseId)}

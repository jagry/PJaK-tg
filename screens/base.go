package screens

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	baseCaption          = "Ошибка"
	baseUserCaption      = baseUserCaptionEmoji + " " + baseUserCaptionText
	baseUserCaptionEmoji = "🤦"
	baseUserCaptionText  = "Ошибка пользователя"
	baseCloseId          = "base.close"

	baseText        = "Что-то пошло не так"
	baseCloseText   = "❌ Закрыть"
	baseCommandText = "Не известная команда"
	baseUpdateText  = "Не надо бездумно слать всякую шнягу"
)

func NewBase(user int8) Base {
	return Base{user: user}
}

func (Base) Channel() chan Event { return nil }

func (Base) Close() {}

func (Base) Init() chan bool { return nil }

func (base Base) Handle(update telegram.Update) (Interface, bool, telegram.Chattable) {
	if update.CallbackQuery != nil {
		if update.CallbackQuery.Data == baseCloseId {
			return base, false, telegram.NewCallback(update.CallbackQuery.ID, baseCloseText)
		}
	}
	if update.Message != nil {
		dmc := telegram.NewDeleteMessage(update.FromChat().ID, update.Message.MessageID)
		if update.Message.IsCommand() {
			switch update.Message.Text {
			case "/start":
				return NewLoading(base, base, NewBetsMainFactory(), betsCaption, betsLoadTournamentsText), true, dmc
			case "/tournaments":
				return NewLoading(base, base, NewBetsMainFactory(), betsCaption, betsLoadTournamentsText), true, dmc
			}
			return NewError(base, base, baseUserCaption, baseCommandText), true, dmc
		}
		return NewError(base, base, baseUserCaption, baseUpdateText), true, dmc
	}
	return NewError(base, base, baseCaption, baseText), true, nil
}

func (Base) Hook(argument Interface) Interface { return argument }

func (Base) Out() *InterfaceOut {
	return nil
}

func (base Base) User() int8 {
	return base.user
}

type Base struct {
	//chat int64
	user int8
}

var baseCloseButton = telegram.NewInlineKeyboardButtonData(baseCloseText, baseCloseId)

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

func NewBase(channel chan Event, user int8) Base {
	return Base{channel: channel, user: user}
}

func NewBaseAction(section string) BaseAction {
	return BaseAction{section: section}
}

//func (Base) Channel() chan Event { return nil }

func (Base) Close() {}

func (Base) Message() {}

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
			case "/bets":
				return NewLoading(base, base, NewMainFactory(betsMainManager, betsCaption), mainLoadText), true, dmc
			case "/results":
				factory := NewMainFactory(resultsMainManager, resultsCaption)
				return NewLoading(base, base, factory, mainLoadText), true, dmc
			case "/start":
				//return NewLoading(base, base, NewMainFactory(bets()), betsCaption, mainLoadText), true, dmc
			case "/statistic":
				return NewLoading(base, base, NewMainFactory(statisticMainManager, statisticCaption), mainLoadText), true, dmc
			}
			return NewError(base, base, baseUserCaption, baseCommandText), true, dmc
		}
		return NewError(base, base, baseUserCaption, baseUpdateText), true, dmc
	}
	return NewError(base, base, baseCaption, baseText), true, nil
}

func (Base) Hook(argument Interface) Interface { return argument }

func (Base) Out() *InterfaceOut { return nil }

//func (base Base) User() int8 { return base.user }

//func (ba BaseAction) Section() string { return ba.section }

type (
	Base struct {
		channel chan Event
		user    int8
	}

	BaseAction struct {
		section string
	}
)

var baseCloseButton = telegram.NewInlineKeyboardButtonData(baseCloseText, baseCloseId)

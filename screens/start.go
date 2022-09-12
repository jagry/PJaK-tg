package screens

//import (
//	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
//)
//
//func NewStart(message int) Start {
//	return Start{}
//}
//
//func StartHandle(update telegram.Update) (Interface, bool, telegram.Chattable) {
//	if update.CallbackQuery != nil {
//		if update.Message.IsCommand() {
//			if update.Message.Text == "/start" {
//				return NewStart(update.Message.MessageID), false, nil
//			}
//		}
//	}
//	return NewStart(update.Message.MessageID), false, nil
//}
//
//func (start Start) Handle(update telegram.Update) (Interface, bool, telegram.Chattable) {
//	if update.CallbackQuery == nil {
//		return start, false, nil
//	}
//	return StartHandle(update)
//}
//
//func (Start) Out() *InterfaceOut {
//	//config := telegram.NewMessage(id, "Pfuheprf")
//	//config.ReplyMarkup = startKeyboard
//	//message, fail := bot.Send(config)
//	//if fail != nil {
//	//	p_anic("screens.Start.Out" + fail.Error())
//	//}
//	//return message.MessageID
//	return nil
//}
//
//func (Start) Save() {
//
//}
//
//var startKeyboard = telegram.NewInlineKeyboardMarkup(telegram.NewInlineKeyboardRow(telegram.NewInlineKeyboardButtonData(cancelTextConst, cancelIdConst)))
//
//type Start struct {
//	Base
//}

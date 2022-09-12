package screens

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"math/rand"
	"time"
)

const (
	actionButtonId    = "action"
	actionButtonText  = "🛑 Отмена"
	actionCharsCount  = 3
	loadingTextSuffix = ". Подождите, пожалуйста"
)

func NewLoading(b Base, i Interface, e ActionExecutor, t string) *Loading {
	runes := []rune(actionChars[rand.Intn(actionCharsCount)])
	return &Loading{Base: b, caller: i, executor: e, index: -1, runes: runes, view: NewView(e.Caption(), t)}
}

func (loading Loading) Channel() chan Event {
	return loading.channel
}

func (loading Loading) Close() {
}

func (loading Loading) close() {
	//loading.channel <- nil
}

func (loading *Loading) execute(event chan Event) {
	//time.Sleep(time.Second * 9)
	log.Println("screens.Loading.execute: sending event")
	event <- loading.executor.Do(loading)
	log.Println("screens.Loading.execute: sent event")
}

func (loading Loading) Handle(id int, update telegram.Update) (Interface, bool, telegram.Chattable) {
	if update.CallbackQuery != nil && update.CallbackQuery.Data == actionButtonId {
		if update.CallbackQuery.Message != nil && update.CallbackQuery.Message.MessageID == id {
			return loading.caller, false, telegram.NewCallback(update.CallbackQuery.ID, actionButtonText)
		}
	}
	return loading.Base.Handle(id, update)
}

func (loading *Loading) Init() chan bool {
	log.Println("screens.Loading.Init: creating event")
	event := make(chan Event)
	log.Println("screens.Loading.Init: created event")
	go loading.execute(event)
	go loading.timer(event)
	return nil
}

func (loading Loading) Out() *InterfaceOut {
	view := loading.view
	if loading.index >= 0 {
		view.text = string(loading.runes[loading.index]) + " " + view.text
	}
	return &InterfaceOut{Keyboard: [][]telegram.InlineKeyboardButton{{actionButton}}, Text: view.Text()}
}

func (loading *Loading) timer(event chan Event) {
	for {
		log.Println("screens.Loading.timer: select event and timer")
		select {
		case finish := <-event:
			log.Println("screens.Loading.timer: received event")
			if finish.Screen != nil {
				log.Println("screens.Loading.timer: sending to chat 0")
				loading.channel <- finish
				log.Println("screens.Loading.timer: sent to chat 0")
				log.Println("screens.Loading.timer: closing event 0")
				close(event)
				log.Println("screens.Loading.timer: closed event 0")
				return
			}
			//for {
			log.Println("screens.Loading.timer: wait event")
			finish = <-event
			log.Println("screens.Loading.timer: received event")
			if finish.Screen != nil {
				log.Println("screens.Loading.timer: closing event 1")
				close(event)
				log.Println("screens.Loading.timer: closed event 1")
				return
			}
			//}
		case <-time.After(time.Second):
			log.Println("screens.Loading.timer: received timer")
			loading.index++
			if loading.index == len(loading.runes) {
				loading.index = 0
			}
			log.Println("screens.Loading.timer: sending to chat 1", loading.channel)
			loading.channel <- Event{}
			log.Println("screens.Loading.timer: sent to chat 1")
		}
	}
}

type (
	Loading struct {
		Base
		caller   Interface
		executor ActionExecutor
		index    int
		runes    []rune
		view     View
	}

	LoadingChannel chan Interface

	LoadingCharacters struct {
		data    string
		shuffle bool
	}

	ActionExecutor interface {
		Caption() string
		Do(action *Loading) Event
	}
)

var (
	actionChars = [actionCharsCount]string{
		"⌛⏳",
		"⚫🔴\U0001F7E0\U0001F7E1\U0001F7E4\U0001F7E2🔵\U0001F7E3⚪",
		"🕛🕧🕐🕜🕑🕝🕒🕞🕓🕟🕔🕠🕕🕡🕖🕢🕗🕣🕘🕤🕙🕥🕚🕦"}
	actionButton = telegram.NewInlineKeyboardButtonData(actionButtonText, actionButtonId)
)

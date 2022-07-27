package screens

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"math/rand"
	"time"
)

const (
	loadingCharactersCount = 3
	loadingButton          = "🛑 Отмена"
	loadingTextSuffix      = ". Подождите, пожалуйста"
)

func NewLoading(b Base, i Interface, f ActionFactory, t string) *Loading {
	runes, view := []rune(loadingCharacters[rand.Intn(loadingCharactersCount)]), NewView(f.Caption(), t)
	return &Loading{Base: b, caller: i, factory: f, index: -1, runes: runes, view: view}
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
	//time.Sleep(time.Second)
	log.Println("screens.Loading.execute: sending event")
	event <- loading.factory.Execute(loading)
	log.Println("screens.Loading.execute: sent event")
}

func (loading Loading) Handle(update telegram.Update) (Interface, bool, telegram.Chattable) {
	if update.CallbackQuery != nil && update.CallbackQuery.Data == cancelIdConst {
		return loading.caller, false, telegram.NewCallback(update.CallbackQuery.ID, cancelTextConst)
	}
	return loading.Base.Handle(update)
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
	return &InterfaceOut{
		Keyboard: [][]telegram.InlineKeyboardButton{{actionKeyboardButton}},
		Text:     view.Text()}
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
		caller  Interface
		factory ActionFactory
		index   int
		runes   []rune
		view    View
	}

	LoadingChannel chan Interface

	LoadingCharacters struct {
		data    string
		shuffle bool
	}

	ActionFactory interface {
		Caption() string
		Execute(action *Loading) Event
	}
)

var (
	loadingCharacters = [loadingCharactersCount]string{
		"⌛⏳",
		"⚫🔴\U0001F7E0\U0001F7E1\U0001F7E4\U0001F7E2🔵\U0001F7E3⚪",
		"🕛🕧🕐🕜🕑🕝🕒🕞🕓🕟🕔🕠🕕🕡🕖🕢🕗🕣🕘🕤🕙🕥🕚🕦"}
	actionKeyboardButton = telegram.NewInlineKeyboardButtonData(loadingButton, cancelIdConst)
)

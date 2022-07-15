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
	loadingText            = ". Подождите, пожалуйста"
)

func NewLoading(base Base, caller Interface, factory ActionFactory, caption, text string) *Loading {
	return &Loading{
		Base:      base,
		caller:    caller,
		channel:   make(LoadingChannel),
		eventChan: make(chan Event),
		factory:   factory,
		index:     -1,
		runes:     []rune(loadingCharacters[rand.Intn(loadingCharactersCount)]),
		view:      NewView(caption, text)}
}

func (loading Loading) Channel() chan Event {
	return loading.eventChan
}

func (loading Loading) Close() {
	go loading.close()
}

func (loading Loading) close() {
	log.Println("screens.Loading.Close: start")
	loading.channel <- nil
	log.Println("screens.Loading.Close: finish")
}
func (loading *Loading) execute() {
	time.Sleep(time.Second)
	loading.channel <- loading.factory.Execute(loading.Base)
}

func (loading Loading) Handle(update telegram.Update) (Interface, bool, telegram.Chattable) {
	if update.CallbackQuery != nil && update.CallbackQuery.Data == cancelIdConst {
		return loading.caller, false, telegram.NewCallback(update.CallbackQuery.ID, cancelTextConst)
	}
	return loading.Base.Handle(update)
}

func (loading *Loading) Init() chan bool {
	//wg := sync.WaitGroup{}
	go loading.execute()
	go loading.timer()
	//wg.Wait()
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

func (loading *Loading) timer() {
	log.Println("screens.Loading.timer 0: start")
	for {
		select {
		case finish := <-loading.channel:
			log.Println("screens.Loading.timer: channel0 =", finish)
			if finish != nil {
				log.Println("screens.Loading.timer: channel0: start")
				loading.eventChan <- finish
				log.Println("screens.Loading.timer: channel0: notify")
				close(loading.eventChan)
				log.Println("screens.Loading.timer: channel0: close")
				loading.eventChan = nil
				log.Println("screens.Loading.timer: channel0: finish")
				return
			}
			for {
				finish = <-loading.channel
				log.Println("screens.Loading.timer: channel1")
				if finish != nil {
					log.Println("screens.Loading.timer: channel1: finish")
					close(loading.channel)
					return
				}
			}
		case <-time.After(time.Second):
			log.Println("screens.Loading.timer: timer")
			loading.index++
			if loading.index == len(loading.runes) {
				loading.index = 0
			}
			log.Println("screens.Loading.timer: timer: pre")
			loading.eventChan <- nil
			log.Println("screens.Loading.timer: timer: post")
		}
	}
}

type (
	Loading struct {
		Base
		caller    Interface
		channel   LoadingChannel
		eventChan chan Event
		factory   ActionFactory
		index     int
		runes     []rune
		view      View
	}

	LoadingChannel chan Interface

	LoadingCharacters struct {
		data    string
		shuffle bool
	}

	ActionFactory interface{ Execute(base Base) Interface }
)

var (
	loadingCharacters = [loadingCharactersCount]string{
		"⌛⏳",
		"⚫🔴\U0001F7E0\U0001F7E1\U0001F7E4\U0001F7E2🔵\U0001F7E3⚪",
		"🕛🕧🕐🕜🕑🕝🕒🕞🕓🕟🕔🕠🕕🕡🕖🕢🕗🕣🕘🕤🕙🕥🕚🕦"}
	actionKeyboardButton = telegram.NewInlineKeyboardButtonData(loadingButton, cancelIdConst)
)

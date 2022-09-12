package main

import (
	"PJaK/screens"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"sync"
	"time"
)

func (chat *Chat) end() {
	//EndChatChannel <- chat.Screen.Chat()
}

func (chat *Chat) out(screen screens.Interface, new bool) {
	if screen != nil {
		chat.Screen = screen
	}
	out := chat.Screen.Out()
	if new {
		if chat.Message != 0 {
			if _, fail := BotAPI.Request(tgbotapi.NewDeleteMessage(chat.id, chat.Message)); fail != nil {
				log.Println("Chat.out-0: " + fail.Error())
			}
			chat.Message = 0
		}
		if out != nil {
			config := tgbotapi.NewMessage(chat.id, out.Text)
			config.ParseMode = "html"
			if len(out.Keyboard) > 0 {
				config.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(out.Keyboard...)
			}
			message, fail := BotAPI.Send(config)
			if fail != nil {
				log.Println("Chat.out-1: " + fail.Error())
			}
			chat.Message = message.MessageID
		}
	} else {
		if chat.Message == 0 {
			if out != nil {
				config := tgbotapi.NewMessage(chat.id, out.Text)
				config.ParseMode = "html"
				if len(out.Keyboard) > 0 {
					config.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(out.Keyboard...)
				}
				message, fail := BotAPI.Send(config)
				if fail != nil {
					log.Println("Chat.out-2: " + fail.Error())
				}
				chat.Message = message.MessageID
			}
		} else {
			if out == nil {
				if _, fail := BotAPI.Request(tgbotapi.NewDeleteMessage(chat.id, chat.Message)); fail != nil {
					log.Println("Chat.out-3: " + fail.Error())
				}
				chat.Message = 0
			} else {
				config := tgbotapi.NewEditMessageText(chat.id, chat.Message, out.Text)
				if len(out.Keyboard) > 0 {
					keyboard := tgbotapi.NewInlineKeyboardMarkup(out.Keyboard...)
					config.ReplyMarkup = &keyboard
				}
				config.ParseMode = "html"
				_, fail := BotAPI.Request(config)
				if fail != nil {
					log.Println("Chat.out-4: " + fail.Error())
				}
			}
		}
	}
}

func (chat *Chat) routine(channel chan screens.Event) {
	for {
		select {
		case update := <-chat.updateChan:
			if len(chat.taskSlice) == 0 {
				chat.taskChan = UpdateTask{update}.execute(chat)
			} else {
				chat.taskSlice = append(chat.taskSlice, UpdateTask{update})
			}
		case event := <-channel:
			if event.Screen != nil {
				event.Screen.Init()
				chat.Screen.Close()
			}
			if event.Text != "" {
				BotAPI.Request(tgbotapi.NewDeleteMessage(chat.id, chat.Message))
				message := tgbotapi.NewMessage(chat.id, event.Text)
				message.ParseMode = "html"
				BotAPI.Request(message)
				chat.Message = 0
			}
			chat.out(event.Screen, false)
		case <-time.After(time.Minute):
			chat.Screen.Close()
			chat.out(chat.Screen.GetBase(), false)
		}
	}
}

func (chat *Chat) update(screen screens.Interface) {
	if chat.Next != nil {
		chat.Next, chat.Next.Previous = nil, chat.Previous
	}
	if chat.Previous != nil {
		chat.Previous.Next = chat.Next
	}
	chat.Time, LastChat = time.Now(), chat
	chat.out(screen, false)
}

type (
	Chat struct {
		Next, Previous *Chat
		Message        int
		user           int8
		id             int64
		Screen         screens.Interface
		taskChan       TaskChan
		taskSlice      []TaskInterface
		Time           time.Time
		updateChan     UpdateChan
	}
)

var (
	ChatMap      = map[int64]*Chat{}
	ChatMutex    sync.Mutex
	FirstChat    *Chat
	LastChat     *Chat
	EventChannel = make(chan screens.Event)
)

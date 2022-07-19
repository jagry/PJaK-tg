package main

import (
	"PJaK/screens"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
				panic("Chat.out-0: " + fail.Error())
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
				panic("Chat.out-0: " + fail.Error())
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
					panic("Chat.out-0: " + fail.Error())
				}
				chat.Message = message.MessageID
			}
		} else {
			if out == nil {
				if _, fail := BotAPI.Request(tgbotapi.NewDeleteMessage(chat.id, chat.Message)); fail != nil {
					panic("Chat.out-1: " + fail.Error())
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
					//fail = errors.Unwrap(fail)
					panic("Chat.out-2: " + fail.Error())
				}
			}
		}
	}
}

func (chat *Chat) routine(channel chan screens.Event) {
	for {
		//if channel := chat.Screen.Channel(); channel == nil {
		//	update := <-chat.updateChan
		//	if len(chat.taskSlice) == 0 {
		//		chat.taskChan = UpdateTask{update}.execute(chat)
		//	} else {
		//		chat.taskSlice = append(chat.taskSlice, UpdateTask{update})
		//	}
		//} else {
		select {
		case update := <-chat.updateChan:
			if len(chat.taskSlice) == 0 {
				chat.taskChan = UpdateTask{update}.execute(chat)
			} else {
				chat.taskSlice = append(chat.taskSlice, UpdateTask{update})
			}
		case event := <-channel:
			if event != nil {
				event.Init()
				chat.Screen.Close()
			}
			chat.out(event, false)
		}
		//		}
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

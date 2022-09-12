package main

import (
	"PJaK/core"
	"PJaK/screens"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/go-yaml/yaml"
	_ "github.com/lib/pq"
	"math/rand"
	"os"
	"time"
)

var BotAPI *tgbotapi.BotAPI

func main() {
	var settings Settings
	bytes, fail := os.ReadFile("./.yml")
	if fail != nil {
		fmt.Println(fail.Error())
		os.Exit(1)
		return
	}
	fail = yaml.Unmarshal([]byte(bytes), &settings)
	if fail != nil {
		fmt.Println(fail.Error())
		os.Exit(2)
		return
	}
	fail = core.Db(settings.Database.Type, settings.Database.Arguments)
	if fail != nil {
		fmt.Println(fail.Error())
		os.Exit(3)
		return
	}
	BotAPI, fail = tgbotapi.NewBotAPI(settings.Telegram)
	if fail != nil {
		fmt.Println(fail.Error())
		os.Exit(4)
		return
	}
	rand.Seed(time.Now().UnixNano())
	update := tgbotapi.NewUpdate(0)
	update.Timeout = 9
	channel := BotAPI.GetUpdatesChan(update)
	for {
		select {
		case message := <-channel:
			messageChat := message.FromChat()
			if messageChat != nil {
				chat, ok := ChatMap[messageChat.ID]
				if ok {
					/*if chat.Next != nil {
						chat.Next.Previous = chat.Previous
					}
					if chat.Previous != nil {
						chat.Previous.Next = chat.Next
					} else {
						FirstChat = chat.Next
					}
					chat.Time = time.Now()*/
				} else {
					userName := message.SentFrom().UserName
					userId := int8(-1)
					if "Whiteseaer" == userName {
						userId = 1
					} else if "petrperke" == userName {
						userId = 0
					} else if "El_Kardo" == userName {
						userId = 2
					}
					chatChannel := make(chan screens.Event)
					screen := screens.NewBase(chatChannel, userId)
					chat = &Chat{id: messageChat.ID, Screen: screen, Time: time.Now(), updateChan: make(UpdateChan)}
					ChatMap[messageChat.ID] = chat
					go chat.routine(chatChannel)
				}
				/*chat.Time = time.Now()*/
				chat.updateChan <- message
				/*if LastChat != nil {
					LastChat.Next = chat
				} else {
					FirstChat = chat
				}
				LastChat, chat.Next, chat.Previous = chat, nil, LastChat*/
			}
		}
	}
}

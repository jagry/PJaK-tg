package main

import (
	"PJaK/core"
	"PJaK/screens"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/go-yaml/yaml"
	_ "github.com/lib/pq"
	"log"
	"math/rand"
	"os"
	"time"
)

var BotAPI *tgbotapi.BotAPI

func main() {
	var settings Settings
	bytes, fail := os.ReadFile("./.yml")
	if fail != nil {
		log.Panic(fail.Error())
	}
	fail = yaml.Unmarshal([]byte(bytes), &settings)
	if fail != nil {
		log.Panic(fail.Error())
	}
	log.Println(settings)
	fail = core.Db(settings.Database.Type, settings.Database.Arguments)
	if fail != nil {
		log.Panic(fail.Error())
	}
	BotAPI, fail = tgbotapi.NewBotAPI("5538271712:AAH-1G6d5LEQXsDnGZY7qTwS9hNVu_0XPoQ")
	if fail != nil {
		log.Panic(fail)
	}
	rand.Seed(time.Now().UnixNano())
	update := tgbotapi.NewUpdate(0)
	update.Timeout = 9
	channel := BotAPI.GetUpdatesChan(update)
	for {
		select {
		case message := <-channel:
			log.Println("main.Update")
			messageChat := message.FromChat()

			log.Println("main.Update", messageChat)
			if messageChat != nil {
				chat, ok := ChatMap[message.FromChat().ID]
				if ok {
					if chat.Next != nil {
						chat.Next.Previous = chat.Previous
					}
					if chat.Previous != nil {
						chat.Previous.Next = chat.Next
					} else {
						FirstChat = chat.Next
					}
					chat.Time = time.Now()
				} else {
					userName := message.SentFrom().UserName
					userId := int8(-1)
					if "Whiteseaer" == userName {
						userId = 1
					}
					chat = &Chat{
						id:         message.FromChat().ID,
						Screen:     screens.NewBase(userId),
						Time:       time.Now(),
						updateChan: make(UpdateChan)}
					ChatMap[message.FromChat().ID] = chat
					go chat.routine()
				}
				chat.Time = time.Now()
				chat.updateChan <- message
				if LastChat != nil {
					LastChat.Next = chat
				} else {
					FirstChat = chat
				}
				LastChat, chat.Next, chat.Previous = chat, nil, LastChat
				//ChatMutex.Unlock()
				log.Println("main.Update")
			}
		}
	}
}

package screens

import (
	"PJaK/core"
	"PJaK/views"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
)

const (
	mainLoadText  = "Идет загрузка футбольных турниров" + loadingTextSuffix
	mainErrorText = "Ошибка загрузки футбольных турниров"
)

func loadMain(base Base, manager MainManager, section string) *Loading {
	return NewLoading(base, base, NewLoadMain(manager, section), mainLoadText)
}

func NewMain(base Base, manager MainManager, section string, tournaments []core.Tournament) Main {
	tournamentMap := MapTournament{}
	for offset, tournament := range tournaments {
		tournamentMap[tournament.Id] = offset
	}
	return Main{Base: base, manager: manager, section: section, tournamentMap: tournamentMap, tournamentSlice: tournaments}
}

func NewLoadMain(manager MainManager, section string) LoadMain {
	return LoadMain{manager: manager, section: section}
}

func NewMainManager(tournament MainManagerTournament) MainManager {
	return MainManager{tournament: tournament}
}

func (main Main) Handle(id int, update tgbotapi.Update) (Interface, bool, tgbotapi.Chattable) {
	if update.CallbackQuery != nil && strings.HasPrefix(update.CallbackQuery.Data, betsTournamentIdPrefix) {
		if update.CallbackQuery.Message != nil && update.CallbackQuery.Message.MessageID == id {
			callbackQuery := update.CallbackQuery
			dataString := callbackQuery.Data[len(betsTournamentIdPrefix):]
			dataInt, fail := strconv.ParseUint(dataString, 10, 16)
			if fail == nil {
				offset, ok := main.tournamentMap[core.TournamentId(dataInt)]
				if ok {
					tournament := main.tournamentSlice[offset]
					text := views.Tournament(tournament).Caption(main.section)
					callback := tgbotapi.NewCallback(callbackQuery.ID, text)
					return main.manager.tournament(main, tournament), false, callback
				} else {
					log.Println("screens.Main.Handle-1")
				}
			} else {
				log.Println("screens.Main.Handle-0:" + fail.Error())
			}
		}
	}
	return main.Base.Handle(id, update)
}

func (main Main) Out() *InterfaceOut {
	if len(main.tournamentMap) == 0 {
		return &InterfaceOut{Keyboard: [][]tgbotapi.InlineKeyboardButton{{baseCloseButton}}, Text: betsEmptyText}
	}
	keyboard := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, tournament := range main.tournamentSlice {
		id := betsTournamentIdPrefix + strconv.FormatUint(uint64(tournament.Id), 10)
		button := tgbotapi.NewInlineKeyboardButtonData(views.Tournament(tournament).Name(), id)
		keyboard = append(keyboard, []tgbotapi.InlineKeyboardButton{button})
	}
	keyboard = append(keyboard, []tgbotapi.InlineKeyboardButton{baseCloseButton})
	return &InterfaceOut{Keyboard: keyboard, Text: NewView(main.section, betsTournamentsText).Text()}
}

func (lm LoadMain) Caption() string { return lm.section }

func (lm LoadMain) Do(action *Loading) Event {
	if tournaments, fail := core.GetTournaments(); fail == nil {
		return NewEvent(NewMain(action.Base, lm.manager, lm.section, tournaments), "")
	}
	return NewEvent(NewError(action.Base, action.Base, lm.section, mainErrorText), "")
}

type (
	Main struct {
		Base
		manager         MainManager
		section         string
		tournamentMap   MapTournament
		tournamentSlice []core.Tournament
	}

	LoadMain struct {
		manager MainManager
		section string
	}

	MainManager struct {
		tournament MainManagerTournament
	}

	MainManagerTournament func(main Main, tournament core.Tournament) *Loading
)

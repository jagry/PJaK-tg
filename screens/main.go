package screens

import (
	"PJaK/core"
	"PJaK/views"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

const (
	mainLoadText  = "Идет загрузка футбольных турниров" + loadingTextSuffix
	mainErrorText = "Ошибка загрузки футбольных турниров"
)

func LoadMain(base Base, manager MainManager, section string) *Loading {
	return NewLoading(base, base, NewMainFactory(manager, section), mainLoadText)
}

func NewMain(base Base, manager MainManager, section string, tournaments []core.Tournament) Main {
	tournamentMap := MapTournament{}
	for offset, tournament := range tournaments {
		tournamentMap[tournament.Id] = offset
	}
	return Main{Base: base, manager: manager, section: section, tournamentMap: tournamentMap, tournamentSlice: tournaments}
}

func NewMainFactory(manager MainManager, section string) MainFactory {
	return MainFactory{manager: manager, section: section}
}

func NewMainManager(tournament MainManagerTournament) MainManager {
	return MainManager{tournament: tournament}
}

func (main Main) Handle(update tgbotapi.Update) (Interface, bool, tgbotapi.Chattable) {
	if update.CallbackQuery != nil && strings.HasPrefix(update.CallbackQuery.Data, betsTournamentIdPrefix) {
		id, fail := strconv.ParseUint(update.CallbackQuery.Data[len(betsTournamentIdPrefix):], 10, 16)
		if fail != nil {
			panic("screens.BetsMain.Handle-0:" + fail.Error())
		}
		offset, ok := main.tournamentMap[core.TournamentId(id)]
		if !ok {
			panic("screens.BetsMain.Handle-1")
		}
		tournament := main.tournamentSlice[offset]
		text := views.Tournament(tournament).Caption(main.section)
		return main.manager.tournament(main, tournament), false, tgbotapi.NewCallback(update.CallbackQuery.ID, text)
	}
	return main.Base.Handle(update)
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

func (mf MainFactory) Execute(action *Loading) Interface {
	if tournaments, fail := core.GetTournaments(); fail == nil {
		return NewMain(action.Base, mf.manager, mf.section, tournaments)
	}
	return NewError(action.Base, action.Base, mf.section, mainErrorText)
}

func (mf MainFactory) Caption() string { return mf.section }

type (
	Main struct {
		Base
		manager         MainManager
		section         string
		tournamentMap   MapTournament
		tournamentSlice []core.Tournament
	}

	MainFactory struct {
		manager MainManager
		section string
	}

	MainManager struct {
		tournament MainManagerTournament
	}

	MainManagerTournament func(main Main, tournament core.Tournament) *Loading
)

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

func LoadMain(base Base, section string) *Loading {
	return NewLoading(base, section, base, NewMainFactory(), betsCaption, mainLoadText)
}

func NewMain(base Base, section string, tournaments []*core.Tournament) Main {
	tournamentMap := MapTournament{}
	for _, tournament := range tournaments {
		tournamentMap[tournament.Id] = tournament
	}
	super := NewSection(base, section)
	return Main{Section: super, tournamentMap: tournamentMap, tournamentSlice: tournaments}
}

func NewMainFactory() MainFactory {
	return MainFactory{}
}

func (main Main) Handle(update tgbotapi.Update) (Interface, bool, tgbotapi.Chattable) {
	if update.CallbackQuery != nil && strings.HasPrefix(update.CallbackQuery.Data, betsTournamentIdPrefix) {
		id, fail := strconv.ParseUint(update.CallbackQuery.Data[len(betsTournamentIdPrefix):], 10, 16)
		if fail != nil {
			panic("screens.BetsMain.Handle-0:" + fail.Error())
		}
		tournament, ok := main.tournamentMap[core.TournamentId(id)]
		if !ok {
			panic("screens.BetsMain.Handle-1")
		}
		text := views.Tournament(*tournament).Caption(main.Section.name)
		callback := tgbotapi.NewCallback(update.CallbackQuery.ID, text)
		return LoadTournament(main.Base, main.Section.name, main, tournament), false, callback
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
		button := tgbotapi.NewInlineKeyboardButtonData(views.Tournament(*tournament).Name(), id)
		keyboard = append(keyboard, []tgbotapi.InlineKeyboardButton{button})
	}
	keyboard = append(keyboard, []tgbotapi.InlineKeyboardButton{baseCloseButton})
	return &InterfaceOut{Keyboard: keyboard, Text: NewView(main.Section.name, betsTournamentsText).Text()}
}

func (mf MainFactory) Execute(action *Loading) Interface {
	if tournaments, fail := core.GetTournaments(); fail == nil {
		return NewMain(action.Base, action.Section.name, tournaments)
	}
	return NewError(action.Base, action.Base, action.Section.name, mainErrorText)
}

type (
	Main struct {
		Section
		tournamentMap   MapTournament
		tournamentSlice []*core.Tournament
	}

	MainFactory struct{}
)

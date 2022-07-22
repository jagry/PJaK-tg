package screens

import (
	"PJaK/core"
	"PJaK/views"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

const (
	statisticAllText      = "Общая"
	statisticCaption      = statisticCaptionEmoji + " " + statisticCaptionText
	statisticCaptionEmoji = "🗒"
	statisticCaptionText  = "Статистика"
)

func newLoadStatisticTournament(caller Interface, core core.Tournament) loadStatisticTournament {
	return loadStatisticTournament{caller: caller, core: core}
}

func newStatisticTournament(base Base, caller, loader Interface,
	core core.Tournament, rounds []core.Round) statisticTournament {
	roundMap := initTournament(MapRound{}, rounds)
	return statisticTournament{Back: NewBack(base, caller, loader), core: core, roundMap: roundMap, roundSlice: rounds}
}

func (loadStatisticTournament) Caption() string { return statisticCaption }

func (lst loadStatisticTournament) Execute(action *Loading) Interface {
	if rounds, fail := core.GetTournamentRounds(lst.core); fail == nil {
		return newStatisticTournament(action.Base, lst.caller, action, lst.core, rounds)
	}
	return NewError(action.Base, action, "!!!", "!!!")
}

func (st statisticTournament) Handle(update tgbotapi.Update) (Interface, bool, tgbotapi.Chattable) {
	if update.CallbackQuery != nil {
		if strings.HasPrefix(update.CallbackQuery.Data, betsTournamentIdPrefix) {
			id, fail := strconv.ParseUint(update.CallbackQuery.Data[len(betsTournamentIdPrefix):], 36, 16)
			if fail != nil {
				panic("screens.BetsTournament.Handle:" + fail.Error())
			}
			if round, ok := st.roundMap[core.RoundId(id)]; ok {
				return nil, false, tgbotapi.NewCallback(update.CallbackQuery.ID, round.Name)
			} else {
				return nil, false, tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			}
		}
	}
	return st.Back.Handle(update)
}

func (st statisticTournament) Out() *InterfaceOut {
	if len(st.roundMap) == 0 {
		text := views.Tournament(st.core).Caption(statisticCaption)
		return &InterfaceOut{Keyboard: backKeyboard, Text: NewView(text, betsTournamentEmptyText).Text()}
	}
	keyboard := [][]tgbotapi.InlineKeyboardButton{{tgbotapi.NewInlineKeyboardButtonData(statisticAllText, "ddddd")}}
	rounds := core.GetRounds(st.roundSlice)
	for counter := 0; counter < len(rounds); counter = counter + betsRowCount {
		row := make([]tgbotapi.InlineKeyboardButton, 0)
		bound := counter + betsRowCount
		empty := 0
		if bound > len(rounds) {
			empty = bound - len(rounds)
			bound = len(rounds)
		}
		for offset, item := range rounds[counter:bound] {
			id := betsTournamentIdPrefix + strconv.FormatUint(uint64(item.Id), 36)
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(counter+offset+1), id))
		}
		for empty > 0 {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(" ", betsTournamentIdPrefix+"0"))
			empty--
		}
		keyboard = append(keyboard, row)
	}
	text := NewView(views.Tournament(st.core).Caption(statisticCaption), betsRoundsText).Text()
	return &InterfaceOut{Keyboard: append(keyboard, backRow), Text: text}
}

func statisticMainManagerTournament(main Main, core core.Tournament) *Loading {
	return NewLoading(main.Base, LoadMain(main.Base, main.manager, statisticCaption),
		newLoadStatisticTournament(LoadMain(main.Base, main.manager, statisticCaption), core), loadTournamentText)
}

type (
	loadStatisticTournament struct {
		caller Interface
		core   core.Tournament
	}

	statisticTournament struct {
		Back
		core       core.Tournament
		manager    matchManager
		roundMap   MapRound
		roundSlice []core.Round
	}
)

var statisticMainManager = MainManager{tournament: statisticMainManagerTournament}

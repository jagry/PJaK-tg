package screens

import (
	"PJaK/core"
	"PJaK/views"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

const (
	loadTournamentText = "Идет загрузка туров" + loadingTextSuffix
)

func initTournament(matchMap *MapRound, matchSlice []*core.Round) {
	for _, match := range matchSlice {
		if len(match.Rounds) == 0 {
			(*matchMap)[match.Id] = match
		} else {
			initTournament(matchMap, match.Rounds)
		}
	}
}

func LoadTournament(base Base, section string, caller Interface, core *core.Tournament) *Loading {
	caption := views.Tournament(*core).Caption(section)
	return NewLoading(base, section, caller, NewTournamentFactory(caller, core), caption, loadTournamentText)
}

func NewTournament(section Section, caller, loader Interface, core *core.Tournament, rounds []*core.Round) Tournament {
	roundMap := MapRound{}
	initTournament(&roundMap, rounds)
	return Tournament{Back: NewBack(section, caller, loader), core: core, roundMap: roundMap, roundSlice: rounds}
}

func NewTournamentFactory(caller Interface, core *core.Tournament) TournamentFactory {
	return TournamentFactory{caller: caller, core: core}
}

func (tournament Tournament) Handle(update tgbotapi.Update) (Interface, bool, tgbotapi.Chattable) {
	if update.CallbackQuery != nil {
		if strings.HasPrefix(update.CallbackQuery.Data, betsTournamentIdPrefix) {
			id, fail := strconv.ParseUint(update.CallbackQuery.Data[len(betsTournamentIdPrefix):], 36, 16)
			if fail != nil {
				panic("screens.BetsTournament.Handle:" + fail.Error())
			}
			if round, ok := tournament.roundMap[core.RoundId(id)]; ok {
				factory := NewRoundFactory(tournament.loader, tournament.core, round)
				text := views.Round(*round).Caption(tournament.Section.name, views.Tournament(*tournament.core))
				loading := NewLoading(
					tournament.Base, tournament.Section.name, tournament, factory, text, betsLoadMatchesText)
				return loading, false, tgbotapi.NewCallback(update.CallbackQuery.ID, round.Name)
			} else {
				return nil, false, tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			}
		}
	}
	return tournament.Back.Handle(update)
}

func (tournament Tournament) Out() *InterfaceOut {
	if len(tournament.roundMap) == 0 {
		text := views.Tournament(*tournament.core).Caption(tournament.Section.name)
		return &InterfaceOut{Keyboard: backKeyboard, Text: NewView(text, betsTournamentEmptyText).Text()}
	}
	rounds := core.GetRounds(tournament.roundSlice)
	keyboard := make([][]tgbotapi.InlineKeyboardButton, 0)
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
	text := NewView(views.Tournament(*tournament.core).Caption(tournament.Section.name), betsRoundsText).Text()
	return &InterfaceOut{Keyboard: append(keyboard, backRow), Text: text}
}

func (tf TournamentFactory) Execute(action *Loading) Interface {
	if rounds, fail := core.GetTournamentRounds(tf.core); fail == nil {
		section := NewSection(action.Base, action.Section.name)
		return NewTournament(section, LoadMain(action.Base, action.Section.name), action, tf.core, rounds)
	}
	text := views.Tournament(*tf.core).Caption(action.Section.name)
	return NewError(action.Base, action.Base, text, "!!!")
}

type (
	MapTournament map[core.TournamentId]*core.Tournament

	Tournament struct {
		Back
		core       *core.Tournament
		roundMap   MapRound
		roundSlice []*core.Round
		section    string
	}

	TournamentFactory struct {
		caller Interface
		core   *core.Tournament
	}
)

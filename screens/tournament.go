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

func initTournament(matchMap MapRound, matchSlice []core.Round) MapRound {
	for _, match := range matchSlice {
		if len(match.Rounds) == 0 {
			matchMap[match.Id] = match
		} else {
			initTournament(matchMap, match.Rounds)
		}
	}
	return matchMap
}

func LoadTournament(base Base, caller Interface, manager matchManager, section string, core core.Tournament) *Loading {
	return NewLoading(base, caller, NewTournamentFactory(caller, manager, section, core), loadTournamentText)
}

func NewTournament(b Base, c, l Interface, m matchManager, s string, t core.Tournament, r []core.Round) Tournament {
	roundMap := initTournament(MapRound{}, r)
	return Tournament{Back: NewBack(b, c, l), core: t, manager: m, roundMap: roundMap, roundSlice: r, section: s}
}

func NewTournamentFactory(c Interface, m matchManager, s string, t core.Tournament) TournamentFactory {
	return TournamentFactory{BaseAction: NewBaseAction(s), caller: c, core: t, manager: m}
}

func (tournament Tournament) Handle(update tgbotapi.Update) (Interface, bool, tgbotapi.Chattable) {
	if update.CallbackQuery != nil {
		if strings.HasPrefix(update.CallbackQuery.Data, betsTournamentIdPrefix) {
			id, fail := strconv.ParseUint(update.CallbackQuery.Data[len(betsTournamentIdPrefix):], 36, 16)
			if fail != nil {
				panic("screens.BetsTournament.Handle:" + fail.Error())
			}
			if round, ok := tournament.roundMap[core.RoundId(id)]; ok {
				factory := NewRoundFactory(tournament.loader, tournament.manager,
					tournament.section, tournament.core, round)
				// !!! text := views.Round(round).Caption(tournament.section, views.Tournament(*tournament.core))
				loading := NewLoading(tournament.Base, tournament, factory, betsLoadMatchesText)
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
		text := views.Tournament(tournament.core).Caption(tournament.section)
		return &InterfaceOut{Keyboard: backKeyboard, Text: NewView(text, betsTournamentEmptyText).Text()}
	}
	keyboard, rounds := make([][]tgbotapi.InlineKeyboardButton, 0), core.GetRounds(tournament.roundSlice)
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
	text := NewView(views.Tournament(tournament.core).Caption(tournament.section), betsRoundsText).Text()
	return &InterfaceOut{Keyboard: append(keyboard, backRow), Text: text}
}

func (tf TournamentFactory) Caption() string { return views.Tournament(tf.core).Caption(tf.section) }

func (tf TournamentFactory) Execute(action *Loading) Interface {
	if rounds, fail := core.GetTournamentRounds(tf.core); fail == nil {
		return NewTournament(action.Base, action.caller, action, tf.manager, tf.section, tf.core, rounds)
	}
	return NewError(action.Base, action.Base, tf.Caption(), "!!!")
}

type (
	MapTournament map[core.TournamentId]int

	Tournament struct {
		Back
		core       core.Tournament
		manager    matchManager
		roundMap   MapRound
		roundSlice []core.Round
		section    string
	}

	TournamentFactory struct {
		BaseAction
		caller  Interface
		core    core.Tournament
		manager matchManager
	}
)

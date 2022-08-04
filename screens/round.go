package screens

import (
	"PJaK/core"
	"PJaK/views"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

const (
	roundIdPrefix = "round."
)

func NewRound(base Base, caller, loader Interface, manager matchManager,
	section string, tournament core.Tournament, round core.Round, matches []core.Match) Round {
	matchMap := mapMatch{}
	for index, match := range matches {
		matchMap[match.Id] = index
	}
	return Round{Back: NewBack(base, caller, loader), core: round, manager: manager,
		matchMap: matchMap, matchSlice: matches, section: section, tournament: tournament}
}

func NewRoundFactory(c Interface, m matchManager, s string, t core.Tournament, r core.Round) RoundFactory {
	return RoundFactory{BaseAction: NewBaseAction(s), caller: c, manager: m, core: r, tournament: t}
}

func (round Round) Handle(update tgbotapi.Update) (Interface, bool, tgbotapi.Chattable) {
	if update.CallbackQuery != nil {
		if strings.HasPrefix(update.CallbackQuery.Data, roundIdPrefix) {
			id, fail := strconv.ParseUint(update.CallbackQuery.Data[len(roundIdPrefix):], 36, 16)
			if fail != nil {
				panic("screens.BetsRound.Handle:" + fail.Error())
			}
			if index, ok := round.matchMap[core.MatchId(id)]; ok {
				match := round.matchSlice[index]
				// !!! view := views.Match(match)
				// !!! text := view.Caption(round.section, views.Round(round.core), views.Tournament(round.tournament))
				executor := newLoadMatch(round.loader, round.manager, round.section, round.tournament, round.core, match)
				loading := NewLoading(round.Base, round, executor, loadTournamentText)
				text := views.Match(match).Players("", ":")
				return loading, false, tgbotapi.NewCallback(update.CallbackQuery.ID, text)
			} else {
				return nil, false, tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			}
		}
	}
	return round.Back.Handle(update)
}

func (round Round) Out() *InterfaceOut {
	caption := views.Round(round.core).Caption(round.section, views.Tournament(round.tournament))
	if len(round.matchMap) == 0 {
		return &InterfaceOut{Keyboard: backKeyboard, Text: NewView(caption, betsRoundEmptyText).Text()}
	}
	keyboard := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, match := range round.matchSlice {
		id := roundIdPrefix + strconv.FormatUint(uint64(match.Id), 36)
		text := round.manager.button(match)
		keyboard = append(keyboard, []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(text, id)})
	}
	return &InterfaceOut{Keyboard: append(keyboard, backRow), Text: NewView(caption, betsMatchesText).Text()}
}

func (rf RoundFactory) Caption() string {
	return views.Round(rf.core).Caption(rf.section, views.Tournament(rf.tournament))
}

func (rf RoundFactory) Do(action *Loading) Event {
	if matches, fail := core.GetMatches(rf.core, action.user); fail == nil {
		return NewEvent(NewRound(action.Base, rf.caller, action, rf.manager, rf.section, rf.tournament, rf.core, matches), "")
	}
	return NewEvent(NewError(action.Base, action.Base, "!!!", "!!!"), "")
}

type (
	MapRound map[core.RoundId]core.Round

	Round struct {
		Back
		core       core.Round
		matchMap   mapMatch
		matchSlice []core.Match
		manager    matchManager
		section    string
		tournament core.Tournament
	}

	RoundFactory struct {
		BaseAction
		caller     Interface
		core       core.Round
		manager    matchManager
		tournament core.Tournament
	}
)

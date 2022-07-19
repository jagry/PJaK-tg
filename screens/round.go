package screens

import (
	"PJaK/core"
	"PJaK/views"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
	"time"
)

const (
	roundIdPrefix = "round."
)

func NewRound(s Section, c, l Interface, t *core.Tournament, r *core.Round, m []*core.Match) Round {
	matchMap := MapMatch{}
	for _, match := range m {
		matchMap[match.Id] = match
	}
	return Round{Back: NewBack(s, c, l), core: r, matchMap: matchMap, matchSlice: m, tournament: t}
}

func NewRoundFactory(caller Interface, tournament *core.Tournament, core *core.Round) RoundFactory {
	return RoundFactory{caller: caller, core: core, tournament: tournament}
}

func (round Round) Handle(update tgbotapi.Update) (Interface, bool, tgbotapi.Chattable) {
	if update.CallbackQuery != nil {
		if strings.HasPrefix(update.CallbackQuery.Data, roundIdPrefix) {
			id, fail := strconv.ParseUint(update.CallbackQuery.Data[len(roundIdPrefix):], 36, 16)
			if fail != nil {
				panic("screens.BetsRound.Handle:" + fail.Error())
			}
			if match, ok := round.matchMap[core.MatchId(id)]; ok {
				view := views.Match(*match)
				text := view.Caption(round.Section.name, views.Round(*round.core), views.Tournament(*round.tournament))
				factory := NewMatchFactory(round.loader, round.tournament, round.core, match)
				loading := NewLoading(round.Base, round.Section.name, round, factory, text, loadTournamentText)
				text = views.Match(*match).Players("", ":")
				return loading, false, tgbotapi.NewCallback(update.CallbackQuery.ID, text)
			} else {
				return nil, false, tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			}
		}
	}
	return round.Back.Handle(update)
}

func (round Round) Out() *InterfaceOut {
	caption := views.Round(*round.core).Caption(round.Section.name, views.Tournament(*round.tournament))
	if len(round.matchMap) == 0 {
		return &InterfaceOut{Keyboard: backKeyboard, Text: NewView(caption, betsRoundEmptyText).Text()}
	}
	keyboard := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, match := range round.matchSlice {
		id := roundIdPrefix + strconv.FormatUint(uint64(match.Id), 36)
		view := views.Match(*match)
		text := " " + view.Players(matchUndefined, ":") + " " + view.Bet("")
		if match.Team1.Result() == nil && match.Team2.Result() == nil {
			if match.Time().Sub(time.Now()) > 0 {
				text = "\U0001F7E1" + text + " / " + view.Time()
			} else {
				text = "\U0001F7E2" + text
			}
		} else {
			if match.Time().Sub(time.Now()) > 0 {
				text = "🔵" + text
			} else {
				text = "\U0001F7E0" + text + " / " + view.Result("")
			}
		}
		keyboard = append(keyboard, []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(text, id)})
	}
	return &InterfaceOut{Keyboard: append(keyboard, backRow), Text: NewView(caption, betsMatchesText).Text()}
}

func (rf RoundFactory) Execute(action *Loading) Interface {
	if matches, fail := core.GetMatches(rf.core, action.User()); fail == nil {
		return NewRound(action.Section, rf.caller, action, rf.tournament, rf.core, matches)
	}
	return NewError(action.Base, action.Base, "!!!", "!!!")
}

type (
	MapRound map[core.RoundId]*core.Round

	Round struct {
		Back
		core       *core.Round
		matchMap   MapMatch
		matchSlice []*core.Match
		tournament *core.Tournament
	}

	RoundFactory struct {
		caller     Interface
		core       *core.Round
		matchMap   MapMatch
		matchSlice []*core.Match
		tournament *core.Tournament
	}
)

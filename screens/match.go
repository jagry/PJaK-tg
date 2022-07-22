package screens

import (
	"PJaK/core"
	"PJaK/views"
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
	"strconv"
	"strings"
)

const (
	matchSaveId   = "match.save"
	matchSaveText = "💾 Сохранить"

	matchUndefined  = "<i>не определена</i>"
	matchUndefined1 = "Команда 1"
	matchUndefined2 = "Команда 2"
)

func matchHandle(c tgbotapi.CallbackQuery, p, n string, s matchManagerTeam, t *core.MatchTeam) tgbotapi.Chattable {
	idString := c.Data[len(p):]
	idInt64, fail := strconv.ParseUint(idString, 10, 8)
	if fail != nil {
		panic("screens.BetsRound.Handle:" + fail.Error())
	}
	s(t, byte(idInt64))
	return tgbotapi.NewCallback(c.ID, views.MatchTeam(*t).Short(n)+" "+idString)
}

func newLoadMatch(c Interface, m matchManager, s string, t core.Tournament, r core.Round, cm core.Match) loadMatch {
	return loadMatch{BaseAction: NewBaseAction(s), caller: c, core: cm, manager: m, round: r, tournament: t}
}

func newMatch(b Base, c, l Interface, m matchManager, s string, t core.Tournament, r core.Round, cm core.Match) match {
	return match{Back: NewBack(b, c, l), core: cm, manager: m, round: r, section: s, tournament: t}
}

func newMatchManager(
	button matchManagerButton,
	modify matchManagerModify,
	save matchManagerSave,
	screen matchManagerScreen,
	team matchManagerTeam) matchManager {
	return matchManager{button: button, modify: modify, save: save, screen: screen, team: team}
}

func newSaveMatch(base Base, success, error Interface, manager matchManager, section string,
	core core.Match, round core.Round, tournament core.Tournament) saveMatch {
	return saveMatch{BaseAction: NewBaseAction(section), base: base, core: core, error: error,
		manager: manager, round: round, success: success, tournament: tournament}
}

func (lm loadMatch) Caption() string {
	return views.Match(lm.core).Caption(lm.section, views.Round(lm.round), views.Tournament(lm.tournament))
}

func (lm loadMatch) Execute(action *Loading) Interface {
	if match, fail := core.GetMatch(lm.core.Id, action.User()); fail == nil {
		return newMatch(action.Base, lm.caller, action, lm.manager, lm.section, lm.tournament, lm.round, match)
	}
	return NewError(action.Base, action.Base, betsCaption, "!!!")
}

func (match match) Handle(update tgbotapi.Update) (Interface, bool, tgbotapi.Chattable) {
	if update.CallbackQuery != nil {
		if strings.HasPrefix(update.CallbackQuery.Data, betsMatchTeam1IdPrefix) {
			callback := matchHandle(*update.CallbackQuery, betsMatchTeam1IdPrefix,
				matchUndefined1, match.manager.team, &match.core.Team1)
			return match, false, callback
		} else if strings.HasPrefix(update.CallbackQuery.Data, betsMatchTeam2IdPrefix) {
			callback := matchHandle(*update.CallbackQuery, betsMatchTeam2IdPrefix,
				matchUndefined2, match.manager.team, &match.core.Team2)
			return match, false, callback
		} else if update.CallbackQuery.Data == matchSaveId {
			factory := newSaveMatch(match.Base, match.caller, match.loader,
				match.manager, match.section, match.core, match.round, match.tournament)
			loading := NewLoading(match.Base, match.loader, factory, "!!!ggg")
			return loading, false, tgbotapi.NewCallback(update.CallbackQuery.ID, matchSaveText)
		}
	} else if update.Message != nil && !update.Message.IsCommand() {
		dmc := tgbotapi.NewDeleteMessage(update.FromChat().ID, update.Message.MessageID)
		regexp := regexp.MustCompile("^(\\d+)[-:\\s](\\d+)$")
		if regexpResult := regexp.FindStringSubmatch(update.Message.Text); len(regexpResult) == 3 {
			newGoal1, fail := strconv.ParseUint(regexpResult[1], 10, 8)
			if fail != nil {
				if errors.Is(fail, strconv.ErrRange) {
					return NewError(match.Base, match, baseUserCaptionText, "!!!! Число голов должно быть в диапазоне от 0 до 255 включительно"), false, dmc
				}
				// TODO: Проверить другую ошибку
			}
			newGoal2, fail := strconv.ParseUint(regexpResult[2], 10, 8)
			if fail != nil {
				if errors.Is(fail, strconv.ErrRange) {
					return NewError(match.Base, match, baseUserCaptionText, "!!!! Число голов должно быть в диапазоне от 0 до 255 включительно"), false, dmc
				}
				// TODO: Проверить другую ошибку
			}
			modify, oldGoal1, oldGoal2 := match.manager.modify(match.core)
			if modify {
				if oldGoal1 == nil || oldGoal2 == nil || *oldGoal1 != byte(newGoal1) || *oldGoal2 != byte(newGoal2) {
					match.manager.team(&match.core.Team1, byte(newGoal1))
					match.manager.team(&match.core.Team2, byte(newGoal2))
					return match, false, dmc
				}
			} else {
				return NewError(match.Base, match.loader, "!!! njj", "!!! njj"), false, dmc
			}
			return nil, false, dmc
		}
		return NewError(match.Base, match, baseUserCaptionText, "!!!! X-X"), false, dmc
	}
	return match.Back.Handle(update)
}

func (match match) Out() *InterfaceOut {
	caption := views.Round(match.round).Caption(match.section, views.Tournament(match.tournament))
	text, keys := match.manager.screen(match.core)
	return &InterfaceOut{Keyboard: append(keys, match.Back.buttons), Text: NewView(caption, text).Text()}
}

func (sm saveMatch) Caption() string {
	return views.Match(sm.core).Caption(sm.section, views.Round(sm.round), views.Tournament(sm.tournament))
}

func (sm saveMatch) Execute(action *Loading) Interface {
	if fail := sm.manager.save(sm.core, action.user); fail == nil {
		return sm.success
	}
	return NewError(sm.base, sm.error, sm.Caption(), "!!! Не удалось сохранить данные")
}

type (
	loadMatch struct {
		BaseAction
		caller     Interface
		core       core.Match
		manager    matchManager
		round      core.Round
		tournament core.Tournament
	}

	mapMatch map[core.MatchId]int

	match struct {
		Back
		core       core.Match
		manager    matchManager
		round      core.Round
		section    string
		tournament core.Tournament
	}

	matchManager struct {
		button matchManagerButton
		modify matchManagerModify
		save   matchManagerSave
		screen matchManagerScreen
		team   matchManagerTeam
	}

	matchManagerButton func(core.Match) string

	matchManagerModify func(core.Match) (bool, *byte, *byte)

	matchManagerSave func(core.Match, int8) error

	matchManagerScreen func(core.Match) (string, [][]tgbotapi.InlineKeyboardButton)

	matchManagerTeam func(*core.MatchTeam, byte)

	saveMatch struct {
		BaseAction
		base           Base
		error, success Interface
		core           core.Match
		manager        matchManager
		round          core.Round
		tournament     core.Tournament
	}
)

var (
	matchSaveButton = tgbotapi.NewInlineKeyboardButtonData(matchSaveText, matchSaveId)
)

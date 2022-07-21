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

func matchHandle(c tgbotapi.CallbackQuery, p, n string, s SectionTeamModify, t *core.MatchTeam) tgbotapi.Chattable {
	idString := c.Data[len(p):]
	idInt64, fail := strconv.ParseUint(idString, 10, 8)
	if fail != nil {
		panic("screens.BetsRound.Handle:" + fail.Error())
	}
	s(t, byte(idInt64))
	return tgbotapi.NewCallback(c.ID, views.MatchTeam(*t).Short(n)+" "+idString)
}

func NewMatch(b Base, c, l Interface, s Section, t *core.Tournament, r *core.Round, m *core.Match) Match {
	return Match{Back: NewBack(b, c, l), core: m, round: r, section: s, tournament: t}
}

func NewMatchFactory(c Interface, s Section, t *core.Tournament, r *core.Round, m *core.Match) MatchFactory {
	return MatchFactory{caller: c, core: m, round: r, section: s, tournament: t}
}

func NewMatchSave(b Base, s, e Interface, c *core.Match, r *core.Round, t *core.Tournament, sc Section) MatchSave {
	return MatchSave{base: b, core: c, error: e, round: r, section: sc, success: s, tournament: t}
}

func (match Match) Handle(update tgbotapi.Update) (Interface, bool, tgbotapi.Chattable) {
	if update.CallbackQuery != nil {
		if strings.HasPrefix(update.CallbackQuery.Data, betsMatchTeam1IdPrefix) {
			callback := matchHandle(
				*update.CallbackQuery,
				betsMatchTeam1IdPrefix,
				matchUndefined1,
				match.section.teamModify,
				&match.core.Team1)
			return match, false, callback
		} else if strings.HasPrefix(update.CallbackQuery.Data, betsMatchTeam2IdPrefix) {
			callback := matchHandle(
				*update.CallbackQuery,
				betsMatchTeam2IdPrefix,
				matchUndefined2,
				match.section.teamModify,
				&match.core.Team2)
			return match, false, callback
		} else if update.CallbackQuery.Data == matchSaveId {
			factory := NewMatchSave(
				match.Base,
				match.caller,
				match.loader,
				// !!! match.user,
				match.core,
				match.round,
				match.tournament,
				match.section)
			loading := NewLoading(match.Base, match.loader, factory, "ddd", "ggg")
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
			modify, oldGoal1, oldGoal2 := match.section.modify(*match.core)
			if modify {
				if oldGoal1 == nil || oldGoal2 == nil || *oldGoal1 != byte(newGoal1) || *oldGoal2 != byte(newGoal2) {
					match.section.teamModify(&match.core.Team1, byte(newGoal1))
					match.section.teamModify(&match.core.Team2, byte(newGoal2))
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

func (match Match) Out() *InterfaceOut {
	caption := views.Round(*match.round).Caption(match.section.name, views.Tournament(*match.tournament))
	text, keys := match.section.screen(*match.core)
	return &InterfaceOut{Keyboard: append(keys, match.Back.buttons), Text: NewView(caption, text).Text()}
}

func (mf MatchFactory) Execute(action *Loading) Interface {
	if match, fail := core.GetMatch(mf.core.Id, action.User()); fail == nil {
		return NewMatch(action.Base, mf.caller, action, mf.section, mf.tournament, mf.round, &match)
	}
	return NewError(action.Base, action.Base, betsCaption, "!!!")
}

func (ms MatchSave) Execute(action *Loading) Interface {
	if fail := ms.section.save(*ms.core, action.user); fail == nil {
		return ms.success
	}
	tournament := views.Tournament(*ms.tournament)
	caption := views.Match(*ms.core).Caption(ms.section.name, views.Round(*ms.round), tournament)
	return NewError(ms.base, ms.error, caption, "!!! Не удалось сохранить данные")
}

type (
	MapMatch map[core.MatchId]*core.Match

	Match struct {
		Back
		core       *core.Match
		round      *core.Round
		section    Section
		tournament *core.Tournament
	}

	MatchFactory struct {
		caller     Interface
		core       *core.Match
		round      *core.Round
		section    Section
		tournament *core.Tournament
	}

	MatchSave struct {
		base       Base
		core       *core.Match
		error      Interface
		round      *core.Round
		section    Section
		success    Interface
		tournament *core.Tournament
	}
)

var (
	matchSaveButton = tgbotapi.NewInlineKeyboardButtonData(matchSaveText, matchSaveId)
)

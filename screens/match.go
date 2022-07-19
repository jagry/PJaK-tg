package screens

import (
	"PJaK/core"
	"PJaK/views"
	"errors"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	matchUndefined  = "<i>не определена</i>"
	matchUndefined1 = "Команда 1"
	matchUndefined2 = "Команда 2"
)

func matchHandle(cq tgbotapi.CallbackQuery, prefix, null string, team *core.MatchTeam) tgbotapi.Chattable {
	idString := cq.Data[len(prefix):]
	idInt64, fail := strconv.ParseUint(idString, 10, 8)
	if fail != nil {
		panic("screens.BetsRound.Handle:" + fail.Error())
	}
	team.SetBet(byte(idInt64))
	return tgbotapi.NewCallback(cq.ID, views.MatchTeam(*team).Short(null)+" "+idString)
}

func NewMatch(s Section, c, l Interface, t *core.Tournament, r *core.Round, m *core.Match) Match {
	return Match{Back: NewBack(s, c, l), core: m, round: r, tournament: t}
}

func NewMatchFactory(caller Interface, tournament *core.Tournament, round *core.Round, core *core.Match) MatchFactory {
	return MatchFactory{caller: caller, core: core, round: round, tournament: tournament}
}

func NewMatchSave(b Base, s, e Interface, u int8, c *core.Match, r *core.Round, t *core.Tournament) MatchSave {
	return MatchSave{base: b, core: c, error: e, round: r, success: s, tournament: t, user: u}
}

func (match Match) Handle(update tgbotapi.Update) (Interface, bool, tgbotapi.Chattable) {
	if update.CallbackQuery != nil {
		if strings.HasPrefix(update.CallbackQuery.Data, betsMatchTeam1IdPrefix) {
			callback := matchHandle(*update.CallbackQuery, betsMatchTeam1IdPrefix, matchUndefined1, &match.core.Team1)
			return match, false, callback
		} else if strings.HasPrefix(update.CallbackQuery.Data, betsMatchTeam2IdPrefix) {
			callback := matchHandle(*update.CallbackQuery, betsMatchTeam2IdPrefix, matchUndefined2, &match.core.Team2)
			return match, false, callback
		} else if update.CallbackQuery.Data == betsSaveId {
			factory := NewMatchSave(
				match.Base, match.caller, match.loader, match.user, match.core, match.round, match.tournament)
			loading := NewLoading(match.Base, match.Section.name, match.loader, factory, "ddd", "ggg")
			return loading, false, tgbotapi.NewCallback(update.CallbackQuery.ID, betsSaveText)
		}
	} else if update.Message != nil && !update.Message.IsCommand() {
		// TODO: Проверить не начался ли матч
		dmc := tgbotapi.NewDeleteMessage(update.FromChat().ID, update.Message.MessageID)
		regexp := regexp.MustCompile("^(\\d+)-(\\d+)$")
		if regexpResult := regexp.FindStringSubmatch(update.Message.Text); len(regexpResult) == 3 {
			goal1, fail := strconv.ParseUint(regexpResult[1], 10, 8)
			if fail != nil {
				if errors.Is(fail, strconv.ErrRange) {
					return NewError(match.Base, match, baseUserCaptionText, "!!!! Число голов должно быть в диапазоне от 0 до 255 включительно"), false, dmc
				}
				// TODO: Проверить другую ошибку
			}
			goal2, fail := strconv.ParseUint(regexpResult[2], 10, 8)
			if fail != nil {
				if errors.Is(fail, strconv.ErrRange) {
					return NewError(match.Base, match, baseUserCaptionText, "!!!! Число голов должно быть в диапазоне от 0 до 255 включительно"), false, dmc
				}
				// TODO: Проверить другую ошибку
			}
			team1, team2 := &match.core.Team1, &match.core.Team2
			bet1, bet2 := team1.Bet(), team2.Bet()
			if bet1 == nil || bet2 == nil || *bet1 != byte(goal1) || *bet2 != byte(goal2) {
				team1.SetBet(byte(goal1))
				team2.SetBet(byte(goal2))
				return match, false, dmc
			}
			return nil, false, dmc
		}
		return NewError(match.Base, match, baseUserCaptionText, "!!!! X-X"), false, dmc
	}
	return match.Back.Handle(update)
}

func (match Match) Out() *InterfaceOut {
	caption := views.Round(*match.round).Caption(betsCaption, views.Tournament(*match.tournament))
	helper1, helper2 := views.MatchTeam(match.core.Team1), views.MatchTeam(match.core.Team2)
	keys := make([][]tgbotapi.InlineKeyboardButton, 0)
	text := "<pre>" +
		helper1.Table(matchUndefined1, betsMatchTeam1IdPrefix, betsPrefixSelected, betsSuffixSelected) + "\n" +
		helper2.Table(matchUndefined2, betsMatchTeam2IdPrefix, betsPrefixSelected, betsSuffixSelected) + "</pre>"
	if match.core.Team1.Result() == nil && match.core.Team2.Result() == nil {
		if match.core.Time().Sub(time.Now()) < 0 {
			text += "Матч начался"
		} else {
			text += "Начало матча: " + views.Match(*match.core).Time()
			keys = append(keys, helper1.Keys("", betsMatchTeam1IdPrefix))
			keys = append(keys, helper2.Keys("", betsMatchTeam2IdPrefix))
			keys = append(keys, []tgbotapi.InlineKeyboardButton{betsSaveButton})
		}
	} else {
		if match.core.Time().Sub(time.Now()) < 0 {
			text += "Счет матча: " + views.Match(*match.core).Result("")
		} else {
			text = "Хуйня какая-то: матч не начался, а счет есть"
		}
	}
	keys = append(keys, match.Back.buttons)
	return &InterfaceOut{Keyboard: keys, Text: NewView(caption, text).Text()}
}

func (mf MatchFactory) Execute(action *Loading) Interface {
	if match, fail := core.GetMatch(mf.core.Id, action.User()); fail == nil {
		return NewMatch(action.Section, mf.caller, action, mf.tournament, mf.round, &match)
	}
	return NewError(action.Base, action.Base, betsCaption, "!!!")
}

func (ms MatchSave) Execute(action *Loading) Interface {
	if fail := core.SaveBets(*ms.core.Team1.Bet(), *ms.core.Team2.Bet(), ms.core.Id, ms.user); fail == nil {
		return ms.success
	}
	tournament := views.Tournament(*ms.tournament)
	caption := views.Match(*ms.core).Caption(betsCaption, views.Round(*ms.round), tournament)
	return NewError(ms.base, ms.error, caption, "!!! Не удалось сохранить данные")
}

type (
	MapMatch map[core.MatchId]*core.Match

	Match struct {
		Back
		core       *core.Match
		round      *core.Round
		tournament *core.Tournament
	}

	MatchFactory struct {
		caller     Interface
		core       *core.Match
		round      *core.Round
		tournament *core.Tournament
	}

	MatchSave struct {
		base       Base
		core       *core.Match
		error      Interface
		round      *core.Round
		success    Interface
		tournament *core.Tournament
		user       int8
	}
)

var (
	betsSaveButton = tgbotapi.NewInlineKeyboardButtonData(betsSaveText, betsSaveId)
)

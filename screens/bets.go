package screens

import (
	"PJaK/core"
	"PJaK/views"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

const (
	betsCaption      = betsCaptionEmoji + " " + betsCaptionText
	betsCaptionText  = "Прогнозы"
	betsCaptionEmoji = "🎲"
)

const (
	betsRowCount = 6

	betsMatchIdPrefix      = "match."
	betsMatchTeam1IdPrefix = betsMatchIdPrefix + "1."
	betsMatchTeam2IdPrefix = betsMatchIdPrefix + "2."
	betsTournamentIdPrefix = "tournament."

	betsPrefixSelected = ""
	betsSuffixSelected = ""

	betsEmptyText           = "Еще нет футбольных турниров в текущих сезонах"
	betsLoadMatchesText     = "Идет загрузка матчей" + loadingTextSuffix
	betsMatchesText         = "Матчи:"
	betsRoundsText          = "Туры:"
	betsRoundEmptyText      = "Расписание матчей будет определено позже"
	betsTournamentsText     = "Футбольные турниры:"
	betsTournamentEmptyText = "Расписание туров будет определено позже"
)

func bets() Section {
	return NewSection(betsCaption, betsButton, betsModify, betsSave, betsScreen, betsTeamModify)
}

func betsButton(match core.Match) string {
	view := views.Match(match)
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
	return text
}

func betsModify(match core.Match) (bool, *byte, *byte) {
	if match.Time().Sub(time.Now()) < 0 {
		return false, nil, nil
	}
	return true, match.Team1.Bet(), match.Team2.Bet()
}

func betsSave(match core.Match, user int8) error {
	return core.SaveBets(*match.Team1.Bet(), *match.Team2.Bet(), match.Id, user)
}

func betsScreen(match core.Match) (string, [][]tgbotapi.InlineKeyboardButton) {
	bet1, bet2 := match.Team1.Bet(), match.Team2.Bet()
	helper1, helper2 := views.MatchTeam(match.Team1), views.MatchTeam(match.Team2)
	keys := make([][]tgbotapi.InlineKeyboardButton, 0, 3)
	text := "<pre>" + helper1.Table(matchUndefined1, bet1) + "\n" + helper2.Table(matchUndefined2, bet2) + "</pre>"
	if match.Time().Sub(time.Now()) > 0 {
		if bet1 == nil && bet2 == nil {
			text += "Начало матча: " + views.Match(match).Time()
		} else {
			text += "Хуйня какая-то: матч не начался, а счет есть"
		}
		keys = append(keys, helper1.Keys("", betsMatchTeam1IdPrefix, bet1))
		keys = append(keys, helper2.Keys("", betsMatchTeam2IdPrefix, bet2))
		keys = append(keys, []tgbotapi.InlineKeyboardButton{matchSaveButton})
	} else if bet1 == nil && bet2 == nil {
		text += "Матч начался"
	} else {
		text += "Счет матча: " + views.Match(match).Result("")
	}
	return text, keys
}

func betsTeamModify(team *core.MatchTeam, bet byte) { team.SetBet(bet) }

func (bm BetsManager) Section(main Main, tournament *core.Tournament) string {
	return bm.section
}

func (bm BetsManager) Tournament(main Main, tournament *core.Tournament) *Loading {
	return LoadTournament(main.Base, bets(), LoadMain(main.Base, bets()), tournament)
}

type BetsManager struct {
	section string
}

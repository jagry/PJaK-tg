package screens

import (
	"PJaK/core"
	"PJaK/views"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

const (
	resultsCaption      = resultsCaptionEmoji + " " + resultsCaptionText
	resultsCaptionEmoji = "⚽️"
	resultsCaptionText  = "Результаты"
)

func resultButton(match core.Match) string {
	view := views.Match(match)
	text := view.Players(matchUndefined, ":")
	if match.Team1.Result() == nil && match.Team2.Result() == nil {
		if match.Time().Sub(time.Now()) > 0 {
			text += " \U0001F7E1 " + view.Time()
		} else {
			text += " \U0001F7E2 " + view.Result("")
		}
	} else {
		if match.Time().Sub(time.Now()) > 0 {
			text += " 🔵"
		} else {
			text += " \U0001F7E0 " + view.Result("")
		}
	}
	return text
}

func resultsMainManagerTournament(main Main, tournament core.Tournament) *Loading {
	load := loadMain(main.Base, main.manager, resultsCaption)
	return LoadTournament(main.Base, load, resultsMatchManager, main.section, tournament)
}

func resultsModify(match core.Match) (bool, *byte, *byte) {
	if match.Time().Sub(time.Now()) > 0 {
		return false, nil, nil
	}
	return true, match.Team1.Result(), match.Team2.Result()
}

func resultsSave(match core.Match, user int8) (byte, byte, error) {
	result1, result2 := *match.Team1.Result(), *match.Team2.Result()
	return result1, result2, core.SaveResult(result1, result2, match.Id, user)
}

func resultsScreen(match core.Match) (string, [][]tgbotapi.InlineKeyboardButton) {
	helper1, helper2 := views.MatchTeam(match.Team1), views.MatchTeam(match.Team2)
	result1, result2 := match.Team1.Result(), match.Team2.Result()
	keys := make([][]tgbotapi.InlineKeyboardButton, 0, 3)
	text := "<pre>" + helper1.Table(matchUndefined1, result1) +
		"\n" + helper2.Table(matchUndefined2, result2) + "</pre>"
	if match.Time().Sub(time.Now()) > 0 {
		if result1 == nil || result2 == nil {
			text += "Начало матча: " + views.Match(match).Time()
		} else {
			text = "Хуйня какая-то: матч не начался, а счет есть"
		}
	} else {
		keys = append(keys, helper1.Keys("", betsMatchTeam1IdPrefix, result1))
		keys = append(keys, helper2.Keys("", betsMatchTeam2IdPrefix, result2))
		keys = append(keys, []tgbotapi.InlineKeyboardButton{matchSaveButton})
	}
	return text, keys
}

func resultsTeam(team *core.MatchTeam, result byte) {
	team.SetResult(result)
}

var (
	resultsMainManager  = NewMainManager(resultsMainManagerTournament)
	resultsMatchManager = newMatchManager(resultButton, resultsModify, resultsSave, resultsScreen, resultsTeam)
)

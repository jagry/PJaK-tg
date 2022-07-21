package screens

import (
	"PJaK/core"
	"PJaK/views"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

const (
	resultsCaption      = resultsCaptionEmoji + " " + resultsCaptionText
	resultsCaptionText  = "Результаты"
	resultsCaptionEmoji = "⚽️"
)

func results() Section {
	return NewSection(resultsCaption, resultButton, resultsModify, resultsSave, resultsScreen, resultsTeamModify)
}

func resultButton(match core.Match) string {
	view := views.Match(match)
	text := " " + view.Players(matchUndefined, ":") + " "
	if match.Team1.Result() == nil && match.Team2.Result() == nil {
		if match.Time().Sub(time.Now()) > 0 {
			text = "\U0001F7E1" + text + " / " + view.Time()
		} else {
			text = "\U0001F7E2" + text + view.Result("")
		}
	} else {
		if match.Time().Sub(time.Now()) > 0 {
			text = "🔵" + text
		} else {
			text = "\U0001F7E0" + text + view.Result("")
		}
	}
	return text
}

func resultsModify(match core.Match) (bool, *byte, *byte) {
	if match.Time().Sub(time.Now()) > 0 {
		return false, nil, nil
	}
	return true, match.Team1.Result(), match.Team2.Result()
}

func resultsSave(match core.Match, user int8) error {
	return core.SaveResult(*match.Team1.Result(), *match.Team2.Result(), match.Id, user)
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

func resultsTeamModify(team *core.MatchTeam, result byte) {
	team.SetResult(result)
}

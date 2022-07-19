package views

import (
	"PJaK/core"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math"
	"strconv"
	"strings"
)

const (
	matchKeyCount = 6

	matchTeamNameLen = 32
	matchTeamGoalLen = 3
	matchTeamGoalNil = "-"
)

func (match Match) Bet(separator string) (result string) {
	if bet := core.Match(match).Team1.Bet(); bet == nil {
		result = matchTeamGoalNil
	} else {
		result = strconv.Itoa(int(*bet))
	}
	result += separator + ":" + separator
	if bet := core.Match(match).Team2.Bet(); bet == nil {
		result += matchTeamGoalNil
	} else {
		result += strconv.Itoa(int(*bet))
	}
	return
}

func (match Match) Caption(section string, round Round, tournament Tournament) string {
	return round.Caption(section, tournament) + formCaptionDelimiter +
		MatchTeam(match.Team1).Full("-") + ":" + MatchTeam(match.Team2).Full("-")
}

func (match Match) Players(null, separator string) string {
	return MatchTeam(core.Match(match).Team1).Short(null) + separator + MatchTeam(core.Match(match).Team2).Short(null)
}

func (match Match) Result(separator string) (result string) {
	if result1 := core.Match(match).Team1.Result(); result1 == nil {
		result = matchTeamGoalNil
	} else {
		result = strconv.Itoa(int(*result1))
	}
	result += separator + ":" + separator
	if result2 := core.Match(match).Team1.Result(); result2 == nil {
		result += matchTeamGoalNil
	} else {
		result += strconv.Itoa(int(*result2))
	}
	return
}

func (match Match) Time() string {
	time := core.Match(match).Time()
	text := strconv.Itoa(time.Day()) + "."
	switch time.Month() {
	case 1:
		text += "Ⅰ"
	case 2:
		text += "Ⅱ"
	case 3:
		text += "Ⅲ"
	case 4:
		text += "Ⅳ"
	case 5:
		text += "Ⅴ"
	case 6:
		text += "Ⅵ"
	case 7:
		text += "Ⅶ"
	case 8:
		text += "Ⅷ"
	case 9:
		text += "Ⅸ"
	case 10:
		text += "Ⅹ"
	case 11:
		text += "Ⅺ"
	case 12:
		text += "Ⅻ"
	}
	text += " "
	temp := time.Hour()
	if temp < 10 {
		text += "0"
	}
	text += strconv.Itoa(temp)
	temp = time.Minute()
	switch temp {
	case 0:
		text += "↉"
	case 10:
		text += "⅙"
	case 15:
		text += "¼"
	case 20:
		text += "⅓"
	case 30:
		text += "½"
	case 40:
		text += "⅔"
	case 45:
		text += "¾"
	case 50:
		text += "⅚"
	default:
		text += "."
		if temp < 10 {
			text += "0"
		}
		strconv.Itoa(temp)
	}
	return text
}

func (mt MatchTeam) Keys(null, prefixId string) []tgbotapi.InlineKeyboardButton {
	bet, keys, start := core.MatchTeam(mt).Bet(), make([]tgbotapi.InlineKeyboardButton, matchKeyCount), 0
	if bet == nil {
		for counter := 0; counter < matchKeyCount; counter++ {
			valueStr := strconv.Itoa(counter)
			keys[counter] = tgbotapi.NewInlineKeyboardButtonData(valueStr, prefixId+valueStr)
		}
	} else {
		start = int(*bet) - (matchKeyCount+1)>>1
		if start < 0 {
			start = 0
		} else if start > math.MaxUint8-matchKeyCount {
			start = math.MaxUint8 - matchKeyCount
		}
		for counter := 0; counter < int(*bet)-start; counter++ {
			value := strconv.Itoa(counter + start)
			keys[counter] = tgbotapi.NewInlineKeyboardButtonData(value, prefixId+value)
		}
		for counter := int(*bet) - start + 1; counter < matchKeyCount+1; counter++ {
			value := strconv.Itoa(counter + start)
			keys[counter-1] = tgbotapi.NewInlineKeyboardButtonData(value, prefixId+value)
		}
	}
	return keys
}

func (mt MatchTeam) Full(null string) string {
	if name := core.MatchTeam(mt).Full(); name != nil {
		return *name
	}
	return null
}

func (mt MatchTeam) Short(null string) string {
	if name := core.MatchTeam(mt).Short(); name != nil {
		return *name
	}
	return null
}

func (mt MatchTeam) Table(null, prefixId, prefix, suffix string) string {
	bet, text := core.MatchTeam(mt).Bet(), mt.Full(null)
	text += strings.Repeat(" ", matchTeamNameLen-len([]rune(text)))
	if bet == nil {
		text += strings.Repeat(" ", matchTeamGoalLen-len([]rune(matchTeamGoalNil))) + matchTeamGoalNil
	} else {
		goal := strconv.Itoa(int(*bet))
		text += strings.Repeat(" ", matchTeamGoalLen-len([]rune(goal))) + goal
	}
	return text
}

type (
	Match core.Match

	MatchTeam core.MatchTeam
)

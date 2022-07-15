package helpers

import (
	"PJaK/core"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"math"
	"strconv"
	"strings"
)

//import (
//	"PJaK/core"
//	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
//	"strconv"
//	"strings"
//)
//
const (
	matchKeyCount = 6

	matchTeamNameLen = 32
	matchTeamGoalLen = 3
	matchTeamGoalNil = "-"
)

func (match Match) Bet(separator string) (result string) {
	if bet := core.Match(match).Team1.Bet(); bet == nil {
		result = "-"
	} else {
		result = strconv.Itoa(int(*bet))
	}
	result += separator + ":" + separator
	if bet := core.Match(match).Team1.Bet(); bet == nil {
		result += "-"
	} else {
		result += strconv.Itoa(int(*bet))
	}
	return
}

func (match Match) Players(null, separator string) string {
	return MatchTeam(core.Match(match).Team1).Name(null) + separator + MatchTeam(core.Match(match).Team2).Name(null)
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

func (mt MatchTeam) Name(null string) string {
	if name := core.MatchTeam(mt).Name(); name != nil {
		return *name
	}
	return null
}

func (mt MatchTeam) Table(null, prefixId, prefix, suffix string) (string, []tgbotapi.InlineKeyboardButton) {
	bet, keys, start, text := core.MatchTeam(mt).Bet(), make([]tgbotapi.InlineKeyboardButton, matchKeyCount), 0, mt.Name(null)
	text += strings.Repeat(" ", matchTeamNameLen-len([]rune(text)))
	if bet == nil {
		start = 0
		text += strings.Repeat(" ", matchTeamGoalLen-len([]rune(matchTeamGoalNil))) + matchTeamGoalNil
	} else {
		start = int(*bet) - matchKeyCount>>1
		if start < 0 {
			start = 0
		} else if start > math.MaxUint8-matchKeyCount {
			start = math.MaxUint8 - matchKeyCount
		}
		goal := strconv.Itoa(int(*bet))
		text += strings.Repeat(" ", matchTeamGoalLen-len([]rune(goal))) + goal
	}
	for counter := 0; counter < matchKeyCount; counter++ {
		valueInt := start + counter
		valueStr := strconv.Itoa(valueInt)
		if bet != nil && int(*bet) == valueInt {
			keys[counter] = tgbotapi.NewInlineKeyboardButtonData(prefix+valueStr+suffix, prefixId+valueStr)
		} else {
			keys[counter] = tgbotapi.NewInlineKeyboardButtonData(valueStr, prefixId+valueStr)
		}
	}
	return text, keys
}

type (
	Match core.Match

	MatchTeam core.MatchTeam
)

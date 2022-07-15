package screens

//import (
//	"PJaK/core"
//	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
//	"strconv"
//	"strings"
//)
//
//const (
//	matchKeyCount       = 6
//	matchOffsetKeyCount = 2
//	matchSelectKeyLeft  = "<"
//	matchSelectKeyRight = ">"
//
//	matchTeamNameLen = 20
//	matchTeamGoalLen = 2
//	matchTeamGoalNil = "-"
//)
//
//func (mh Match) Players(separator string) string {
//	return core.MatchTeam(mh).Team1.Name() + separator + ":" + separator + match.Team2.Name()
//}
//
//func (mh MatchHelper) TimeString() string {
//	text := strconv.Itoa(match.Time.Day()) + "."
//	switch match.Time.Month() {
//	case 1:
//		text += "Ⅰ"
//	case 2:
//		text += "Ⅱ"
//	case 3:
//		text += "Ⅲ"
//	case 4:
//		text += "Ⅳ"
//	case 5:
//		text += "Ⅴ"
//	case 6:
//		text += "Ⅵ"
//	case 7:
//		text += "Ⅶ"
//	case 8:
//		text += "Ⅷ"
//	case 9:
//		text += "Ⅸ"
//	case 10:
//		text += "Ⅹ"
//	case 11:
//		text += "Ⅺ"
//	case 12:
//		text += "Ⅻ"
//	}
//	text += " "
//	temp := match.Time.Hour()
//	if temp < 10 {
//		text += "0"
//	}
//	text += strconv.Itoa(temp)
//	temp = match.Time.Minute()
//	switch temp {
//	case 0:
//		text += "↉"
//	case 10:
//		text += "⅙"
//	case 15:
//		text += "¼"
//	case 20:
//		text += "⅓"
//	case 30:
//		text += "½"
//	case 40:
//		text += "⅔"
//	case 45:
//		text += "¾"
//	case 50:
//		text += "⅚"
//	default:
//		text += "."
//		if temp < 10 {
//			text += "0"
//		}
//		strconv.Itoa(temp)
//	}
//	return text
//}
//
//func (mt MatchTeam) Bet(null string) string {
//	if bet := core.MatchTeam(mt).Bet(); bet != nil {
//		return strconv.Itoa(int(*bet))
//	}
//	return null
//}
//
//func (mt MatchTeam) Name(arg string) string {
//	if name := core.MatchTeam(mt).Name(); name == nil {
//		return arg
//	} else {
//		return *name
//	}
//}
//
//func (mt MatchTeam) Table(prefix string) (string, []tgbotapi.InlineKeyboardButton) {
//	core := core.MatchTeam(mt)
//	keys := make([]tgbotapi.InlineKeyboardButton, matchKeyCount)
//	text := mt.Name("не определено")
//	text += strings.Repeat(" ", matchTeamNameLen-len([]rune(text)))
//	if bet := core.Bet(); bet == nil {
//		text += strings.Repeat(" ", matchTeamGoalLen-len([]rune(matchTeamGoalNil))) + matchTeamGoalNil
//		for counter := 0; counter < matchKeyCount; counter++ {
//			value := strconv.Itoa(counter)
//			keys[counter] = tgbotapi.NewInlineKeyboardButtonData(value, prefix+strconv.Itoa(counter))
//		}
//	} else if *bet < matchKeyCount-matchOffsetKeyCount {
//		for counter := 0; counter < matchKeyCount; counter++ {
//			value := strconv.Itoa(counter)
//			if int(*bet) == counter {
//				text += strings.Repeat(" ", matchTeamGoalLen-len([]rune(value))) + value
//				value = matchSelectKeyLeft + value + matchSelectKeyRight
//			}
//			keys[counter] = tgbotapi.NewInlineKeyboardButtonData(value, prefix+strconv.Itoa(counter))
//		}
//	} else {
//		counter := 0
//		for ; counter < matchKeyCount-matchOffsetKeyCount<<2-1; counter++ {
//			value := strconv.Itoa(counter)
//			keys[counter] = tgbotapi.NewInlineKeyboardButtonData(value, prefix+strconv.Itoa(counter))
//		}
//		for ; counter < matchKeyCount; counter++ {
//			value := strconv.Itoa(int(*bet) - matchOffsetKeyCount + counter)
//			if int(*bet) == counter {
//				text += strings.Repeat(" ", matchTeamGoalLen-len([]rune(value))) + value
//				value = matchSelectKeyLeft + value + matchSelectKeyRight
//			}
//			keys[counter] = tgbotapi.NewInlineKeyboardButtonData(value, prefix+strconv.Itoa(counter))
//		}
//	}
//	return text, keys
//}
//
//type (
//	Match core.Match
//
//	MatchTeam core.MatchTeam
//)

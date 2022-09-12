package screens

import (
	"PJaK/core"
	"PJaK/views"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	statisticAllText      = "Общая"
	statisticCaption      = statisticCaptionEmoji + " " + statisticCaptionText
	statisticCaptionEmoji = "🗒"
	statisticCaptionText  = "Статистика"
	statisticLength       = 23
)

func newLoadStatisticBase(caller Interface) loadStatisticBase {
	return loadStatisticBase{caller}
}

func newLoadStatisticRound(caller Interface, core core.Round) loadStatisticRound {
	return loadStatisticRound{loadStatisticBase: newLoadStatisticBase(caller), core: core}
}

func newLoadStatisticTournament(caller Interface, core core.Tournament) loadStatisticTournament {
	return loadStatisticTournament{loadStatisticBase: newLoadStatisticBase(caller), core: core}
}

func newLoadStatisticTournamentAll(caller Interface, core core.Tournament) loadStatisticTournamentAll {
	return loadStatisticTournamentAll{loadStatisticBase: newLoadStatisticBase(caller), core: core}
}

func newStatisticTournament(base Base, caller, loader Interface,
	core core.Tournament, rounds []core.Round) statisticTournament {
	roundMap := initTournament(MapRound{}, rounds)
	return statisticTournament{Back: NewBack(base, caller, loader), core: core, roundMap: roundMap, roundSlice: rounds}
}

func statisticMainManagerTournament(main Main, core core.Tournament) *Loading {
	executor := loadMain(main.Base, main.manager, statisticCaption)
	return NewLoading(main.Base, executor, newLoadStatisticTournament(executor, core), loadTournamentText)
}

func (loadStatisticBase) Caption() string { return statisticCaption }

func (lsr loadStatisticRound) Do(action *Loading) Event {
	if matches, fail := core.GetMatches(lsr.core, action.user); fail == nil {
		var wg sync.WaitGroup
		statisticRoundBets := make([]StatisticRoundBets, len(matches))
		wg.Add(len(matches))
		for index, item := range matches {
			go func(match core.Match, index int) {
				bets, fail0 := core.GetBets(match, action.user)
				if fail0 == nil {
					statisticRoundBets[index].Data = bets
					statisticRoundBets[index].Match = match
				}
				wg.Done()
			}(item, index)
		}
		wg.Wait()
		texts := make([]string, len(matches))
		text := "<pre>" +
			"    П р о г н о з ы    |  О ч к и\n" +
			" Perke | Jagry |Kardick| P | J | K\n" +
			"===================================\n"
		pointsMap := map[int8]int16{0: 0, 1: 0, 2: 0}

		for index, item := range statisticRoundBets {
			teamStrings := ""
			teamString1 := views.MatchTeam(item.Match.Team1).Full("---")
			teamString2 := views.MatchTeam(item.Match.Team2).Full("---")
			teamRunes1, teamRunes2 := []rune(teamString1), []rune(teamString2)
			if len(teamRunes1)+len(teamRunes2) < statisticLength {
				teamStrings = teamString1 + ":" + teamString2
			} else {
				teamStringShort1 := views.MatchTeam(item.Match.Team1).Short("---")
				teamStringShort2 := views.MatchTeam(item.Match.Team2).Short("---")
				teamStrings = teamStringShort1 + ":" + teamStringShort2
			}
			teamRunes := []rune(teamStrings)
			result := views.Match(item.Match).Result("")
			text := teamStrings + strings.Repeat(" ", statisticLength-len(teamRunes)) +
				"|" + strings.Repeat(" ", 10-len([]rune(result))) + result + "\n"
			betsText := make([]string, 3)
			pointsText := make([]string, 3)
			// TODO : отойти от четких ID. Например, чтоб раьотало с юзерами 1,2 и 3, а не про пеорядку
			timeNow := time.Now()
			for counter := int8(0); counter < 3; counter++ {
				switch item.Data[counter].Goals.(type) {
				case core.BetGoalsHidden:
					if item.Time().Sub(timeNow) > 0 {
						if item.Data[counter].Goals.(core.BetGoalsHidden) {
							betsText[counter] = "   +   "
						} else {
							betsText[counter] = "   -   "
						}
					} else {
						betsText[counter] = "  -:-  "
					}
				case core.BetGoalsReal:
					betGoals := item.Data[counter].Goals.(core.BetGoalsReal)
					if betGoals.Bet1 == nil || betGoals.Bet2 == nil {
						betsText[counter] = "  -:-  "
					} else {
						bet1, bet2 := strconv.Itoa(int(*betGoals.Bet1)), strconv.Itoa(int(*betGoals.Bet2))
						betsText[counter] = strings.Repeat(" ", 3-len([]rune(bet1))) +
							bet1 + ":" + bet2 + strings.Repeat(" ", 3-len([]rune(bet2)))
					}
				}
				if item.Data[counter].Points == nil {
					pointsText[counter] = "   "
				} else {
					pointsMap[counter] += int16(*item.Data[counter].Points)
					points := strconv.Itoa(int(*item.Data[counter].Points))
					pointsText[counter] = strings.Repeat(" ", 3-len([]rune(points))) + points
				}
			}
			texts[index] = text + strings.Join(betsText, "|") + "|" + strings.Join(pointsText, "|") + "\n"
		}
		text += strings.Join(texts, "-------+-------+-------+---+---+---\n")
		text += "===================================\n"
		text += "                       |"
		pointsSlice := []string{}
		for counter := int8(0); counter < 3; counter++ {
			points := strconv.Itoa(int(pointsMap[counter]))
			points = strings.Repeat(" ", 3-len([]rune(points))) + points
			pointsSlice = append(pointsSlice, points)
		}
		text += strings.Join(pointsSlice, "|") + "</pre>"
		return NewEvent(lsr.caller, text)
	}
	return NewEvent(NewError(action.Base, action, "!!!", "!!!"), "")
}

func (lst loadStatisticTournament) Do(action *Loading) Event {
	if rounds, fail := core.GetTournamentRounds(lst.core); fail == nil {
		return NewEvent(newStatisticTournament(action.Base, lst.caller, action, lst.core, rounds), "")
	}
	return NewEvent(NewError(action.Base, action, "!!!", "!!!"), "")
}

func (lst loadStatisticTournamentAll) Do(action *Loading) Event {
	if rounds, fail := core.GetTournamentRounds(lst.core); fail == nil {
		return NewEvent(newStatisticTournament(action.Base, lst.caller, action, lst.core, rounds), "")
	}
	return NewEvent(NewError(action.Base, action, "!!!", "!!!"), "")
}

func (st statisticTournament) Handle(id int, update tgbotapi.Update) (Interface, bool, tgbotapi.Chattable) {
	if cq := update.CallbackQuery; cq != nil && cq.Message != nil && cq.Message.MessageID == id {
		if cq.Data == betsTournamentIdPrefix+"all" {
			action := NewLoading(st.Base, st.loader, newLoadStatisticTournamentAll(st.loader, st.core), "!!!")
			return action, false, tgbotapi.NewCallback(cq.ID, "Общая")
		} else if strings.HasPrefix(cq.Data, betsTournamentIdPrefix) {
			dataId, fail := strconv.ParseUint(cq.Data[len(betsTournamentIdPrefix):], 36, 16)
			if fail != nil {
				panic("screens.BetsTournament.Handle:" + fail.Error())
			}
			if round, ok := st.roundMap[core.RoundId(dataId)]; ok {
				action := NewLoading(st.Base, st.loader, newLoadStatisticRound(st.loader, round), "!!!")
				return action, false, tgbotapi.NewCallback(cq.ID, round.Name)
			} else {
				return nil, false, tgbotapi.NewCallback(cq.ID, "")
			}
		}
	}
	return st.Back.Handle(id, update)
}

func (st statisticTournament) Out() *InterfaceOut {
	if len(st.roundMap) == 0 {
		text := views.Tournament(st.core).Caption(statisticCaption)
		return &InterfaceOut{Keyboard: backKeyboard, Text: NewView(text, betsTournamentEmptyText).Text()}
	}
	keyboard := [][]tgbotapi.InlineKeyboardButton{{tgbotapi.NewInlineKeyboardButtonData(statisticAllText, betsTournamentIdPrefix+"all")}}
	rounds := core.GetRounds(st.roundSlice)
	for counter := 0; counter < len(rounds); counter = counter + betsRowCount {
		row := make([]tgbotapi.InlineKeyboardButton, 0)
		bound := counter + betsRowCount
		empty := 0
		if bound > len(rounds) {
			empty = bound - len(rounds)
			bound = len(rounds)
		}
		for offset, item := range rounds[counter:bound] {
			id := betsTournamentIdPrefix + strconv.FormatUint(uint64(item.Id), 36)
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(counter+offset+1), id))
		}
		for empty > 0 {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(" ", betsTournamentIdPrefix+"0"))
			empty--
		}
		keyboard = append(keyboard, row)
	}
	text := NewView(views.Tournament(st.core).Caption(statisticCaption), betsRoundsText).Text()
	return &InterfaceOut{Keyboard: append(keyboard, backRow), Text: text}
}

type (
	loadStatisticRound struct {
		loadStatisticBase
		core core.Round
	}

	loadStatisticTournament struct {
		loadStatisticBase
		core core.Tournament
	}

	loadStatisticTournamentAll struct {
		loadStatisticBase
		core core.Tournament
	}

	loadStatisticBase struct {
		caller Interface
	}

	StatisticRoundBets struct {
		core.Match
		Data map[int8]core.Bets
	}

	statisticTournament struct {
		Back
		core       core.Tournament
		manager    matchManager
		roundMap   MapRound
		roundSlice []core.Round
	}
)

var statisticMainManager = MainManager{tournament: statisticMainManagerTournament}

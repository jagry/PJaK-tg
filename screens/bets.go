package screens

import (
	"PJaK/core"
	"PJaK/helpers"
	"errors"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const (
	betsCaption      = betsCaptionEmoji + " " + betsCaptionText
	betsCaptionText  = "Ставки"
	betsCaptionEmoji = "🎲"
	betsBackId       = betsIdPrefix + backIdConst
	// !!! betsCloseId      = betsIdPrefix + closeIdConst
	betsSaveId   = betsIdPrefix + saveIdConst
	betsRowCount = 7

	betsIdPrefix           = "bets."
	betsMatchIdPrefix      = betsIdPrefix + "match."
	betsMatchTeam1IdPrefix = betsMatchIdPrefix + "1."
	betsMatchTeam2IdPrefix = betsMatchIdPrefix + "2."
	betsRoundIdPrefix      = betsIdPrefix + "round."
	betsTournamentIdPrefix = betsIdPrefix + "tournament."

	betsPrefixSelected = "<"
	betsSuffixSelected = ">"

	betsBackText            = backTextConst
	betsEmptyText           = "Еще нет футбольных турниров в текущих сезонах"
	betsErrorText           = "Ошибка загрузки футбольных турниров"
	betsLoadTournamentsText = "Идет загрузка футбольных турниров" + loadingText
	betLoadRoundsText       = "Идет загрузка туров" + loadingText
	betsMatchesText         = "Матчи:"
	betsRoundsText          = "Туры:"
	betsRoundEmptyText      = "Расписание матчей будет определено позже"
	betsSaveText            = saveTextConst
	betsTournamentsText     = "Футбольные турниры:"
	betsTournamentEmptyText = "Расписание туров будет определено позже"

	betsNameUndefined  = "не определена"
	betsNameUndefined1 = "Команда 1"
	betsNameUndefined2 = "Команда 2"
)

func initBetsTournament(matchMap *BetsMapRound, matchSlice []*core.Round) {
	for _, match := range matchSlice {
		if len(match.Rounds) == 0 {
			(*matchMap)[match.Id] = match
		} else {
			initBetsTournament(matchMap, match.Rounds)
		}
	}
}

func LoadBetsMain(base Base) *Loading {
	return NewLoading(base, base, NewBetsMainFactory(), betsCaption, betsLoadTournamentsText)
}

func LoadBetsTournament(base Base, caller Interface, core *core.Tournament) *Loading {
	caption := betsCaption + dividerTextConst + core.Emoji + " " + core.Full
	return NewLoading(base, caller, NewBetsTournamentFactory(caller, core), caption, betsLoadTournamentsText)
}

func NewBetsMain(base Base, tournaments []*core.Tournament) BetsMain {
	tournamentMap := BetsMapTournament{}
	for _, tournament := range tournaments {
		tournamentMap[tournament.Id] = tournament
	}
	return BetsMain{Base: base, tournamentMap: tournamentMap, tournamentSlice: tournaments}
}

func NewBetsMainFactory() BetsMainFactory {
	return BetsMainFactory{}
}

func NewBetsMatch(b Base, i Interface, t *core.Tournament, r *core.Round, c *core.Match) BetsMatch {
	return BetsMatch{Base: b, caller: i, core: c, round: r, tournament: t}
}

func NewBetsMatchFactory(i Interface, t *core.Tournament, r *core.Round, c *core.Match) BetsMatchFactory {
	return BetsMatchFactory{caller: i, core: c, round: r, tournament: t}
}

func NewBetsRound(b Base, i Interface, t *core.Tournament, c *core.Round, m []*core.Match) BetsRound {
	matchMap := BetsMapMatch{}
	for _, match := range m {
		matchMap[match.Id] = match
	}
	return BetsRound{Base: b, caller: i, core: c, matchMap: matchMap, matchSlice: m, tournament: t}
}

func NewBetsRoundFactory(caller Interface, tournament *core.Tournament, core *core.Round) BetsRoundFactory {
	return BetsRoundFactory{caller: caller, core: core, tournament: tournament}
}

func NewBetsTournament(base Base, caller Interface, core *core.Tournament, rounds []*core.Round) BetsTournament {
	roundMap := BetsMapRound{}
	initBetsTournament(&roundMap, rounds)
	return BetsTournament{Base: base, caller: caller, core: core, roundMap: roundMap, roundSlice: rounds}
}

func NewBetsTournamentFactory(caller Interface, core *core.Tournament) BetsTournamentFactory {
	return BetsTournamentFactory{caller: caller, core: core}
}

func (bm BetsMain) Handle(update telegram.Update) (Interface, bool, telegram.Chattable) {
	if update.CallbackQuery != nil && strings.HasPrefix(update.CallbackQuery.Data, betsTournamentIdPrefix) {
		id, fail := strconv.ParseUint(update.CallbackQuery.Data[len(betsTournamentIdPrefix):], 10, 16)
		if fail != nil {
			panic("screens.BetsMain.Handle:" + fail.Error())
		}
		tournament := bm.tournamentMap[core.TournamentId(id)]
		callback := telegram.NewCallback(update.CallbackQuery.ID, tournament.Emoji+" "+tournament.Full)
		//!!!text := betsCaption + dividerTextConst + tournament.Emoji + " " + tournament.Full
		//factory := NewBetsTournamentFactory(bm, tournament)
		//return NewLoading(bm.Base, bm, factory, text, betLoadRoundsText), false, callback
		return LoadBetsTournament(bm.Base, bm, tournament), false, callback
	}
	return bm.Base.Handle(update)
}

func (bm BetsMain) Out() *InterfaceOut {
	if len(bm.tournamentMap) == 0 {
		return &InterfaceOut{Keyboard: [][]telegram.InlineKeyboardButton{{baseCloseButton}}, Text: betsEmptyText}
	}
	keyboard := make([][]telegram.InlineKeyboardButton, 0)
	for _, tournament := range bm.tournamentSlice {
		id := betsTournamentIdPrefix + strconv.FormatUint(uint64(tournament.Id), 10)
		button := telegram.NewInlineKeyboardButtonData(tournament.Emoji+" "+tournament.Full, id)
		keyboard = append(keyboard, []telegram.InlineKeyboardButton{button})
	}
	keyboard = append(keyboard, []telegram.InlineKeyboardButton{baseCloseButton})
	return &InterfaceOut{Keyboard: keyboard, Text: NewView(betsCaption, betsTournamentsText).Text()}
}

func (BetsMainFactory) Execute(base Base) Interface {
	if tournaments, fail := core.GetTournaments(); fail == nil {
		return NewBetsMain(base, tournaments)
	}
	return NewError(base, base, betsCaption, betsErrorText)
}

func (bm BetsMatch) Handle(update telegram.Update) (Interface, bool, telegram.Chattable) {
	if update.CallbackQuery != nil {
		if strings.HasPrefix(update.CallbackQuery.Data, betsMatchTeam1IdPrefix) {
			idString := update.CallbackQuery.Data[len(betsMatchTeam1IdPrefix):]
			idUint64, fail := strconv.ParseUint(idString, 10, 8)
			if fail != nil {
				panic("screens.BetsRound.Handle:" + fail.Error())
			}
			bm.core.Team1.SetBet(byte(idUint64))
			text := helpers.MatchTeam(bm.core.Team1).Name(betsNameUndefined1) + " " + idString
			return bm, false, telegram.NewCallback(update.CallbackQuery.ID, text)
		} else if update.CallbackQuery.Data == betsSaveId {
			return bm.caller, false, telegram.NewCallback(update.CallbackQuery.ID, betsBackText)
		} else if update.CallbackQuery.Data == betsBackId {
			return bm.caller, false, telegram.NewCallback(update.CallbackQuery.ID, betsBackText)
		}
	} else if update.Message != nil && !update.Message.IsCommand() {
		dmc := telegram.NewDeleteMessage(update.FromChat().ID, update.Message.MessageID)
		regexp := regexp.MustCompile("^(\\d+)-(\\d+)$")
		if regexpResult := regexp.FindStringSubmatch(update.Message.Text); len(regexpResult) == 3 {
			goal1, fail := strconv.ParseUint(regexpResult[1], 10, 8)
			if fail != nil {
				if errors.Is(fail, strconv.ErrRange) {
					return NewError(bm.Base, bm, baseUserCaptionText, "!!!! Число голов должно быть в диапазоне от 0 до 255 включительно"), false, dmc
				}
			}
			goal2, fail := strconv.ParseUint(regexpResult[2], 10, 8)
			if fail != nil {
				if errors.Is(fail, strconv.ErrRange) {
					return NewError(bm.Base, bm, baseUserCaptionText, "!!!! Число голов должно быть в диапазоне от 0 до 255 включительно"), false, dmc
				}
			}
			team1, team2 := &bm.core.Team1, &bm.core.Team2
			bet1, bet2 := team1.Bet(), team2.Bet()
			if bet1 == nil || bet2 == nil || *bet1 != byte(goal1) || *bet2 != byte(goal2) {
				team1.SetBet(byte(goal1))
				team2.SetBet(byte(goal2))
				return bm, false, dmc
			}
			return nil, false, dmc
		}
		return NewError(bm.Base, bm, baseUserCaptionText, "!!!! X-X"), false, dmc
	}
	return bm.Base.Handle(update)
}

func (bm BetsMatch) Out() *InterfaceOut {
	helper, keys, team1, team2 := helpers.MatchTeam(bm.core.Team1), make([][]telegram.InlineKeyboardButton, 4), "", ""
	caption := betsCaption + dividerTextConst + bm.tournament.Emoji + " "
	caption += bm.tournament.Full + dividerTextConst + bm.round.Name
	team1, keys[0] = helper.Table(betsNameUndefined1, betsMatchTeam1IdPrefix, betsPrefixSelected, betsSuffixSelected)
	helper = helpers.MatchTeam(bm.core.Team2)
	team2, keys[1] = helper.Table(betsNameUndefined2, betsMatchTeam1IdPrefix, betsPrefixSelected, betsSuffixSelected)
	keys[2] = []telegram.InlineKeyboardButton{telegram.NewInlineKeyboardButtonData(betsSaveText, saveIdConst)}
	keys[3] = betsBackCloseRow
	text := "<pre>" + team1 + "\n" + team2 + "</pre>"
	return &InterfaceOut{Keyboard: keys, Text: NewView(caption, text).Text()}
}

func (bmf BetsMatchFactory) Execute(base Base) Interface {
	if match, fail := core.GetMatch(bmf.core.Id, base.User()); fail == nil {
		return NewBetsMatch(base, bmf.caller, bmf.tournament, bmf.round, &match)
	}
	return NewError(base, base, betsCaption, betsErrorText)
}

func (br BetsRound) Handle(update telegram.Update) (Interface, bool, telegram.Chattable) {
	if update.CallbackQuery != nil {
		if strings.HasPrefix(update.CallbackQuery.Data, betsRoundIdPrefix) {
			id, fail := strconv.ParseUint(update.CallbackQuery.Data[len(betsRoundIdPrefix):], 36, 16)
			if fail != nil {
				panic("screens.BetsRound.Handle:" + fail.Error())
			}
			if match, ok := br.matchMap[core.MatchId(id)]; ok {
				text := helpers.Match(*match).Players(betsNameUndefined, ":")
				factory := NewBetsMatchFactory(br, br.tournament, br.core, match)
				//!!!factory := NewBetsMatchFactory(LoadBetsTournament(br.Base, br, br.tournament), br.tournament, br.core, match)
				loading := NewLoading(br.Base, br, factory, text, betLoadRoundsText)
				return loading, false, telegram.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			} else {
				return nil, false, telegram.NewCallback(update.CallbackQuery.ID, "")
			}
		} else if update.CallbackQuery.Data == betsBackId {
			return br.caller, false, telegram.NewCallback(update.CallbackQuery.ID, betsBackText)
		}
	}
	return br.Base.Handle(update)
}

func (br BetsRound) Out() *InterfaceOut {
	caption := betsCaption + dividerTextConst + br.tournament.Emoji
	caption += " " + br.tournament.Full + dividerTextConst + br.core.Name
	if len(br.matchMap) == 0 {
		return &InterfaceOut{Keyboard: betsBackCloseKeyboard, Text: NewView(caption, betsRoundEmptyText).Text()}
	}
	keyboard := make([][]telegram.InlineKeyboardButton, 0)
	for _, match := range br.matchSlice {
		id := betsRoundIdPrefix + strconv.FormatUint(uint64(match.Id), 36)
		helper := helpers.Match(*match)
		text := helper.Time() + ": " + helper.Players(betsNameUndefined, ":") + " " + helper.Bet("")
		keyboard = append(keyboard, []telegram.InlineKeyboardButton{telegram.NewInlineKeyboardButtonData(text, id)})
	}
	return &InterfaceOut{Keyboard: append(keyboard, betsBackCloseRow), Text: NewView(caption, betsMatchesText).Text()}
}

func (brf BetsRoundFactory) Execute(base Base) Interface {
	if matches, fail := core.GetMatches(brf.core, base.User()); fail == nil {
		return NewBetsRound(base, brf.caller, brf.tournament, brf.core, matches)
	}
	return NewError(base, base, betsCaption, betsErrorText)
}

func (bt BetsTournament) Handle(update telegram.Update) (Interface, bool, telegram.Chattable) {
	if update.CallbackQuery != nil {
		if strings.HasPrefix(update.CallbackQuery.Data, betsTournamentIdPrefix) {
			id, fail := strconv.ParseUint(update.CallbackQuery.Data[len(betsTournamentIdPrefix):], 36, 16)
			if fail != nil {
				log.Println("screens.BetsTournament.Handle:" + fail.Error())
				panic("screens.BetsTournament.Handle:" + fail.Error())
			}
			if round, ok := bt.roundMap[core.RoundId(id)]; ok {
				text := betsCaptionText + dividerTextConst + bt.core.Full + dividerTextConst + round.Name
				loading := NewLoading(bt.Base, bt, NewBetsRoundFactory(bt, bt.core, round), text, betLoadRoundsText)
				return loading, false, telegram.NewCallback(update.CallbackQuery.ID, round.Name)
			} else {
				return nil, false, telegram.NewCallback(update.CallbackQuery.ID, "")
			}
		} else if update.CallbackQuery.Data == betsBackId {
			return LoadBetsMain(bt.Base), false, telegram.NewCallback(update.CallbackQuery.ID, betsBackText)
		}
	}
	return bt.Base.Handle(update)
}

func (bt BetsTournament) Out() *InterfaceOut {
	if len(bt.roundMap) == 0 {
		keyboard := [][]telegram.InlineKeyboardButton{{betsBackButton, baseCloseButton}}
		text := NewView(betsCaption+bt.core.Emoji+" "+bt.core.Full, betsTournamentEmptyText).Text()
		return &InterfaceOut{Keyboard: keyboard, Text: text}
	}
	rounds := core.GetRounds(bt.roundSlice)
	keyboard := make([][]telegram.InlineKeyboardButton, 0)
	for counter := 0; counter < len(rounds); counter = counter + betsRowCount {
		row := make([]telegram.InlineKeyboardButton, 0)
		bound := counter + betsRowCount
		empty := 0
		if bound > len(rounds) {
			empty = bound - len(rounds)
			bound = len(rounds)
		}
		for offset, item := range rounds[counter:bound] {
			id := betsTournamentIdPrefix + strconv.FormatUint(uint64(item.Id), 36)
			row = append(row, telegram.NewInlineKeyboardButtonData(strconv.Itoa(counter+offset+1), id))
		}
		for empty > 0 {
			row = append(row, telegram.NewInlineKeyboardButtonData(" ", betsTournamentIdPrefix+"0"))
			empty--
		}
		keyboard = append(keyboard, row)
	}
	text := NewView(betsCaption+dividerTextConst+bt.core.Emoji+" "+bt.core.Full, betsRoundsText).Text()
	return &InterfaceOut{Keyboard: append(keyboard, betsBackCloseRow), Text: text}
}

func (btf BetsTournamentFactory) Execute(base Base) Interface {
	if rounds, fail := core.GetTournamentRounds(btf.core); fail == nil {
		return NewBetsTournament(base, btf.caller, btf.core, rounds)
	}
	return NewError(base, base, betsCaptionText+dividerTextConst+btf.core.Full, betsErrorText)
}

type (
	BetsMain struct {
		Base
		tournamentMap   BetsMapTournament
		tournamentSlice []*core.Tournament
	}

	BetsMainFactory struct{}

	BetsMapMatch map[core.MatchId]*core.Match

	BetsMapRound map[core.RoundId]*core.Round

	BetsMapTournament map[core.TournamentId]*core.Tournament

	BetsMatch struct {
		Base
		caller     Interface
		core       *core.Match
		round      *core.Round
		tournament *core.Tournament
	}

	BetsMatchFactory struct {
		caller     Interface
		core       *core.Match
		round      *core.Round
		tournament *core.Tournament
	}

	BetsRound struct {
		Base
		caller     Interface
		core       *core.Round
		matchMap   BetsMapMatch
		matchSlice []*core.Match
		tournament *core.Tournament
	}

	BetsRoundFactory struct {
		caller     Interface
		core       *core.Round
		matchMap   BetsMapMatch
		matchSlice []*core.Match
		tournament *core.Tournament
	}

	BetsTournament struct {
		Base
		caller     Interface
		core       *core.Tournament
		roundMap   BetsMapRound
		roundSlice []*core.Round
	}

	BetsTournamentFactory struct {
		caller Interface
		core   *core.Tournament
	}
)

var (
	betsBackButton = telegram.NewInlineKeyboardButtonData(betsBackText, betsBackId)
	betsSaveButton = telegram.NewInlineKeyboardButtonData(betsSaveText, betsSaveId)

	betsBackCloseKeyboard = [][]telegram.InlineKeyboardButton{betsBackCloseRow}

	betsBackCloseRow = []telegram.InlineKeyboardButton{betsBackButton, baseCloseButton}
)

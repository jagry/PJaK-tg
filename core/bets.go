package core

import (
	"log"
	"time"
)

const (
	getBets  = `SELECT "goal1", "goal2", "user" FROM "bet" WHERE "match" = $1`
	saveBets = `INSERT INTO "bet"("goal1", "goal2", "match", "user")
		VALUES($1, $2, $3, $4) ON CONFLICT ("match", "user") DO UPDATE SET "goal1" = $1, "goal2" = $2`

	NoneBetsState = 0
	HaveBetsState = 1
)

func GetBets(match Match, user int8) (map[int8]Bets, error) {
	// TODO: Сделать запрос данных матча
	rows, fail := db.Query(getBets, match.Id)
	if rows != nil {
		defer rows.Close()
	}
	if fail != nil {
		log.Println("core.GetBets: ", fail.Error())
		return nil, fail
	}
	result := make(map[int8]Bets, 0)
	for counter := int8(0); counter < 3; counter++ {
		//if match.time.Sub(time.Now()) < 0 {
		//	points := int8(-1)
		//	result[counter] = Bets{Goals: BetGoalsReal{}, Points: &points, User: counter}
		//} else if match.Team1.bets != nil && match.Team2.bets != nil {
		//	result[counter] = Bets{Goals: BetGoalsReal{}, User: counter}
		//} else if counter == user {
		//	result[counter] = Bets{Goals: BetGoalsReal{}, User: counter}
		//} else {
		result[counter] = Bets{Goals: BetGoalsReal{}, User: counter}
		//}
	}
	var betsUser int8
	betGoals := BetGoalsReal{}
	for rows.Next() {
		fail = rows.Scan(&betGoals.Bet1, &betGoals.Bet2, &betsUser)
		if fail != nil {
			return nil, fail
		}
		if betsUser == user {
			match.Team1.bets = betGoals.Bet1
			match.Team2.bets = betGoals.Bet2
		}
		result[betsUser] = Bets{Goals: BetGoalsReal{Bet1: betGoals.Bet1, Bet2: betGoals.Bet2}, User: betsUser}
	}
	timeNow := time.Now()
	for counter := int8(0); counter < 3; counter++ {
		bets := result[counter]
		if match.Team1.result == nil || match.Team2.result == nil {
			if match.time.Sub(timeNow) < 0 {
				if bets.Goals.(BetGoalsReal).Bet1 == nil || bets.Goals.(BetGoalsReal).Bet2 == nil {
					points := int8(-1)
					bets.Points = &points
				}
			} else if match.Team1.bets == nil || match.Team2.bets == nil {
				if bets.Goals.(BetGoalsReal).Bet1 == nil || bets.Goals.(BetGoalsReal).Bet2 == nil {
					bets.Goals = BetGoalsHidden(false)
				} else {
					bets.Goals = BetGoalsHidden(true)
				}
			}
		} else {
			if match.time.Sub(timeNow) < 0 {
				betsGoals := bets.Goals.(BetGoalsReal)
				if betsGoals.Bet1 == nil || betsGoals.Bet2 == nil {
					points := int8(-1)
					bets.Points = &points
				} else if *betsGoals.Bet1 == *match.Team1.result && *betsGoals.Bet2 == *match.Team2.result {
					points := int8(7)
					bets.Points = &points
				} else if *betsGoals.Bet1-*betsGoals.Bet2 == *match.Team1.result-*match.Team2.result {
					points := int8(5)
					bets.Points = &points
				} else if *betsGoals.Bet1 < *betsGoals.Bet2 && *match.Team1.result < *match.Team2.result {
					points := int8(4)
					bets.Points = &points
				} else if *betsGoals.Bet1 > *betsGoals.Bet2 && *match.Team1.result > *match.Team2.result {
					points := int8(4)
					bets.Points = &points
				} else {
					points := int8(0)
					bets.Points = &points
				}
			}
		}
		result[counter] = bets
	}
	return result, nil
}

func SaveBets(bet1, bet2 byte, match MatchId, user int8) error {
	_, fail := db.Exec(saveBets, bet1, bet2, match, user)
	return fail
}

func (bgh BetGoalsHidden) Set(goal1, goal2 byte) {

}

func (bgh BetGoalsReal) Set(goal1, goal2 byte) {

}

type (
	Bets struct {
		User   int8
		Goals  interface{}
		Points *int8
	}

	BetGoals interface{}

	BetGoalsHidden bool

	BetGoalsReal struct{ Bet1, Bet2 *byte }
)

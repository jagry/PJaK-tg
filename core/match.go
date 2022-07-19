package core

import (
	"log"
	"time"
)

const (
	get = `SELECT "m"."id", "b"."goal1", "1"."full", "m"."goal1", "1"."short", "b"."goal2",
		"2"."full", "m"."goal2", "2"."short", "m"."time" FROM "match" "m"
    	JOIN "team" "1" ON "1"."id" = "m"."team1"
		JOIN "team" "2" ON "2"."id" = "m"."team2"
		LEFT JOIN "bet" "b" ON "b"."match" = "m"."id" AND "b"."user" = $1
		WHERE `

	getMatch = get + `"m"."id" = $2`

	getMatches = get + `"m"."round" = $2 ORDER BY "m"."time"`
)

func GetMatch(matchId MatchId, user int8) (match Match, fail error) {
	vars := []interface{}{&match.Id, &match.Team1.bets, &match.Team1.full, &match.Team1.result, &match.Team1.short}
	vars = append(vars, &match.Team2.bets, &match.Team2.full, &match.Team2.result, &match.Team2.short, &match.time)
	fail = db.QueryRow(getMatch, user, matchId).Scan(vars...)
	return
}

func GetMatches(round *Round, user int8) ([]*Match, error) {
	result := make([]*Match, 0)
	rows, fail := db.Query(getMatches, user, round.Id)
	if rows != nil {
		defer rows.Close()
	}
	if fail != nil {
		log.Println("core.GetMatches:", fail.Error())
		return nil, fail
	}
	for rows.Next() {
		match := &Match{}
		vars := []interface{}{&match.Id, &match.Team1.bets, &match.Team1.full, &match.Team1.result, &match.Team1.short}
		vars = append(vars, &match.Team2.bets, &match.Team2.full, &match.Team2.result, &match.Team2.short, &match.time)
		fail = rows.Scan(vars...)
		if fail != nil {
			return nil, fail
		}
		result = append(result, match)
	}
	return result, nil
}

func (match Match) Time() time.Time {
	return match.time
}

func (mt MatchTeam) Bet() *byte {
	return mt.bets
}

func (mt MatchTeam) Full() *string {
	return mt.full
}

func (mt MatchTeam) Result() *byte {
	return mt.result
}

func (mt *MatchTeam) SetBet(arg byte) {
	mt.bets = &arg
}

func (mt MatchTeam) Short() *string {
	return mt.short
}

type (
	Match struct {
		Id    MatchId
		Team1 MatchTeam
		Team2 MatchTeam
		time  time.Time
	}

	MatchId uint16

	MatchTeam struct {
		bets, result *byte
		short, full  *string
	}
)

package core

import (
	"log"
	"time"
)

const (
	get = `SELECT "m"."id", "b"."goal1", "1"."name", "m"."goal1", "b"."goal2", "2"."name", "m"."goal2", "m"."time" FROM "match" "m"
    	JOIN "team" "1" ON "1"."id" = "m"."team1"
		JOIN "team" "2" ON "2"."id" = "m"."team2"
		LEFT JOIN "bet" "b" ON "b"."match" = "m"."id" AND "b"."user" = $1
		WHERE `

	getMatch = get + `"m"."id" = $2`

	getMatches = get + `"m"."round" = $2 ORDER BY "m"."time"`
)

func GetMatch(matchId MatchId, user int8) (match Match, fail error) {
	vars := []interface{}{&match.Id, &match.Team1.bets, &match.Team1.name, &match.Team1.result}
	vars = append(vars, &match.Team2.bets, &match.Team2.name, &match.Team2.result, &match.time)
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
		vars := []interface{}{&match.Id, &match.Team1.bets, &match.Team1.name, &match.Team1.result}
		vars = append(vars, &match.Team2.bets, &match.Team2.name, &match.Team2.result, &match.time)
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

func (mt MatchTeam) Name() *string {
	return mt.name
}

func (mt MatchTeam) Bet() *byte {
	return mt.bets
}

func (mt *MatchTeam) SetBet(arg byte) {
	mt.bets = &arg
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
		bets   *byte
		name   *string
		result *byte
	}
)

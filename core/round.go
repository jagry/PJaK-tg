package core

import "log"

const (
	getChildRounds = `SELECT "r"."id", "r"."name" FROM "round" "r"
		JOIN "roundLink" "l" ON "l"."child" = "r"."id" WHERE "l"."parent" = $1`

	getTournamentRounds = `SELECT "r"."id", "r"."name" FROM "round" "r"
		JOIN "seasonRound" "s" ON "r"."id" = "s"."round" WHERE "s"."ts" = $1`
)

func GetRounds(rounds []*Round) (result []*Round) {
	for _, round := range rounds {
		if len(round.Rounds) == 0 {
			result = append(result, round)
		} else {
			result = append(result, GetRounds(round.Rounds)...)
		}
	}
	return
}

func GetRoundsPrefix(prefix string, rounds []*Round) (result []*Round) {
	for _, round := range rounds {
		if len(round.Rounds) == 0 {
			round.Name = prefix + " " + round.Name
			result = append(result, round)
		} else {
			result = append(result, GetRoundsPrefix(round.Name+":", round.Rounds)...)
		}
	}
	return
}

func getRounds(query string, id uint16) ([]*Round, error) {
	rows, fail := db.Query(query, id)
	if rows != nil {
		defer rows.Close()
	}
	if fail != nil {
		log.Println("core.getRounds:", fail.Error())
		return nil, fail
	}
	result := make([]*Round, 0)
	for rows.Next() {
		round := &Round{}
		fail = rows.Scan(&round.Id, &round.Name)
		if fail != nil {
			return nil, fail
		}
		round.Rounds, fail = getRounds(getChildRounds, uint16(round.Id))
		result = append(result, round)
	}
	return result, nil
}

func GetTournamentRounds(tournament *Tournament) ([]*Round, error) {
	return getRounds(getTournamentRounds, uint16(tournament.Id))
}

type (
	Round struct {
		Id         RoundId
		Name       string
		Rounds     []*Round
		Tournament Tournament
	}

	RoundId uint16
)

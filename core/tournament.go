package core

import "log"

const getActiveTournaments = `SELECT "t"."emoji", "t"."full", "ts"."handler", "ts"."id" FROM "season" "s"
        JOIN "tournamentSeason" "ts" ON "s"."id" = "ts"."season"
        JOIN "tournament" "t" ON "t"."id" = "ts"."tournament" WHERE "s"."previous" IS NULL ORDER BY "ts"."index"`

func GetTournaments() ([]Tournament, error) {
	result := make([]Tournament, 0)
	rows, fail := db.Query(getActiveTournaments)
	if rows != nil {
		defer rows.Close()
	}
	if fail != nil {
		log.Println("core.GetTournaments:", fail.Error())
		return nil, fail
	}
	tournament := Tournament{}
	for rows.Next() {
		fail = rows.Scan(&tournament.Emoji, &tournament.Full, &tournament.Handler, &tournament.Id)
		if fail != nil {
			return nil, fail
		}
		result = append(result, tournament)
	}
	return result, nil
}

type (
	Tournament struct {
		Emoji   string
		Full    string
		Handler string
		Id      TournamentId
	}

	TournamentId uint16
)

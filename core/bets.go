package core

const (
	saveBets = `INSERT INTO "bet"("goal1", "goal2", "match", "user")
		VALUES($1, $2, $3, $4) ON CONFLICT ("match", "user") DO UPDATE SET "goal1" = $1, "goal2" = $2`
)

func SaveBets(bet1, bet2 byte, match MatchId, user int8) error {
	_, fail := db.Exec(saveBets, bet1, bet2, match, user)
	return fail
}

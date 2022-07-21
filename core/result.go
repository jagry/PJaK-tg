package core

const (
	saveResult = `UPDATE "match" SET "goal1" = $1, "goal2" = $2 WHERE "id" = $3`
)

func SaveResult(result1, result2 byte, match MatchId, user int8) error {
	_, fail := db.Exec(saveResult, result1, result2, match)
	return fail
}

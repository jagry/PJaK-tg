package views

import (
	"strconv"
)

func Goals(goal1, goal2 *byte) (result string) {
	if goal1 == nil {
		result = matchTeamGoalNil
	} else {
		result = strconv.Itoa(int(*goal1))
	}
	result += ":"
	if goal2 == nil {
		result += matchTeamGoalNil
	} else {
		result += strconv.Itoa(int(*goal2))
	}
	return
}

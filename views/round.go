package views

import "PJaK/core"

func (round Round) Caption(section string, tournament Tournament) string {
	return tournament.Caption(section) + formCaptionDelimiter + round.Name
}

type Round core.Round

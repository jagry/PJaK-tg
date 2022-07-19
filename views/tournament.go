package views

import "PJaK/core"

func (tournament Tournament) Caption(section string) string {
	return section + formCaptionDelimiter + tournament.Name()
}

func (tournament Tournament) Name() string {
	return tournament.Emoji + " " + tournament.Full
}

type Tournament core.Tournament

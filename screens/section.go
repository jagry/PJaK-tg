package screens

import (
	"PJaK/core"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewSection(
	n string, b SectionButton, m SectionModify, sv SectionSave, sn SectionScreen, t SectionTeamModify) Section {
	return Section{button: b, modify: m, name: n, save: sv, screen: sn, teamModify: t}
}

type (
	Section struct {
		button     SectionButton
		modify     SectionModify
		name       string
		save       SectionSave
		screen     SectionScreen
		teamModify SectionTeamModify
	}

	SectionButton func(core.Match) string

	SectionModify func(core.Match) (bool, *byte, *byte)

	SectionSave func(core.Match, int8) error

	SectionScreen func(core.Match) (string, [][]tgbotapi.InlineKeyboardButton)

	SectionTeamModify func(*core.MatchTeam, byte)
)

package screens

const (
	betsCaption      = betsCaptionEmoji + " " + betsCaptionText
	betsCaptionText  = "Прогнозы"
	betsCaptionEmoji = "🎲"
)

const (
	betsSaveId   = saveIdConst
	betsRowCount = 7

	betsMatchIdPrefix      = "match."
	betsMatchTeam1IdPrefix = betsMatchIdPrefix + "1."
	betsMatchTeam2IdPrefix = betsMatchIdPrefix + "2."
	betsTournamentIdPrefix = "tournament."

	betsPrefixSelected = ""
	betsSuffixSelected = ""

	betsEmptyText           = "Еще нет футбольных турниров в текущих сезонах"
	betsLoadMatchesText     = "Идет загрузка матчей" + loadingTextSuffix
	betsMatchesText         = "Матчи:"
	betsRoundsText          = "Туры:"
	betsRoundEmptyText      = "Расписание матчей будет определено позже"
	betsSaveText            = saveTextConst
	betsTournamentsText     = "Футбольные турниры:"
	betsTournamentEmptyText = "Расписание туров будет определено позже"
)

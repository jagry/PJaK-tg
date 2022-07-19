-- 0
CREATE TABLE "round"
(
    "data" varchar(9),
    "id"   smallserial PRIMARY KEY,
    "name" varchar(99) NOT NULL
);

CREATE TABLE "tournament"
(
    "emoji"  varchar(9),
    "id"     smallserial PRIMARY KEY,
    "full"   varchar(99),
    "sign"   varchar(9),
    "short"  varchar(9),
    "status" smallint
);

CREATE TABLE "tournamentSeasonHandler"
(
    "id"   smallserial PRIMARY KEY,
    "sign" varchar(9)
);

CREATE TABLE "type"
(
    "id"   smallserial PRIMARY KEY,
    "sign" varchar(9)
);

CREATE TABLE "user"
(
    "emoji"    varchar(9)  NOT NULL,
    "id"       smallserial PRIMARY KEY,
    "name"     varchar(9)  NOT NULL,
    "sign"     varchar(9)  NOT NULL,
    "telegram" varchar(99) NOT NULL
);
-- 1
CREATE TABLE "roundLink"
(
    "child"  smallint NOT NULL REFERENCES "round",
    "index"  smallint NOT NULL,
    "parent" smallint NOT NULL REFERENCES "round",
    CONSTRAINT "roundLinkP" PRIMARY KEY ("parent", "child"),
    CONSTRAINT "roundLinkU" UNIQUE ("index", "parent")
);

CREATE TABLE "season"
(
    "finish"   smallint NOT NULL,
    "id"       smallserial PRIMARY KEY,
    "previous" smallint,
    "start"    smallint NOT NULL,
    "type"     smallint REFERENCES "type",
    CONSTRAINT "seasonUPrevious" UNIQUE ("previous"),
    CONSTRAINT "seasonUStartFinish" UNIQUE ("start", "finish")
);
-- 2
CREATE TABLE "tournamentSeason"
(
    "handler"    smallint NOT NULL REFERENCES "tournamentSeasonHandler",
    "id"         smallserial PRIMARY KEY,
    "index"      smallint NOT NULL,
    --"name"       varchar(99) NOT NULL,
    "season"     smallint NOT NULL REFERENCES "season",
    "tournament" smallint NOT NULL REFERENCES "tournament",
    CONSTRAINT "tournamentSeasonUSeasonTournament" UNIQUE ("season", "tournament"),
    CONSTRAINT "tournamentSeasonUIndexSeason" UNIQUE ("index", "season")
);
-- 3
CREATE TABLE "seasonRound"
(
    "index" smallint NOT NULL,
    "round" smallint NOT NULL REFERENCES "round",
    "ts"    smallint NOT NULL REFERENCES "tournamentSeason",
    CONSTRAINT "seasonRoundP" PRIMARY KEY ("round", "ts"),
    CONSTRAINT "seasonRoundU" UNIQUE ("index", "ts")
);

CREATE TABLE "team"
(
    "full" varchar(99) NOT NULL,
    "id"   smallserial PRIMARY KEY,
    "short" varchar(99) NOT NULL,
    "ts"   smallint    NOT NULL REFERENCES "tournamentSeason",
    CONSTRAINT "teamUShort" UNIQUE ("short", "ts"),
    CONSTRAINT "teamUFull" UNIQUE ("full", "ts")
);
-- 4
CREATE TABLE "match"
(
    "goal1" smallint,
    "goal2" smallint,
    "id"    smallserial PRIMARY KEY,
    "round" smallint  NOT NULL REFERENCES "round",
    "team1" smallint REFERENCES "team",
    "team2" smallint REFERENCES "team",
    "time"  timestamp NOT NULL
);
-- 5
CREATE TABLE "bet"
(
    "goal1" smallint,
    "goal2" smallint,
    "match" smallint NOT NULL REFERENCES "match",
    "user"  smallint NOT NULL REFERENCES "user",
    CONSTRAINT "betP" PRIMARY KEY ("match", "user")
);

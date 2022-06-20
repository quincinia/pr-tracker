// Package models contains struct definitions for the database tables.
package models

// Could probably use the sql.NullXxx types, but keeping it simple for now
// Used when scanning
type NullTourneyType struct {
	TypeID *int
	Name   *string
	Valid  bool
}


type NullTier struct {
	TierID     *int
	Name       *string
	Multiplier *int
	Valid      bool
}

// Challonge or smash.gg
type TourneyType struct {
	TypeID int
	Name   string
}

// See types.sql
var (
	CHALLONGE = &TourneyType{1, "Challonge"}
	STARTGG   = &TourneyType{2, "start.gg"}
)

// C	B	A	S
// 75	150	200	300
type Tier struct {
	TierID     int
	Name       string
	Multiplier int
}

// See tiers.sql
var (
	C_TIER = &Tier{1, "C", 75}
	B_TIER = &Tier{2, "B", 150}
	A_TIER = &Tier{3, "A", 200}
	S_TIER = &Tier{4, "S", 300}
)

type Tournament struct {
	TourneyID      int
	Type           *TourneyType
	Name           string
	URL            string
	NumEntrants    int
	UniquePlacings int
	BracketReset   bool
	Tier           *Tier
}

type Player struct {
	PlayerID int
	Name     string
}

type Attendee struct {
	AttendeeID int
	Tourney    int
	Player     *Player
	Name       string
	Standing   int
}

type Attendees []Attendee

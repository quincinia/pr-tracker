// Package models contains struct definitions for the database tables.
package models

// Could probably use the sql.NullXxx types, but keeping it simple for now
// Used when scanning
type NullTourneyType struct {
	TypeID *int
	Name   *string
}

func (ntt *NullTourneyType) ToTourneyType() (tt *TourneyType) {
	if ntt.TypeID != nil && ntt.Name != nil {
		tt = &TourneyType{TypeID: *ntt.TypeID, Name: *ntt.Name}
	}
	return
}

type NullTier struct {
	TierID     *int
	Name       *string
	Multiplier *int
}

func (nt *NullTier) ToTier() (t *Tier) {
	if nt.TierID != nil && nt.Name != nil && nt.Multiplier != nil {
		t = &Tier{TierID: *nt.TierID, Name: *nt.Name, Multiplier: *nt.Multiplier}
	}
	return
}

type NullPlayer struct {
	PlayerID *int
	Name     *string
}

func (np *NullPlayer) ToPlayer() (p *Player) {
	if np.PlayerID != nil && np.Name != nil {
		p = &Player{PlayerID: *np.PlayerID, Name: *np.Name}
	}
	return
}

// Challonge or smash.gg
type TourneyType struct {
	TypeID int    `json:"typeID"`
	Name   string `json:"name"`
}

// See types.sql
var (
	CHALLONGE = &TourneyType{1, "Challonge"}
	STARTGG   = &TourneyType{2, "start.gg"}
)

// C	B	A	S
// 75	150	200	300
type Tier struct {
	TierID     int    `json:"tierID"`
	Name       string `json:"name"`
	Multiplier int    `json:"multiplier"`
}

// See tiers.sql
var (
	C_TIER = &Tier{1, "C", 75}
	B_TIER = &Tier{2, "B", 150}
	A_TIER = &Tier{3, "A", 200}
	S_TIER = &Tier{4, "S", 300}
)

type Tournament struct {
	TourneyID      int          `json:"tourneyID"`
	Type           *TourneyType `json:"type"`
	Name           string       `json:"name"`
	URL            string       `json:"url"`
	NumEntrants    int          `json:"numEntrants"`
	UniquePlacings int          `json:"uniquePlacings"`
	BracketReset   bool         `json:"bracketReset"`
	Tier           *Tier        `json:"tier"`
}

type Player struct {
	PlayerID int    `json:"playerID"`
	Name     string `json:"name"`
}

type Attendee struct {
	AttendeeID int     `json:"attendeeID"`
	Tourney    int     `json:"tourney"`
	Player     *Player `json:"player"`
	Name       string  `json:"name"`
	Standing   int     `json:"standing"`
}

type PlayerAttendance struct {
	Player
	Attendance []struct {
		AttendeeID int    `json:"attendeeID"`
		Tourney    int    `json:"tourney"`
		Name       string `json:"name"`
		Standing   int    `json:"standing"`
	} `json:"attendance"`
}

// Keeping this separate rather than putting it under Tournament
// Deprecated
type Attendees []Attendee

// Used when making HTTP requests
type FullTournament struct {
	Tournament
	Attendees []Attendee `json:"attendees"`
}

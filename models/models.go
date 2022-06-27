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

// Keeping this separate rather than putting it under Tournament
type Attendees []Attendee

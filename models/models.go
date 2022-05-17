// Package models contains struct definitions for the database tables.
package models

type TourneyType struct {
	TypeID int
	Name   string
}

type Tier struct {
	TierID     int
	Name       string
	Multiplier int
}

type Tournament struct {
	TourneyID      int
	Type           TourneyType
	Name           string
	URL            string
	NumEntrants    int
	UniquePlacings int
	BracketReset   bool
	Tier           Tier
}

type Player struct {
	PlayerID int
	Name     string
}

type Attendee struct {
	AttendeeID int
	Tourney    Tournament
	Player     Player
	Name       string
	Standing   int
}

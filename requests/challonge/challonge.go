// Package challonge contains functions for handling Challonge-based tournaments.
// Package assumes tournaments are complete.
package challonge

import (
	"encoding/json"
	"io"
)

// The JSON response puts a wrapper around the tournament object.
type Challonge struct {
	Tournament struct {
		Name           string
		URL            string `json:"full_challonge_url"`
		NumEntrants    int    `json:"participants_count"`
		UniquePlacings int    `json:"-"`
		BracketReset   bool   `json:"-"`
		Participants   []Participant
		Matches        []Match
	}
}

type Participant struct {
	Participant struct {
		ID       int
		Name     string
		Standing int `json:"final_rank"`
	}
}

type Match struct {
	Match struct {
		Player1ID int `json:"player1_id"`
		Player2ID int `json:"player2_id"`
		WinnerID  int `json:"winner_id"`
		LoserID   int `json:"loser_id"`
		Order     int `json:"suggested_play_order"`
	}
}

// Create a new tournament from a JSON source (generally an HTTP response)
func NewChallonge(r io.Reader) (c *Challonge, err error) {
	c = &Challonge{}
	decoder := json.NewDecoder(r)
	err = decoder.Decode(c)
	if err != nil {
		return nil, err
	}
	return
}

// In JavaScript: Math.trunc(Math.log2(n))
// Do not use with negative numbers.
// See: https://stackoverflow.com/questions/19339594/truncated-binary-logarithm
func TruncLog2(n int) (exp, value int) {
	if n <= 0 {
		return
	}
	value = 1
	n = n >> 1
	for n != 0 {
		exp++
		value *= 2
		n = n >> 1
	}
	return
}

func CalculatePlacings(numEntrants int) (placings int) {
	// Annoying edge cases
	if numEntrants < 1 {
		return -1
	}
	if numEntrants <= 4 {
		return numEntrants
	}

	exp, value := TruncLog2(numEntrants)
	placings = 2*exp + 1
	if numEntrants > (3*value)/2 {
		placings++
	}
	return
}

// Pass match by reference?
func isPresent(p Participant, m Match) bool {
	match := m.Match
	participant := p.Participant
	return participant.ID == match.Player1ID || participant.ID == match.Player2ID
}

// Returns the index of the first player with the given standing.
func (c *Challonge) FindStanding(standing int) (index int) {
	for i, p := range c.Tournament.Participants {
		if p.Participant.Standing == standing {
			return i
		}
	}
	return -1
}

// Similar to FindStanding, except with matches.
func (c *Challonge) FindMatch(order int) (index int) {
	for i, m := range c.Tournament.Matches {
		if m.Match.Order == order {
			return i
		}
	}
	return -1
}

// Second place will receive bonus points if they made a bracket reset.
// Bracket reset occurred if:
//  1. The last two matches occurred between first and second place.
//  2. The winners of the last two games are different.
// Bracket reset points will not be awarded if the first-place player won from losers. 
func (c *Challonge) ApplyResetPoints() {
	participants := c.Tournament.Participants
	matches := c.Tournament.Matches

	firstPlace := participants[c.FindStanding(1)]
	secondPlace := participants[c.FindStanding(2)]
	grands := matches[c.FindMatch(len(matches))]
	finals := matches[c.FindMatch(len(matches)-1)]

	c.Tournament.BracketReset =
		isPresent(firstPlace, finals) &&
			isPresent(secondPlace, finals) &&
			isPresent(firstPlace, grands) &&
			isPresent(secondPlace, grands) &&
			finals.Match.WinnerID != grands.Match.WinnerID
}

func (c *Challonge) ApplyUniquePlacings() {
	c.Tournament.UniquePlacings = CalculatePlacings(c.Tournament.NumEntrants)
}
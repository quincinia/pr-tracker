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
	decoder := json.NewDecoder(r)
	err = decoder.Decode(c)
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
func isPresent(playerID int, m Match) bool {
	match := m.Match
	return playerID == match.Player1ID || playerID == match.Player2ID
}

// Returns the index of the first player with the given standing.
func (c *Challonge) FindStanding(standing int) (index int) {
	t := c.Tournament
	for i := range t.Participants {
		if t.Participants[i].Participant.Standing == standing {
			return i
		}
	}
	return -1
}

func (c *Challonge) ApplyResetPoints() {
	
}

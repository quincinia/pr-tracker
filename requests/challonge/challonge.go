// Package challonge contains functions for handling Challonge-based tournaments.
package challonge

// The JSON response puts a wrapper around the tournament object.
type Challonge struct {
	Tourney Tournament `json:"tournament"`
}

// Similar to models.Tournament, however only the fields available in the Challonge response are included.
type Tournament struct {
	Name           string
	URL            string `json:"full_challonge_url"`
	NumEntrants    int    `json:"participants_count"`
	UniquePlacings int    `json:"-"`
	BracketReset   bool   `json:"-"`
	Participants   []Attendee
	Matches        []Match
}

// Similar to models.Attendee, but only containing relevant fields.
type Attendee struct {
	ID       int
	Name     string
	Standing int `json:"final_rank"`
}

type Match struct {
	Player1ID int `json:"player1_id"`
	Player2ID int `json:"player2_id"`
	WinnerID  int `json:"winner_id"`
	LoserID   int `json:"loser_id"`
	Round     int
}

// In JavaScript: Math.trunc(Math.log2(n))
// Do not use with negative numbers.
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

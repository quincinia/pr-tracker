// Package challonge contains functions for handling Challonge-based tournaments.
package challonge

// Similar to models.Tournament, however only the fields available in the Challonge response are included.
type Tournament struct {
	Name           string
	URL            string `json:"full_challonge_url"`
	NumEntrants    int    `json:"participants_count"`
	UniquePlacings int    `json:"-"`
	BracketReset   bool   `json:"-"`
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
}

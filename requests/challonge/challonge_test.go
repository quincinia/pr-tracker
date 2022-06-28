package challonge

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"pr-tracker/requests"
)

func TestDecoder(t *testing.T) {
	var challonge *Challonge

	t.Run("Decode tournament fields", func(t *testing.T) {
		response, err := os.Open("challonge.json")
		if err != nil {
			t.Error(err)
		}
		challonge, err = NewChallonge(response)
		if err != nil {
			t.Error(err)
		}

		tourney := challonge.Tournament
		fmt.Println(tourney.Name, tourney.URL, tourney.NumEntrants)

		if len(tourney.Participants) != 20 {
			t.Error("Incorrect number of participants", len(tourney.Participants))
		}

		if len(tourney.Matches) != 39 {
			t.Error("Incorrect number of matches", len(tourney.Matches))
		}
	})

	t.Run("Decode participants", func(t *testing.T) {
		firstIndex := challonge.FindStanding(1)
		winner := challonge.Tournament.Participants[firstIndex].Participant
		if winner.Name != "Scisto" {
			t.Error("got", winner.Name, "want Scisto")
		}
	})

	t.Run("Decode matches", func(t *testing.T) {
		lastIndex := challonge.FindMatch(39)
		grands := challonge.Tournament.Matches[lastIndex].Match
		if grands.WinnerID != 135182824 {
			t.Error("got", grands.WinnerID, "want", 135182824)
		}
	})
}

func TestTrunc(t *testing.T) {
	var n, exp, value int

	t.Run("Negative n", func(t *testing.T) {
		n = -1
		_, value = requests.TruncLog2(n)
		if value != 0 {
			t.Error("got", value, ", want 0")
		}
	})

	t.Run("Power of 2", func(t *testing.T) {
		n = 16
		exp, value = requests.TruncLog2(n)
		if exp != 4 || value != 16 {
			t.Error("got", exp, ", want 4")
			t.Error("got", value, ", want 16")

		}
	})

	t.Run("Between powers", func(t *testing.T) {
		n = 31
		exp, value = requests.TruncLog2(n)
		if exp != 4 || value != 16 {
			t.Error("got", exp, ", want 4")
			t.Error("got", value, ", want 16")

		}
	})
}

func TestCalcPlacings(t *testing.T) {
	tests := []int{1, 2, 3, 4, 5, 7, 9, 13, 17, 25, 33, 49, 65, 97}
	for index, test := range tests {
		placings := requests.CalculatePlacings(test)
		if placings != index+1 {
			t.Error("got", placings, "want", index+1)
		}
	}
}

func TestResetPoints(t *testing.T) {
	var challonge *Challonge

	t.Run("No reset points given", func(t *testing.T) {
		response, err := os.Open("challonge.json")
		if err != nil {
			t.Error(err)
		}
		challonge, err = NewChallonge(response)
		if err != nil {
			t.Error(err)
		}

		challonge.ApplyResetPoints()

		if challonge.Tournament.BracketReset != false {
			t.Error("got", challonge.Tournament.BracketReset, "want", false)
		}
	})

	t.Run("Give reset points", func(t *testing.T) {
		response, err := os.Open("bracketreset.json")
		if err != nil {
			t.Error(err)
		}
		challonge, err = NewChallonge(response)
		if err != nil {
			t.Error(err)
		}

		challonge.ApplyResetPoints()

		if challonge.Tournament.BracketReset != true {
			t.Error("got", challonge.Tournament.BracketReset, "want", true)
		}
	})
}

func TestConversion(t *testing.T) {
	var challonge *Challonge

	response, err := os.Open("challonge.json")
	if err != nil {
		t.Error(err)
	}
	challonge, err = NewChallonge(response)
	if err != nil {
		t.Error(err)
	}

	tourney, attendees := challonge.ToTournament()

	if tourney.Type.Name != "Challonge" {
		t.Error("got", tourney.Type.Name, "want Challonge")
	}

	if len(attendees) != 20 {
		t.Error("got", len(attendees), "want 20")
	}

	output, _ := json.MarshalIndent(tourney, "", "\t")
	fmt.Println(string(output))

	output, _ = json.MarshalIndent(attendees, "", "\t")
	fmt.Println(string(output))
}

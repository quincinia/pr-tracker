package smashgg

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestDecoder(t *testing.T) {
	var smashgg *Smashgg

	t.Run("Decode tournament fields", func(t *testing.T) {
		response, err := os.Open("smashgg.json")
		if err != nil {
			t.Error(err)
		}

		smashgg, err = NewSmashgg(response)
		if err != nil {
			t.Error(err)
		}

		tourney := smashgg.Data.Tournament
		event := smashgg.Data.Event
		fmt.Println(tourney.Name, event.Name, event.NumEntrants)
		fmt.Println(tourney.URL)

		if len(event.Entrants.Nodes) != 128 {
			t.Error("Incorrect number of participants", len(event.Entrants.Nodes))
		}

		if len(event.Sets.Nodes) != 3 {
			t.Error("Incorrect number of matches", len(event.Sets.Nodes))
		}
	})

	t.Run("Decode participants", func(t *testing.T) {
		firstIndex := smashgg.FindStanding(1)
		winner := smashgg.GetEntrant(firstIndex)
		if winner.Name != "FORT | Cless" {
			t.Error("got", winner.Name, "want FORT | Cless")
		}
	})

	t.Run("Decode matches", func(t *testing.T) {
		if smashgg.ApplyResetPoints() != true {
			t.Error("got", smashgg.ApplyResetPoints(), "want", true)
		}
	})
}

func TestConversion(t *testing.T) {
	var smashgg *Smashgg

	response, err := os.Open("smashgg.json")
	if err != nil {
		t.Error(err)
	}
	smashgg, err = NewSmashgg(response)
	if err != nil {
		t.Error(err)
	}

	tourney, attendees := smashgg.ToTournament()

	if tourney.Type.Name != "start.gg" {
		t.Error("got", tourney.Type.Name, "want start.gg")
	}

	if len(attendees) != 128 {
		t.Error("got", len(attendees), "want 128")
	}

	output, _ := json.MarshalIndent(tourney, "", "\t")
	fmt.Println(string(output))

	output, _ = json.MarshalIndent(attendees, "", "\t")
	fmt.Println(string(output))
}

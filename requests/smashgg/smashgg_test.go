package smashgg

import (
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

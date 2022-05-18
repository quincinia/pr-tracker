package challonge

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestDecoder(t *testing.T) {
	response, err := os.Open("challonge.json")
	if err != nil {
		t.Error(err)
	}
	defer response.Close()

	decoder := json.NewDecoder(response)
	var c Challonge
	err = decoder.Decode(&c)
	if err != nil {
		t.Error(err)
	}

	tourney := c.Tourney
	fmt.Println(tourney.Name, tourney.URL, tourney.NumEntrants)

	if len(tourney.Participants) != 20 {
		t.Error("Incorrect number of participants", len(tourney.Participants))
	}

	if len(tourney.Matches) != 39 {
		t.Error("Incorrect number of matches", len(tourney.Matches))
	}
}

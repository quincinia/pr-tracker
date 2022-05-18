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

func TestTrunc(t *testing.T) {
	var n, exp, value int

	t.Run("Negative n", func(t *testing.T) {
		n = -1;
		_, value = TruncLog2(n)
		if value != 0 {
			t.Error("got", value, ", want 0")
		}
	})

	t.Run("Power of 2", func(t *testing.T) {
		n = 16
		exp, value = TruncLog2(n)
		if exp != 4 || value != 16 {
			t.Error("got", exp, ", want 4")
			t.Error("got", value, ", want 16")

		}
	})

	t.Run("Between powers", func(t *testing.T) {
		n = 31
		exp, value = TruncLog2(n)
		if exp != 4 || value != 16 {
			t.Error("got", exp, ", want 4")
			t.Error("got", value, ", want 16")

		}
	})
}

func TestCalcPlacings(t *testing.T) {
	tests := []int{1, 2, 3, 4, 5, 7, 9, 13, 17, 25, 33, 49, 65, 97}
	for index, test := range tests {
		placings := CalculatePlacings(test)
		if placings != index+1 {
			t.Error("got", placings, "want", index+1)
		}
	}
}
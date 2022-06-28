package models

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setUp() {
	// Database connection is made using init()
	DB.Exec("delete from tournaments")
}

func tearDown() {
	DB.Exec("delete from tournaments")
}

func TestTournament(t *testing.T) {
	var ids []int

	t.Run("Tournament creation", func(t *testing.T) {
		tests := []Tournament{
			{Name: "No Type, No Tier"},
			{Name: "No Type", Tier: C_TIER},
			{Name: "No Tier", Type: CHALLONGE},
		}
		for _, test := range tests {
			err := test.Create()
			if err != nil {
				t.Fatal(err)
			} else {
				fmt.Println("ID:", test.TourneyID)
				ids = append(ids, test.TourneyID)
			}
		}
	})

	t.Run("Tournament retrieval", func(t *testing.T) {
		for _, id := range ids {
			tourney, err := GetTournament(id)
			if err != nil {
				t.Error(err)
			}

			var pass bool
			switch tourney.Name {
			case "No Type, No Tier":
				pass = tourney.Type == nil && tourney.Tier == nil
			case "No Type":
				pass = tourney.Type == nil
			case "No Tier":
				pass = tourney.Tier == nil
			}

			if !pass {
				t.Error("got type:", tourney.Type, "tier:", tourney.Tier, "want", tourney.Name)
			}
		}
	})

	t.Run("Tournament update", func(t *testing.T) {
		tourney, err := GetTournament(ids[0])
		if err != nil {
			t.Error(err)
		}

		tourney.Type = STARTGG
		tourney.Tier = S_TIER

		err = tourney.Update()
		if err != nil {
			t.Error(err)
		}

		tourney, err = GetTournament(ids[0])
		if err != nil {
			t.Error(err)
		}

		if tourney.Type.TypeID != 2 && tourney.Tier.TierID != 4 {
			t.Error("got typeid:", tourney.Type.TypeID, "tierid:", tourney.Tier.TierID, "want 2, 4")
		}
	})
}

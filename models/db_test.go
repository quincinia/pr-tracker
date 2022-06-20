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
}

func tearDown() {
	// DB.Exec("delete from tournaments")
}

func TestTournament(t *testing.T) {
	t.Run("Tournament creation", func(t *testing.T) {
		tests := []Tournament{
			{Name: "No Type, No Tier"},
			{Name: "No Type", Tier: C_TIER},
			{Name: "No Tier", Type: CHALLONGE},
		}
		for _, test := range tests {
			err := test.Create()
			if err != nil {
				t.Error(err)
			} else {
				fmt.Println("ID:", test.TourneyID)
			}
		}
	})
}

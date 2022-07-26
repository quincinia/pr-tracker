// Package smashgg contains structures for handling smash.gg/start.gg tournaments.
// Package assumes tournaments are complete.
package smashgg

import (
	"encoding/json"
	"io"
	"pr-tracker/models"
	"pr-tracker/requests"
)

type SGGQuery struct {
	Query string `json:"query"`
	Variables map[string]string `json:"variables"`
}

type Smashgg struct {
	Data struct {
		Tournament Tournament
		Event      struct {
			Name        string
			NumEntrants int
			Entrants    struct {
				Nodes []Entrant
			}
			Sets struct {
				Nodes []Set
			}
		}
	}
}

type Startgg = Smashgg

type Tournament struct {
	Name string
	URL  string
}

type Entrant struct {
	Name     string
	Standing struct {
		Placement int
	}
}

type Set struct {
	FullRoundText string
	LPlacement    int
}

func NewSmashgg(r io.Reader) (s *Smashgg, err error) {
	s = &Smashgg{}
	decoder := json.NewDecoder(r)
	err = decoder.Decode(s)
	if err != nil {
		return nil, err
	}
	return
}

func (s *Smashgg) FindStanding(standing int) (index int) {
	for i, e := range s.Data.Event.Entrants.Nodes {
		if e.Standing.Placement == standing {
			return i
		}
	}
	return -1
}

func (s *Smashgg) GetEntrant(index int) (entrant *Entrant) {
	if index < 0 || index >= len(s.Data.Event.Entrants.Nodes) {
		return
	}
	return &s.Data.Event.Entrants.Nodes[index]
}

func (s *Smashgg) ApplyResetPoints() bool {
	for _, set := range s.Data.Event.Sets.Nodes {
		if set.FullRoundText == "Grand Final Reset" && set.LPlacement == 2 {
			return true
		}
	}
	return false
}

func (s *Smashgg) ApplyUniquePlacings() int {
	return requests.CalculatePlacings(s.Data.Event.NumEntrants)
}

func (s *Smashgg) ToTournament() (t models.Tournament, as []models.Attendee) {
	t = models.Tournament{
		Type:           models.STARTGG,
		Name:           s.Data.Tournament.Name,
		URL:            s.Data.Tournament.URL,
		NumEntrants:    s.Data.Event.NumEntrants,
		UniquePlacings: s.ApplyUniquePlacings(),
		BracketReset:   s.ApplyResetPoints(),
	}
	for _, e := range s.Data.Event.Entrants.Nodes {
		a := models.Attendee{
			Name: e.Name,
			Standing: e.Standing.Placement,
		}
		as = append(as, a)
	}
	return
}

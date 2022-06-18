// Package smashgg contains structures for handling smash.gg/start.gg tournaments.
// Package assumes tournaments are complete.
package smashgg

import (
	"encoding/json"
	"io"
)

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

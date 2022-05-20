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

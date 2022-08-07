// Package smashgg contains structures for handling smash.gg/start.gg tournaments.
// Package assumes tournaments are complete.
package smashgg

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"pr-tracker/models"
	"pr-tracker/requests"
	"strings"
)

type SGGQuery struct {
	Query     string            `json:"query"`
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
			Name:     e.Name,
			Standing: e.Standing.Placement,
		}
		as = append(as, a)
	}
	return
}

func FromURL(url *url.URL, key string) (t models.FullTournament, err error) {
	query := SGGQuery{
		Query:     "query TournamentEventQuery($tournament: String, $event: String) { tournament(slug: $tournament) { name url(relative: false) } event(slug: $event) { name numEntrants entrants(query: { page: 1, perPage: 500 }) { nodes { id name standing { placement } } } sets(page: 1, perPage: 3, sortType: RECENT) { nodes { fullRoundText lPlacement } } }}",
		Variables: make(map[string]string),
	}
	query.Variables["tournament"] = strings.Split(url.Path, "/")[2]
	query.Variables["event"] = url.Path[1:]

	body, _ := json.Marshal(query)
	req, _ := http.NewRequest("POST", "https://api.smash.gg/gql/alpha", bytes.NewBuffer(body))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", "Bearer "+key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	smashgg, err := NewSmashgg(res.Body)
	if err != nil {
		return
	}

	t.Tournament, t.Attendees = smashgg.ToTournament()

	return
}

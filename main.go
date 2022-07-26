package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"pr-tracker/requests/challonge"
	"pr-tracker/requests/smashgg"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

var (
	SMASHGG_KEY   = os.Getenv("SMASHGG_KEY")
	CHALLONGE_KEY = os.Getenv("CHALLONGE_KEY")
	DB_USERNAME   = os.Getenv("DB_USERNAME")
	DB_PASSWORD   = os.Getenv("DB_PASSWORD")
)

func main() {
	// fmt.Printf("%v %T\n", SMASHGG_KEY, SMASHGG_KEY)
	// fmt.Printf("%v %T\n", CHALLONGE_KEY, CHALLONGE_KEY)
	// fmt.Printf("%v %T\n", DB_USERNAME, DB_USERNAME)
	// fmt.Printf("%v %T\n", DB_PASSWORD, DB_PASSWORD)

	flag.Parse()
	args := flag.Args()
	input, err := url.Parse(args[0])
	if err != nil {
		log.Fatalln(err)
	}

	switch input.Host {
	case "challonge.com":
		fmt.Println("challonge.com")
		reqURL := url.URL{
			Scheme: "https",
			Host:   "api.challonge.com",
			Path:   "/v1/tournaments" + input.Path + ".json",
		}
		query := url.Values{}
		query.Add("api_key", CHALLONGE_KEY)
		query.Add("include_participants", "1")
		query.Add("include_matches", "1")
		reqURL.RawQuery = query.Encode()
		fmt.Println(reqURL.String())

		res, err := http.Get(reqURL.String())
		if err != nil {
			fmt.Print("network error: ")
			log.Fatalln(err)
		}
		defer res.Body.Close()

		challonge, err := challonge.NewChallonge(res.Body)
		if err != nil {
			fmt.Print("parse error: ")
			log.Fatalln(err)
		}

		tournament, _ := challonge.ToTournament()
		fmt.Println(tournament)

	case "smash.gg", "start.gg":
		fmt.Println("smash.gg/start.gg")
		query := smashgg.SGGQuery{
			Query:     "query TournamentEventQuery($tournament: String, $event: String) { tournament(slug: $tournament) { name url(relative: false) } event(slug: $event) { name numEntrants entrants(query: { page: 1, perPage: 500 }) { nodes { id name standing { placement } } } sets(page: 1, perPage: 3, sortType: RECENT) { nodes { fullRoundText lPlacement } } }}",
			Variables: make(map[string]string),
		}
		query.Variables["tournament"] = strings.Split(input.Path, "/")[2]
		query.Variables["event"] = input.Path[1:]

		body, _ := json.Marshal(query)
		req, _ := http.NewRequest("POST", "https://api.smash.gg/gql/alpha", bytes.NewBuffer(body))
		req.Header.Add("content-type", "application/json")
		req.Header.Add("authorization", "Bearer " + SMASHGG_KEY)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Print("network error: ")
			log.Fatalln(err)
		}
		defer res.Body.Close()

		smashgg, err := smashgg.NewSmashgg(res.Body)
		if err != nil {
			fmt.Print("parse error: ")
			log.Fatalln(err)
		}

		tournament, _ := smashgg.ToTournament()
		fmt.Println(tournament)
	}
}

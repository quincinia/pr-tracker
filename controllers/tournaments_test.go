package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	. "pr-tracker/models"
	"strconv"
	"testing"
)

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	// tearDown()
	os.Exit(code)
}

func setUp() {
	// Database connection is made using init()
	DB.Exec("delete from tournaments")
}

func tearDown() {
	DB.Exec("delete from tournaments")
}
func TestGetTournament(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/tournament/", TournamentsRouter)
	writer := httptest.NewRecorder()

	tourney := FullTournament{
		Tournament: Tournament{Name: "GET Tournament"},
		Attendees:  []Attendee{{Name: "GET Attendee", Standing: 1}},
	}
	err := tourney.Create()
	if err != nil {
		t.Fatal(err)
	}

	request, _ := http.NewRequest("GET", "/tournament/"+strconv.Itoa(tourney.TourneyID), nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}

	var body FullTournament
	json.Unmarshal(writer.Body.Bytes(), &body)
	if body.Name != "GET Tournament" {
		t.Error("got", body.Name, "want 'GET Tournament'")
	}
	if len(body.Attendees) != 1 {
		t.Error("got", len(body.Attendees), "want 1")
	}
}

func TestPostTournament(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/tournament/", TournamentsRouter)
	writer := httptest.NewRecorder()

	tourney := FullTournament{
		Tournament: Tournament{Name: "POST Tournament"},
		Attendees:  []Attendee{{Name: "POST Attendee", Standing: 1}},
	}

	output, _ := json.Marshal(tourney)
	body := bytes.NewReader(output)

	request, _ := http.NewRequest("POST", "/tournament/", body)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Error("got", writer.Code, "want 200")
	}
}

func TestDeleteTournament(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/tournament/", TournamentsRouter)
	writer := httptest.NewRecorder()

	tourney := FullTournament{
		Tournament: Tournament{Name: "DELETE Tournament"},
		Attendees:  []Attendee{{Name: "DELETE Attendee", Standing: 1}},
	}
	err := tourney.Create()
	if err != nil {
		t.Fatal(err)
	}

	request, _ := http.NewRequest("DELETE", "/tournament/"+strconv.Itoa(tourney.TourneyID), nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Error("got", writer.Code, "want 200")
	}
}

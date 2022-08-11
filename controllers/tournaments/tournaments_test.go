package tournaments

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	. "pr-tracker/models"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setUp() {
	godotenv.Load("../../.env")
	Connect(os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"))
	DB.Exec("delete from tournaments where url = ''")
}

func tearDown() {
	DB.Exec("delete from tournaments where url = ''")
}

func TestGetTournament(t *testing.T) {
	mux := httprouter.New()
	AddRoutes("/tournaments", mux)
	writer := httptest.NewRecorder()

	tourney := FullTournament{
		Tournament: Tournament{Name: "GET Tournament"},
		Attendees:  []Attendee{{Name: "GET Attendee", Standing: 1}},
	}
	err := tourney.Create()
	if err != nil {
		t.Fatal(err)
	}

	request, _ := http.NewRequest("GET", "/tournaments/"+strconv.Itoa(tourney.TourneyID), nil)
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

func TestGetTournaments(t *testing.T) {
	mux := httprouter.New()
	AddRoutes("/tournaments", mux)
	writer := httptest.NewRecorder()

	tests := []FullTournament{
		{
			Tournament: Tournament{Name: "GET Tournaments 1"},
			Attendees:  []Attendee{{Name: "GET Attendee 1", Standing: 1}},
		},
		{
			Tournament: Tournament{Name: "GET Tournaments 2"},
			Attendees:  []Attendee{{Name: "GET Attendee 2", Standing: 1}},
		},
	}

	for i := range tests {
		err := tests[i].Create()
		if err != nil {
			t.Fatal(err)
		}
	}

	request, _ := http.NewRequest("GET", "/tournaments/", nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}

	var body []Tournament
	json.Unmarshal(writer.Body.Bytes(), &body)
	if len(body) < 2 {
		t.Error("got", len(body), "want >=2")
	}
}

func TestPostTournament(t *testing.T) {
	mux := httprouter.New()
	AddRoutes("/tournaments", mux)
	writer := httptest.NewRecorder()

	tourney := FullTournament{
		Tournament: Tournament{Name: "POST Tournament"},
		Attendees:  []Attendee{{Name: "POST Attendee", Standing: 1}},
	}

	output, _ := json.Marshal(tourney)
	body := bytes.NewReader(output)

	request, _ := http.NewRequest("POST", "/tournaments/", body)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Error("got", writer.Code, "want 200")
	}
}

func TestDeleteTournament(t *testing.T) {
	mux := httprouter.New()
	AddRoutes("/tournaments", mux)
	writer := httptest.NewRecorder()

	tourney := FullTournament{
		Tournament: Tournament{Name: "DELETE Tournament"},
		Attendees:  []Attendee{{Name: "DELETE Attendee", Standing: 1}},
	}
	err := tourney.Create()
	if err != nil {
		t.Fatal(err)
	}

	request, _ := http.NewRequest("DELETE", "/tournaments/"+strconv.Itoa(tourney.TourneyID), nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Error("got", writer.Code, "want 200")
	}
}

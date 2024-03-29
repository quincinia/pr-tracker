package attendees

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
	DB.Exec("delete from players where playerid > 6")
	DB.Exec("delete from tournaments where url = ''")
}

func tearDown() {
	DB.Exec("delete from players where playerid > 6")
	DB.Exec("delete from tournaments where url = ''")
}

func TestGetAttendee(t *testing.T) {
	mux := httprouter.New()
	AddRoutes("/attendees", mux)
	writer := httptest.NewRecorder()

	tourney := FullTournament{
		Tournament: Tournament{Name: "GET Attendee Tournament"},
		Attendees:  []Attendee{{Name: "GET Attendee Attendee", Standing: 1}},
	}
	err := tourney.Create()
	if err != nil {
		t.Fatal(err)
	}

	request, _ := http.NewRequest("GET", "/attendees/"+strconv.Itoa(tourney.Attendees[0].AttendeeID), nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}

	var body Attendee
	json.Unmarshal(writer.Body.Bytes(), &body)
	if body.Name != "GET Attendee Attendee" {
		t.Error("got", body.Name, "want 'GET Attendee Attendee'")
	}
}

func TestPutAttendee(t *testing.T) {
	mux := httprouter.New()
	AddRoutes("/attendees", mux)
	writer := httptest.NewRecorder()

	player1 := Player{Name: "PUT Attendee player 1"}
	player1.Create()
	player2 := Player{Name: "PUT Attendee player 2"}
	player2.Create()
	tourney := FullTournament{
		Tournament: Tournament{Name: "GET Attendee Tournament"},
		Attendees:  []Attendee{{Player: &player1, Name: "GET Attendee Attendee 1", Standing: 1}},
	}
	err := tourney.Create()
	if err != nil {
		t.Fatal(err)
	}

	output, _ := json.Marshal(Attendee{Player: &player2})
	body := bytes.NewReader(output)

	request, _ := http.NewRequest("PUT", "/attendees/"+strconv.Itoa(tourney.Attendees[0].AttendeeID), body)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Error("got", writer.Code, "want 200")
	}
}

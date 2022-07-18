package players

import (
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
	DB.Exec("delete from players")
	DB.Exec("delete from tournaments")
}

func tearDown() {
	// DB.Exec("delete from players")
}

func TestGetPlayer(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/players/", PlayersRouter)
	writer := httptest.NewRecorder()

	player := Player{Name: "GET Player test"}
	err := player.Create()
	if err != nil {
		t.Fatal(err)
	}

	tourney := FullTournament{
		Tournament: Tournament{Name: "GET Player Tournament"},
		Attendees:  []Attendee{{Player: &player, Name: "GET Player Attendee", Standing: 1}},
	}
	err = tourney.Create()
	if err != nil {
		t.Fatal(err)
	}

	request, _ := http.NewRequest("GET", "/players/"+strconv.Itoa(player.PlayerID), nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}

	var body PlayerAttendance
	json.Unmarshal(writer.Body.Bytes(), &body)
	if body.Name != player.Name {
		t.Error("got", body.Name, "want", player.Name)
	}
	if len(body.Attendance) != 1 && body.Attendance[0].Name != "GET Player Attendee" {
		t.Error("got", body.Attendance[0].Name, "want 'GET Player Attendee'")
	}
}

func TestGetPlayers(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/players/", PlayersRouter)
	writer := httptest.NewRecorder()

	player1 := Player{Name: "X.1 Player"}
	player1.Create()
	player2 := Player{Name: "X.2 Player"}
	player2.Create()
	tournaments := []FullTournament{
		{
			Tournament: Tournament{Name: "GET Players Tournament 1"},
			Attendees:  []Attendee{
				{Player: &player1, Name: "GET Players Attendee 1.1", Standing: 1}, 
				{Player: &player2, Name: "GET Players Attendee 1.2", Standing: 2},
			},
		},
		{
			Tournament: Tournament{Name: "GET Players Tournament 2"},
			Attendees:  []Attendee{
				{Player: &player1, Name: "GET Players Attendee 2.1", Standing: 1},
			},
		},
	}

	for i := range tournaments {
		err := tournaments[i].Create()
		if err != nil {
			t.Fatal(err)
		}
	}

	request, _ := http.NewRequest("GET", "/players/", nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}

	var body []PlayerAttendance
	json.Unmarshal(writer.Body.Bytes(), &body)
	if len(body) != 2 {
		t.Error("got", len(body), "want 2")
	}
	// Assuming the players are returned in the order defined above
	if body[0].Player != player1 {
		t.Error("got", body[0].Player, "want", player1)
	}
	if body[1].Player != player2 {
		t.Error("got", body[1].Player, "want", player2)
	}

	if len(body[0].Attendance) != 2 {
		t.Error("got", len(body[0].Attendance), "want 2")
	}
	if len(body[1].Attendance) != 1 {
		t.Error("got", len(body[1].Attendance), "want 1")
	}
}

func TestPostPlayer(t *testing.T) {
	// mux := http.NewServeMux()
	// mux.HandleFunc("/players/", PlayersRouter)
	// writer := httptest.NewRecorder()
}

func TestDeletePlayer(t *testing.T) {
	// mux := http.NewServeMux()
	// mux.HandleFunc("/players/", PlayersRouter)
	// writer := httptest.NewRecorder()
}

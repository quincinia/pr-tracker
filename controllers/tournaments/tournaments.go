// Package tournaments contains REST handlers to work with handling tournaments
// tournaments.go contains handlers for managing tournament data
package tournaments

import (
	"encoding/json"
	"net/http"
	"pr-tracker/models"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Prefix paths should not end with "/"
func AddRoutes(prefix string, router *httprouter.Router) {
	router.HandlerFunc("GET", prefix+"/", getTournaments)
	router.Handle("GET", prefix+"/:id", getTournament)
	router.HandlerFunc("POST", prefix+"/", postTournament)
	router.Handle("DELETE", prefix+"/:id", deleteTournament)
}

// Deprecated, using httprouter
// func TournamentsRouter(w http.ResponseWriter, r *http.Request) {
// 	var err error
// 	// fmt.Println(r.URL.Path)
// 	switch r.Method {
// 	case "GET":
// 		err = getTournament(w, r)
// 	case "POST":
// 		err = postTournament(w, r)
// 	// case "PUT":
// 	// 	err = handlePut(w, r)
// 	case "DELETE":
// 		err = deleteTournament(w, r)
// 	}
// 	if err != nil {
// 		fmt.Println(err)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }

// Returns full tournament data (including attendees)
func getTournament(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tournament, err := models.GetTournament(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	attendees, err := models.GetAttendees(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fulltournament := models.FullTournament{Tournament: tournament, Attendees: attendees}

	output, err := json.MarshalIndent(&fulltournament, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

// Returns only the tournament data (no attendees)
func getTournaments(w http.ResponseWriter, r *http.Request) {
	tournaments, err := models.GetTournaments()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output, err := json.MarshalIndent(&tournaments, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

func postTournament(w http.ResponseWriter, r *http.Request) {
	var err error
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)

	var fulltournament models.FullTournament
	err = json.Unmarshal(body, &fulltournament)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = fulltournament.Create()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

func deleteTournament(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tournament, err := models.GetTournament(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	attendees, err := models.GetAttendees(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fulltournament := models.FullTournament{Tournament: tournament, Attendees: attendees}

	err = fulltournament.Tournament.Delete()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// No longer needed since attendees are deleted with the tournament
	// err = models.DeleteAttendees(fulltournament.TourneyID)
	// if err != nil {
	// 	return
	// }

	w.WriteHeader(200)
}

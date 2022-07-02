// Package controllers contains REST handlers to work with various routes of the application
// tournaments.go contains handlers for managing tournament data
package controllers

import (
	"encoding/json"
	"net/http"
	"path"
	"pr-tracker/models"
	"strconv"
)

func TournamentsRouter(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case "GET":
		err = getTournament(w, r)
	case "POST":
		err = postTournament(w, r)
	// case "PUT":
	// 	err = handlePut(w, r)
	case "DELETE":
		err = deleteTournament(w, r)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Returns full tournament data (including attendees)
func getTournament(w http.ResponseWriter, r *http.Request) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return
	}

	tournament, err := models.GetTournament(id)
	if err != nil {
		return
	}

	attendees, err := models.GetAttendees(id)
	if err != nil {
		return
	}

	fulltournament := models.FullTournament{Tournament: tournament, Attendees: attendees}

	output, err := json.MarshalIndent(&fulltournament, "", "\t")
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return
}

func postTournament(w http.ResponseWriter, r *http.Request) (err error) {
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	var fulltournament models.FullTournament
	json.Unmarshal(body, &fulltournament)
	err = fulltournament.Create()
	if err != nil {
		return
	}
	w.WriteHeader(200)
	return
}

func deleteTournament(w http.ResponseWriter, r *http.Request) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return
	}

	tournament, err := models.GetTournament(id)
	if err != nil {
		return
	}

	attendees, err := models.GetAttendees(id)
	if err != nil {
		return
	}

	fulltournament := models.FullTournament{Tournament: tournament, Attendees: attendees}
	if err != nil {
		return
	}

	err = fulltournament.Tournament.Delete()
	if err != nil {
		return
	}

	err = models.DeleteAttendees(fulltournament.TourneyID)
	if err != nil {
		return
	}

	w.WriteHeader(200)
	return
}

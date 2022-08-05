// This file contains functions for parsing and processing the forms from templates/tournament_view.html
package tournaments

import (
	"net/http"
	"pr-tracker/models"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func ProcessTourneyAttendees(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	attendees, err := models.GetAttendees(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// No guarantees are made about the order of the returned attendees
	// If multiple attendees are given the same player, only one update will be performed
	// Please fill out the form correctly when submitting
	for i := range attendees {
		attendeeid := strconv.Itoa(attendees[i].AttendeeID)
		playerid := r.FormValue(attendeeid)
		if playerid == "" {
			attendees[i].Player = nil
			continue
		}
		insert, err := strconv.Atoi(playerid)
		if err != nil {
			continue
		}
		attendees[i].Player = &models.Player{PlayerID: insert}
	}

	for i := range attendees {
		attendees[i].Update()
	}

	// Expects urls of the form: /tournaments/:tourneyid/edit
	// Not performing validation right now
	http.Redirect(w, r, strings.TrimSuffix(r.URL.String(), "/edit"), http.StatusFound)
}

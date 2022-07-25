package attendees

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"pr-tracker/models"
	"strconv"
)

func AttendeesRouter(w http.ResponseWriter, r *http.Request) {
	var err error
	// fmt.Println(r.URL.Path)
	switch r.Method {
	case "GET":
		err = getAttendee(w, r)
	// case "POST":
	// 	err = postAttendee(w, r)
	case "PUT":
		err = putAttendee(w, r)
		// case "DELETE":
		// 	err = deleteAttendee(w, r)
	}
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getAttendee(w http.ResponseWriter, r *http.Request) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return
	}

	attendee, err := models.GetAttendee(id)
	if err != nil {
		return
	}

	output, err := json.Marshal(attendee)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return
}

func putAttendee(w http.ResponseWriter, r *http.Request) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return
	}

	attendee, err := models.GetAttendee(id)
	if err != nil {
		return
	}

	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)

	var a models.Attendee
	json.Unmarshal(body, &a)
	
	// Only the player can be changed
	attendee.Player = a.Player

	err = attendee.Update()
	if err != nil {
		return
	}

	w.WriteHeader(200)
	return
}

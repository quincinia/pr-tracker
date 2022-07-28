package attendees

import (
	"encoding/json"
	"net/http"
	"pr-tracker/models"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func AddRoutes(prefix string, router *httprouter.Router) {
	router.GET(prefix+"/:id", getAttendee)
	router.PUT(prefix+"/:id", putAttendee)
}

// Deprecated
// func AttendeesRouter(w http.ResponseWriter, r *http.Request) {
// 	var err error
// 	// fmt.Println(r.URL.Path)
// 	switch r.Method {
// 	case "GET":
// 		err = getAttendee(w, r)
// 	// case "POST":
// 	// 	err = postAttendee(w, r)
// 	case "PUT":
// 		err = putAttendee(w, r)
// 		// case "DELETE":
// 		// 	err = deleteAttendee(w, r)
// 	}
// 	if err != nil {
// 		fmt.Println(err)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }

func getAttendee(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	attendee, err := models.GetAttendee(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output, err := json.Marshal(attendee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

func putAttendee(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	attendee, err := models.GetAttendee(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)

	var a models.Attendee
	err = json.Unmarshal(body, &a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Only the player can be changed
	attendee.Player = a.Player

	err = attendee.Update()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

package players

import (
	"encoding/json"
	"net/http"
	"pr-tracker/models"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func AddRoutes(prefix string, router *httprouter.Router) {
	router.HandlerFunc("GET", prefix+"/", getPlayers)
	router.Handle("GET", prefix+"/:id", getPlayer)
	router.HandlerFunc("POST", prefix+"/", postPlayer)
	router.Handle("DELETE", prefix+"/:id", deletePlayer)
}

// Deprecated
// func PlayersRouter(w http.ResponseWriter, r *http.Request) {
// 	var err error
// 	// fmt.Println(r.URL.Path)
// 	switch r.Method {
// 	case "GET":
// 		err = getPlayer(w, r)
// 	case "POST":
// 		err = postPlayer(w, r)
// 	// case "PUT":
// 	// 	err = handlePut(w, r)
// 	case "DELETE":
// 		err = deletePlayer(w, r)
// 	}
// 	if err != nil {
// 		fmt.Println(err)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }

func getPlayer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	player, err := models.GetPlayer(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output, err := json.Marshal(player)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

func getPlayers(w http.ResponseWriter, r *http.Request) {
	players, err := models.GetPlayers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output, err := json.Marshal(players)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

func postPlayer(w http.ResponseWriter, r *http.Request) {
	var err error
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)

	var p models.Player
	err = json.Unmarshal(body, &p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = p.Create()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

func deletePlayer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = models.DeletePlayer(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

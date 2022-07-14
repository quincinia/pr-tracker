package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"pr-tracker/models"
	"strconv"
)

func PlayersRouter(w http.ResponseWriter, r *http.Request) {
	var err error
	// fmt.Println(r.URL.Path)
	switch r.Method {
	case "GET":
		err = getPlayer(w, r)
	case "POST":
		err = postPlayer(w, r)
	// case "PUT":
	// 	err = handlePut(w, r)
	case "DELETE":
		err = deletePlayer(w, r)
	}
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getPlayer(w http.ResponseWriter, r *http.Request) (err error) {
	if r.URL.Path == "/players/" {
		return getPlayers(w, r)
	}

	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return
	}

	player, err := models.GetPlayer(id)
	if err != nil {
		return
	}

	output, err := json.Marshal(player)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return
}

func getPlayers(w http.ResponseWriter, r *http.Request) (err error) {
	players, err := models.GetPlayers()
	if err != nil {
		return
	}

	output, err := json.Marshal(players)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return
}

func postPlayer(w http.ResponseWriter, r *http.Request) (err error) {
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)

	var p models.Player
	json.Unmarshal(body, &p)

	err = p.Create()
	if err != nil {
		return
	}

	w.WriteHeader(200)
	return
}

func deletePlayer(w http.ResponseWriter, r *http.Request) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return
	}

	err = models.DeletePlayer(id)
	if err != nil {
		return
	}

	w.WriteHeader(200)
	return
}

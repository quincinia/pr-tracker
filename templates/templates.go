// Package templates contains code needed to render the site's templates
package templates

import (
	"html/template"
	"net/http"
	"pr-tracker/models"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func getPlayers() (ps []models.Player, err error) {
	rows, err := models.DB.Query("select playerid, name from players")
	if err != nil {
		return
	}

	for rows.Next() {
		var p models.Player
		err = rows.Scan(&p.PlayerID, &p.Name)
		if err != nil {
			return
		}

		ps = append(ps, p)
	}
	rows.Close()
	return
}

func getAttendance() (as []models.Attendee, err error) {
	// Technically we don't need all these fields, just including them for completeness
	// That and I don't want to make custom structs just to do this
	rows, err := models.DB.Query("select attendeeid, tourney, player, name, standing from attendees where player is not null")
	if err != nil {
		return
	}

	for rows.Next() {
		// Only need the playerid, not the name
		a := models.Attendee{Player: &models.Player{}}
		err = rows.Scan(&a.AttendeeID, &a.Tourney, &a.Player.PlayerID, &a.Name, &a.Standing)
		if err != nil {
			return
		}

		as = append(as, a)
	}
	rows.Close()
	return
}

func RenderTable(w http.ResponseWriter, r *http.Request) {
	var (
		err         error
		tmpl        *template.Template
		tournaments []models.Tournament
		tmap        map[int]int // maps tourneyids to their place in the array
		players     []models.Player
		pmap        map[int]int // maps playerids to their place in the array
		attendance  []models.Attendee

		// The order of a player's attendance record matches the order of the tournaments array
		// eg. the attendance record for tournaments[i] for a specific player can be found in rows[player][i]
		// The attendance record is nil if the player did not attend that tournament
		rows map[models.Player][]*models.Attendee
	)
	tmpl = template.Must(template.ParseFiles("./templates/layout.html", "./templates/table.html"))
	tmap = make(map[int]int)
	pmap = make(map[int]int)
	rows = make(map[models.Player][]*models.Attendee)

	tournaments, err = models.GetTournaments()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	players, err = getPlayers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	attendance, err = getAttendance()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i := range tournaments {
		tmap[tournaments[i].TourneyID] = i
	}

	for i := range players {
		pmap[players[i].PlayerID] = i
		rows[players[i]] = make([]*models.Attendee, len(tournaments))
	}

	for i := range attendance {
		// Find the player the attendance record refers to
		pindex := pmap[attendance[i].Player.PlayerID]

		// Find the tournament the attendance record refers to
		tindex := tmap[attendance[i].Tourney]

		// Add the tournament to that player's attendance record
		rows[players[pindex]][tindex] = &attendance[i]
	}

	table := struct {
		Tournaments []models.Tournament
		Rows        map[models.Player][]*models.Attendee
		Tmap        map[int]int
	}{
		Tournaments: tournaments,
		Rows:        rows,
		Tmap:        tmap,
	}
	err = tmpl.ExecuteTemplate(w, "layout", table)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RenderTourneySelect(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/layout.html", "./templates/tournament_select.html"))
	tournaments, err := models.GetTournaments()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", tournaments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RenderTourneyView(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tmpl := template.Must(template.ParseFiles("./templates/layout.html", "./templates/tournament_view.html"))
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

	players, err := getPlayers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	attendees, err := models.GetAttendees(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tiers := []string{"", "C", "B", "A", "S"}

	data := struct {
		Players []models.Player
		models.FullTournament
		Tiers []string
	}{
		players,
		models.FullTournament{Tournament: tournament, Attendees: attendees},
		tiers,
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RenderPlayerSelect(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/layout.html", "./templates/player_select.html"))
	players, err := getPlayers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", players)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RenderPlayerView(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tmpl := template.Must(template.ParseFiles("./templates/layout.html", "./templates/player_view.html"))
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tournaments, err := models.GetTournaments()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	player, err := models.GetPlayer(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmap := make(map[int]models.Tournament)
	for _, t := range tournaments {
		tmap[t.TourneyID] = t
	}

	data := struct {
		Tmap map[int]models.Tournament
		models.PlayerAttendance
	}{
		tmap,
		player,
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

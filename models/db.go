// db defines functions for working with our database.
package models

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Method 1 from https://www.alexedwards.net/blog/organising-database-access
var DB *sql.DB

// Do this in main instead
func Connect(username, password string) {
	var err error
	connString := fmt.Sprintf("user=%s password=%s dbname=pr_tracker sslmode=disable", username, password)
	DB, err = sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}
	fmt.Println("connected to db")
}

func (t *Tournament) Create() (err error) {
	query := "insert into tournaments (type, name, url, numentrants, uniqueplacings, bracketreset, tier) values ($1, $2, $3, $4, $5, $6, $7) returning tourneyid"

	stmt, err := DB.Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()

	var typeid, tierid *int
	if t.Type != nil {
		typeid = &t.Type.TypeID
	}
	if t.Tier != nil {
		tierid = &t.Tier.TierID
	}

	err = stmt.QueryRow(typeid, t.Name, t.URL, t.NumEntrants, t.UniquePlacings, t.BracketReset, tierid).Scan(&t.TourneyID)
	return
}

func GetTournament(id int) (t Tournament, err error) {
	var (
		ntt NullTourneyType
		nt  NullTier
	)
	query := `
		select tourneyid, type, tourneytypes.name, tournaments.name, url, numentrants, uniqueplacings, bracketreset, tier, tiers.name, tiers.multiplier
		from tournaments
		left outer join tourneytypes on type = typeid
		left outer join tiers on tier = tierid
		where tourneyid = $1;
	`
	err = DB.QueryRow(query, id).Scan(&t.TourneyID, &ntt.TypeID, &ntt.Name, &t.Name, &t.URL, &t.NumEntrants, &t.UniquePlacings, &t.BracketReset, &nt.TierID, &nt.Name, &nt.Multiplier)
	t.Type = ntt.ToTourneyType()
	t.Tier = nt.ToTier()
	return
}

func GetTournaments() (ts []Tournament, err error) {
	var (
		ntt NullTourneyType
		nt  NullTier
	)

	rows, err := DB.Query(`
		select tourneyid, typeid, tourneytypes.name, tournaments.name, url, numentrants, uniqueplacings, bracketreset, tierid, tiers.name, multiplier
		from tournaments
		left outer join tourneytypes on type = typeid
		left outer join tiers on tier = tierid
	`)
	if err != nil {
		return
	}

	for rows.Next() {
		t := Tournament{}
		err = rows.Scan(&t.TourneyID, &ntt.TypeID, &ntt.Name, &t.Name, &t.URL, &t.NumEntrants, &t.UniquePlacings, &t.BracketReset, &nt.TierID, &nt.Name, &nt.Multiplier)
		if err != nil {
			return
		}
		t.Type = ntt.ToTourneyType()
		t.Tier = nt.ToTier()
		ts = append(ts, t)
	}
	rows.Close()
	return
}

func GetAttendees(tourneyid int) (as []Attendee, err error) {
	var np NullPlayer

	rows, err := DB.Query(`
		select attendeeid, tourney, player, players.name, attendees.name, standing
		from attendees
		left outer join players on player = playerid
		where tourney = $1
	`, tourneyid)
	if err != nil {
		return
	}

	for rows.Next() {
		a := Attendee{}
		err = rows.Scan(&a.AttendeeID, &a.Tourney, &np.PlayerID, &np.Name, &a.Name, &a.Standing)
		if err != nil {
			return
		}
		a.Player = np.ToPlayer()
		as = append(as, a)
	}
	rows.Close()
	return
}

func (t *Tournament) Update() (err error) {
	var typeid, tierid *int
	query := "update tournaments set type = $2, name = $3, url = $4, numentrants = $5, uniqueplacings = $6, bracketreset = $7, tier = $8 where tourneyid = $1"
	if t.Type != nil {
		typeid = &t.Type.TypeID
	}
	if t.Tier != nil {
		tierid = &t.Tier.TierID
	}
	_, err = DB.Exec(query, t.TourneyID, typeid, t.Name, t.URL, t.NumEntrants, t.UniquePlacings, t.BracketReset, tierid)
	return
}

// Deleted tournaments will delete their attendees
func (t *Tournament) Delete() (err error) {
	_, err = DB.Exec("delete from tournaments where tourneyid = $1", t.TourneyID)
	return
}

func (a *Attendee) Create() (err error) {
	// The only time Create() should be called is when a tournament is first created
	// Prioritizing getting all the attendees in the db, even if that means removing any set players
	// A player should not be assigned to more than 1 attendee
	query := `
		insert into attendees (tourney, player, name, standing) values ($1, $2, $3, $4)
		on conflict on constraint attendees_tourney_player_key do update set player = null
		returning attendeeid
	`

	stmt, err := DB.Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()

	var playerid *int
	if a.Player != nil {
		playerid = &a.Player.PlayerID
	}

	err = stmt.QueryRow(a.Tourney, playerid, a.Name, a.Standing).Scan(&a.AttendeeID)
	return
}

func GetAttendee(attendeeid int) (a Attendee, err error) {
	var np NullPlayer
	query := `
		select attendeeid, tourney, player, players.name, attendees.name, standing
		from attendees
		left outer join players on player = playerid
		where attendeeid = $1
	`
	err = DB.QueryRow(query, attendeeid).Scan(&a.AttendeeID, &a.Tourney, &np.PlayerID, &np.Name, &a.Name, &a.Standing)
	a.Player = np.ToPlayer()
	return
}

func (a *Attendee) Update() (err error) {
	var playerid *int
	if a.Player != nil {
		playerid = &a.Player.PlayerID

		// If a player has been assigned to an existing attendee, they will be moved to the currently updated attendee
		// A combination of (tourney, player) has to be unique (see setup.sql)
		DB.Exec("delete from attendees where tourney = $1 and playerid = $2", a.Tourney, playerid)
	}

	query := "update attendees set tourney = $2, player = $3, name = $4, standing = $5 where attendeeid = $1"
	_, err = DB.Exec(query, a.AttendeeID, a.Tourney, playerid, a.Name, a.Standing)
	return
}

func (a *Attendee) Delete() (err error) {
	_, err = DB.Exec("delete from attendees where attendeeid = $1", a.AttendeeID)
	return
}

// Deprecated
func (as *Attendees) Save(tourneyid int) (failindex int, err error) {
	failindex = -1
	if *as == nil {
		return
	}
	for i := range *as {
		(*as)[i].Tourney = tourneyid
		err = (*as)[i].Create()
		if err != nil {
			return i, err
		}
	}
	return
}

// Omitting the failindex because I'm not expecting to use it
func CreateAttendees(attendees []Attendee, tourneyid int) (err error) {
	for i := range attendees {
		attendees[i].Tourney = tourneyid
		err = attendees[i].Create()
		if err != nil {
			return
		}
	}
	return
}

func DeleteAttendees(tourneyid int) (err error) {
	_, err = DB.Exec("delete from attendees where tourney = $1", tourneyid)
	return
}

func (ft *FullTournament) Create() (err error) {
	err = ft.Tournament.Create()
	if err != nil {
		return
	}
	for i := range ft.Attendees {
		ft.Attendees[i].Tourney = ft.TourneyID
		err = ft.Attendees[i].Create()
		if err != nil {
			return
		}
	}
	return
}

func (p *Player) Create() (err error) {
	err = DB.QueryRow("insert into players (name) values ($1) returning playerid", p.Name).Scan(&p.PlayerID)
	return
}

// Player data will also return attendance
func GetPlayer(id int) (pa PlayerAttendance, err error) {
	err = DB.QueryRow("select playerid, name from players where playerid = $1", id).Scan(&pa.PlayerID, &pa.Name)
	if err != nil {
		return
	}

	rows, err := DB.Query("select attendeeid, tourney, name, standing from attendees where player = $1", id)
	if err != nil {
		return
	}

	for rows.Next() {
		// See PlayerAttendance definition
		var a struct {
			AttendeeID int    `json:"attendeeID"`
			Tourney    int    `json:"tourney"`
			Name       string `json:"name"`
			Standing   int    `json:"standing"`
		}
		err = rows.Scan(&a.AttendeeID, &a.Tourney, &a.Name, &a.Standing)
		if err != nil {
			return
		}
		pa.Attendance = append(pa.Attendance, a)
	}
	rows.Close()
	return
}

func GetPlayers() (pas []PlayerAttendance, err error) {
	rows, err := DB.Query("select playerid from players")
	if err != nil {
		return
	}

	for rows.Next() {
		var (
			id int
			pa PlayerAttendance
		)
		err = rows.Scan(&id)
		if err != nil {
			return
		}

		pa, err = GetPlayer(id)
		if err != nil {
			return
		}

		pas = append(pas, pa)
	}
	rows.Close()
	return
}

// Deleted players will remove their attendance
func DeletePlayer(id int) (err error) {
	_, err = DB.Exec("delete from players where playerid = $1", id)
	return
}

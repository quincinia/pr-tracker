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
func init() {
	var err error
	DB, err = sql.Open("postgres", "user=jacob dbname=pr_tracker password=test sslmode=disable")
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

func RetrieveTournament(id int) (t Tournament, err error) {
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

func (t *Tournament) Delete() (err error) {
	_, err = DB.Exec("delete from tournaments where tourneyid = $1", t.TourneyID)
	return
}

func (a *Attendee) Create() (err error) {
	query := "insert into attendees (tourney, player, name, standing) values ($1, $2, $3, $4) returning attendeeid"
	
	stmt, err := DB.Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()

	var playerid *int
	if a.Player != nil {
		playerid = &a.Player.PlayerID
	}

	err = stmt.QueryRow(query, a.Tourney, playerid, a.Name, a.Standing).Scan(&a.AttendeeID)
	return
}

func RetrieveAttendee(id int) (a Attendee, err error) {
	var np NullPlayer
	query := `
		select attendeeid, tourney, player, players.name, attendees.name, standing
		from attendees
		left outer join players on player = playerid
		where attendeeid = $1
	`
	err = DB.QueryRow(query, id).Scan(&a.AttendeeID, &a.Tourney, &np.PlayerID, &np.Name, &a.Name, &a.Standing)
	a.Player = np.ToPlayer()
	return
}

func (a *Attendee) Update() (err error) {
	var playerid *int
	query := "update players set tourney = $2, player = $3, name = $4, standing = $5 where attendeeid = $1"
	if a.Player != nil {
		playerid = &a.Player.PlayerID
	}
	_, err = DB.Exec(query, a.AttendeeID, a.Tourney, playerid, a.Name, a.Standing)
	return
}

func (a *Attendee) Delete() (err error) {
	_, err = DB.Exec("delete from attendees where attendeeid = $1", a.AttendeeID)
	return
}
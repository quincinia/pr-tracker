// db defines functions for working with our database.
package models

import "database/sql"

// Method 1 from https://www.alexedwards.net/blog/organising-database-access
var DB *sql.DB

// Do this in main instead
func init() {
	var err error
	DB, err = sql.Open("postgres", "user=jacob dbname=gwp_chapter6 password=test sslmode=disable")
	if err != nil {
		panic(err)
	}
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
	
	err = stmt.QueryRow(typeid, t.Name, t.NumEntrants, t.UniquePlacings, t.BracketReset, tierid).Scan(&t.TourneyID)
	return
}

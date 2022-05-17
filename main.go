package main

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var (
	SMASHGG_KEY   = os.Getenv("SMASHGG_KEY")
	CHALLONGE_KEY = os.Getenv("CHALLONGE_KEY")
	DB_USERNAME   = os.Getenv("DB_USERNAME")
	DB_PASSWORD   = os.Getenv("DB_PASSWORD")
)

func main() {
	fmt.Printf("%v %T\n", SMASHGG_KEY, SMASHGG_KEY)
	fmt.Printf("%v %T\n", CHALLONGE_KEY, CHALLONGE_KEY)
	fmt.Printf("%v %T\n", DB_USERNAME, DB_USERNAME)
	fmt.Printf("%v %T\n", DB_PASSWORD, DB_PASSWORD)
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"pr-tracker/controllers/attendees"
	"pr-tracker/controllers/players"
	"pr-tracker/controllers/tournaments"
	"pr-tracker/templates"

	_ "github.com/joho/godotenv/autoload"
	"github.com/julienschmidt/httprouter"
)

var (
	SMASHGG_KEY   = os.Getenv("SMASHGG_KEY")
	CHALLONGE_KEY = os.Getenv("CHALLONGE_KEY")
	DB_USERNAME   = os.Getenv("DB_USERNAME")
	DB_PASSWORD   = os.Getenv("DB_PASSWORD")
)

func main() {
	// fmt.Printf("%v %T\n", SMASHGG_KEY, SMASHGG_KEY)
	// fmt.Printf("%v %T\n", CHALLONGE_KEY, CHALLONGE_KEY)
	// fmt.Printf("%v %T\n", DB_USERNAME, DB_USERNAME)
	// fmt.Printf("%v %T\n", DB_PASSWORD, DB_PASSWORD)

	// flag.Parse()
	// args := flag.Args()
	// input, err := url.Parse(args[0])
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	attendeesRouter := httprouter.New()
	playersRouter := httprouter.New()
	tournamentsRouter := httprouter.New()
	fs := http.FileServer(http.Dir("./static"))

	app := http.NewServeMux()

	attendees.AddRoutes("/api/attendees", attendeesRouter)
	players.AddRoutes("/api/players", playersRouter)
	tournaments.AddRoutes("", tournamentsRouter)

	// http.StripPrefix can also work here if you wish to reuse the same handler on different routes
	app.Handle("/api/attendees/", attendeesRouter)
	app.Handle("/api/players/", playersRouter)
	app.Handle("/api/tournaments/", http.StripPrefix("/api/tournaments", tournamentsRouter))
	app.Handle("/static/", http.StripPrefix("/static/", fs))

	site := httprouter.New()
	site.HandlerFunc("GET", "/", templates.RenderTable)
	site.HandlerFunc("GET", "/tournaments/", templates.RenderTourneySelect)
	site.GET("/tournaments/:id", templates.RenderTourneyView)
	site.HandlerFunc("GET", "/players/", templates.RenderPlayerSelect)
	site.GET("/players/:id", templates.RenderPlayerView)

	app.Handle("/", site)

	fmt.Println("Serving on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", app))
}

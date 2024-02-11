package main

import (
	"log"
	"net/http"

	"github.com/andremfp/poker-app"
)

const dbFileName = "game.db.json"

func main() {
	store, close, err := poker.FsPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	game := poker.NewTexasHoldem(store, poker.BlindAlerterFunc(poker.Alerter))

	server, err := poker.NewPlayerServer(store, game)
	if err != nil {
		log.Fatal("problem creating player server", err)
	}

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000, %v", err)
	}
}

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/andremfp/poker-app"
)

const dbFileName = "game.db.json"

func main() {
	store, close, err := poker.FsPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	fmt.Println("Let's play poker")
	fmt.Println("Type '{Name} wins' to record a win")
	game := poker.NewTexasHoldem(store, poker.BlindAlerterFunc(poker.Alerter))
	poker.NewCLI(os.Stdin, os.Stdout, game).PlayPoker()
}

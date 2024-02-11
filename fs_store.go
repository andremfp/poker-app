package poker

import (
	"encoding/json"
	"fmt"
	"os"
)

type FsPlayerStore struct {
	database *json.Encoder
	league   League
}

// only read from disk once
func NewFsPlayerStore(file *os.File) (*FsPlayerStore, error) {

	err := initializePlayerDbFile(file)
	if err != nil {
		return nil, fmt.Errorf("could not initialize player db file, %v", err)
	}

	league, err := NewLeague(file)
	if err != nil {
		return nil, fmt.Errorf("could not load player store form file %s, %v", file.Name(), err)
	}

	return &FsPlayerStore{
		// using the tape type, allows to have a custom Write function
		database: json.NewEncoder(&tape{file}),
		league:   league,
	}, nil
}

func FsPlayerStoreFromFile(path string) (*FsPlayerStore, func(), error) {
	// create db file. Rw or create if does not exist.
	db, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return nil, nil, fmt.Errorf("problem opening %s %v", path, err)
	}

	closeFunc := func() {
		db.Close()
	}

	store, err := NewFsPlayerStore(db)

	if err != nil {
		return nil, nil, fmt.Errorf("problem creating file system player store, %v ", err)
	}

	return store, closeFunc, nil
}

func (f *FsPlayerStore) GetLeague() League {
	return f.league
}

func (f *FsPlayerStore) GetPlayerScore(playerName string) int {
	player := f.league.Find(playerName)
	if player != nil {
		return player.Wins
	}

	return 0
}

func (f *FsPlayerStore) RecordWin(playerName string) {
	player := f.league.Find(playerName)
	if player != nil {
		player.Wins++
	} else {
		f.league = append(f.league, Player{playerName, 1})
	}

	f.database.Encode(f.league)
}

func initializePlayerDbFile(file *os.File) error {
	// always start from the beginning of the reader
	file.Seek(0, 0)

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("could not get file info from file %s, %v", file.Name(), err)
	}

	// if exists but is empty, write empty json and seek back to beginning
	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, 0)
	}

	return nil
}

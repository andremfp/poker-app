package poker

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

type League []Player

func (l League) Find(playerName string) *Player {
	for i, player := range l {
		if player.Name == playerName {
			return &l[i]
		}
	}

	return nil
}

func NewLeague(reader io.Reader) (League, error) {
	var league []Player
	err := json.NewDecoder(reader).Decode(&league)
	if err != nil {
		err = fmt.Errorf("unable to parse league, %v", err)
	}

	sort.Slice(league, func(i, j int) bool {
		return league[i].Wins > league[j].Wins
	})
	return league, err
}

package poker

import "sync"

/*
This store is not being used.
It was used to help build the server code with a dummy db store.
*/

type InMemoryPlayerStore struct {
	lock   sync.RWMutex
	scores map[string]int
}

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{sync.RWMutex{}, map[string]int{}}
}

func (s *InMemoryPlayerStore) GetPlayerScore(playerName string) int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.scores[playerName]
}

func (s *InMemoryPlayerStore) RecordWin(playerName string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.scores[playerName]++
}

func (s *InMemoryPlayerStore) GetLeague() League {
	// pre allocate slice size for efficiency
	league := make([]Player, len(s.scores))
	i := 0
	for name, wins := range s.scores {
		league[i] = Player{name, wins}
		i++
	}

	return league
}

package main

import (
	"fmt"
	"time"

	"kvstore/kvstore"
)

func main() {
	s, _ := kvstore.NewKVStore(10 * time.Second)
	defer s.Close()

	footballPlayers := []string{"Ronaldo", "Messi", "Neymar", "Mbappe", "Salah", "Kane", "Lewandowski", "De Bruyne", "Modric", "Van Dijk"}
	bytes := [][]byte{
		{'C', 'R', '7'},
		{'L', 'M', '1', '0'},
		{'N', 'J', '1', '0'},
		{'K', 'M', '7'},
		{'M', 'S', '1', '1'},
		{'H', 'K', '1', '0'},
		{'R', 'L', '9'},
		{'D', 'B', '1', '0'},
		{'L', 'M', '1', '0'},
		{'V', 'D', '1', '0'},
	}

	for i, player := range footballPlayers {
		s.Set(player, bytes[i], 5*time.Second)
		time.Sleep(2 * time.Second)
		fmt.Println(s.Get(player))
	}

	fmt.Println("Final store state:", s)
}

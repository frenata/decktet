package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/frenata/gaga/decktet"
)

func oneGame() int {
	deck := decktet.NewDecktet(decktet.BasicDeck)
	player := decktet.NewAdamanPlayer()

	deck.Shuffle(-1)

	score := player.Play(deck)
	return score
}

func runStats(runs int) {
	stats := make(map[string]int)
	var total int

	for i := 0; i < runs; i++ {
		score := oneGame()
		total += score
		switch {
		case score == 0:
			stats["total loss"]++
		case score < 70:
			stats["loss"]++
		case score >= 70:
			stats["win"]++
		default:
			stats["error"]++
		}
	}
	average := total / runs

	fmt.Println("Average:", average)
	fmt.Println(stats)
}

func main() {
	var runs int = 1

	if len(os.Args) > 1 {
		runs, _ = strconv.Atoi(os.Args[1])
	}

	runStats(runs)
}

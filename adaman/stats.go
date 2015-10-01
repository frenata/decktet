// Adaman AI Game, run x of games and run stats.
package main

import (
	"fmt"
	"os"
	"strconv"

	//"github.com/frenata/decktet"
)

func oneGame(player *AdamanPlayer) int {
	player.Shuffle(1)
	score := player.Play()

	//fmt.Println(player)
	//fmt.Println(player.discard)
	//fmt.Println(len(player.Cards()))
	player.cleanup()
	//fmt.Printf("cards: %s\ndiscards: %s\n", player.Cards(), player.Discards())
	player.Shuffle(1)
	//fmt.Printf("cards: %s\ndiscards: %s\n", player.Cards(), player.Discards())
	//fmt.Println(len(player.Cards()))
	return score
}

func runStats(player *AdamanPlayer, runs int) {
	stats := make(map[string]int)
	var total, highscore int

	for i := 0; i < runs; i++ {
		score := oneGame(player)
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
		if score > highscore {
			highscore = score
		}
	}
	average := total / runs
	winPer := float64(stats["win"]) / float64(runs) * 100
	lossPer := float64(stats["loss"]) / float64(runs) * 100
	totalLossPer := float64(stats["total loss"]) / float64(runs) * 100

	fmt.Println("Average:", average)
	fmt.Println("Win% :", winPer)
	fmt.Println("Loss% :", lossPer)
	fmt.Println("Total Loss% :", totalLossPer)
	fmt.Println("High Score:", highscore)
	fmt.Println(stats)
}

func main() {
	var runs int = 1

	if len(os.Args) > 1 {
		runs, _ = strconv.Atoi(os.Args[1])
	}

	player := NewAdamanPlayer()
	runStats(player, runs)
}

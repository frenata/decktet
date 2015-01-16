package main

import (
	"fmt"

	"github.com/frenata/gaga/decktet"
)

func main() {
	//ranks := decktet.BasicRanks
	//fmt.Println(ranks)
	//fmt.Println(decktet.Ace < ranks[1])
	//fmt.Println(decktet.Ace + decktet.Two)
	d := decktet.NewDecktet(decktet.BasicDeck)
	d.Shuffle(-1)

	fmt.Printf("%v\n", d.Shuffled[0])
}

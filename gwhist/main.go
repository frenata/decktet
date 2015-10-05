package main

import (
	"fmt"

	"github.com/frenata/decktet"
)

// setup
// 	if depleted, shuffle
// 	deal
// 	bid
// 	play 7 hands
// check bid
// check win

var round int = 1

func main() {
	fmt.Println("Gongor Whist!")

	g := New(newDummy("dummy"), newAi("tommy"))

	g.cards.Seed(2)

	playR(g)
	fmt.Println(len(g.cards.Cards()), len(g.cards.Discards()))
	playR(g)
	fmt.Println(len(g.cards.Cards()), len(g.cards.Discards()))
	playR(g)
	fmt.Println(len(g.cards.Cards()), len(g.cards.Discards()))
}

func playR(g *game) {
	deal(g)
	score := 0
	for i := 0; i < 7; i++ {
		if hand(g) {
			score++
		}
	}
	fmt.Printf("Won %d tricks.\n", score)
	round++
}

func deal(g *game) {
	fmt.Printf("Round %d!\n", round)
	if len(g.cards.Cards()) == 2 || round == 1 {
		fmt.Println("SHUFFLING!")
		g.shuffle()
	}
	g.deal()
	status(g)
}

func status(g *game) {
	fmt.Println(g.dummy)
	fmt.Println(g.player)
}

func hand(g *game) bool {
	dc := g.dummy.Play()
	fmt.Printf("%s plays %s\n", g.dummy.name, dc)

	ace := g.trump()
	//fmt.Printf("Trump is %s\n", ace)

	pc := g.player.Play(g, dc)
	fmt.Printf("%s plays %s\n", g.player.Name(), pc)

	defer g.cards.Discard(dc, pc)
	if testWin(dc, pc, ace) {
		fmt.Println("player wins")
		return true
	}

	return false
}

func testWin(dc, pc, ace *decktet.DecktetCard) bool {
	trump := ace.Suits()[0]
	switch {
	case dc.HasSuit(trump) && !pc.HasSuit(trump):
		return false
	case pc.HasSuit(trump) && !dc.HasSuit(trump):
		return true
	default: // both have trump or neither does
		return pc.Rank() > dc.Rank()
	}
}

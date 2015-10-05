// A CLI interface for Gongor Whist, playable by input or by ai (with stats)
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

	// create some players
	g := newGame(newDummy("dummy"), newAi("tommy"))

	// for testing, can set a random seed here
	g.cards.Seed(-1)

	// full round not yet implemented, for now just play 3 rounds and output results
	/*playR(g)
	fmt.Println(len(g.cards.Cards()), len(g.cards.Discards()))
	playR(g)
	fmt.Println(len(g.cards.Cards()), len(g.cards.Discards()))
	playR(g)
	fmt.Println(len(g.cards.Cards()), len(g.cards.Discards()))*/

	wins := 0
	games := 10000
	for i := 0; i < games; i++ {
		if playR(g) {
			wins++
		}
	}

	fmt.Printf("Played %d games, won %d games.\nWinPer is %f", games, wins, float64(wins)/float64(games))
}

// play a round
func playR(g *game) bool {
	deal(g)
	fmt.Printf("Bid is %d\n", g.player.(*ai).bestbid(g))
	fmt.Printf("Trump is %s\n", g.trump().Suits()[0])
	g.player.(*ai).score = 0

	score := 0
	for i := 0; i < 7; i++ {
		if hand(g) {
			score++
		}
	}
	fmt.Printf("Won %d tricks.\n", score)
	round++
	return g.player.(*ai).score == g.player.(*ai).bid
}

// deal out cards, if the deck is depleted, first shuffle the discards back in
func deal(g *game) {
	fmt.Printf("Round %d!\n", round)
	if len(g.cards.Cards()) == 2 || round == 1 {
		fmt.Println("SHUFFLING!")
		g.shuffle()
	}
	g.deal()
	status(g)
}

// print the hands of the players
func status(g *game) {
	fmt.Println(g.dummy)
	fmt.Println(g.player)
}

// play a hand
func hand(g *game) bool {
	dc := g.dummy.play()
	fmt.Printf("%s plays %s\n", g.dummy.name, dc)

	ace := g.trump()
	//fmt.Printf("Trump is %s\n", ace)

	pc := g.player.play(g, dc)
	fmt.Printf("%s plays %s\n", g.player.Name(), pc)

	defer g.cards.Discard(dc, pc)
	if testWin(dc, pc, ace) {
		g.player.(*ai).score++
		fmt.Printf("player wins, score is %d\n", g.player.(*ai).score)
		return true
	}

	return false
}

// did the player win the hand?
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

// FollowSuit tests whether a given hand of cards *can* follow either suit on a card that has been led.
func FollowSuit(card *decktet.DecktetCard, hand ...*decktet.DecktetCard) bool {
	for _, s := range card.Suits() {
		for _, c := range hand {
			if c.HasSuit(s) {
				return true
			}
		}
	}
	return false
}

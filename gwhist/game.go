package main

import (
	"fmt"

	"github.com/frenata/deck"
	"github.com/frenata/decktet"
)

type game struct {
	aces   deck.Deck
	cards  deck.Deck
	dummy  deck.Player
	player deck.Player
	bid    [8]bool
}

func init() {
	cards := decktet.BasicDeck().Cards()
	d := make([]*decktet.DecktetCard, 30)
	a := make([]*decktet.DecktetCard, 6)
	ai, di := 0, 0

	for _, c2 := range cards {
		c := c2.(*decktet.DecktetCard)

		if c.Rank() == decktet.Ace {
			a[ai] = c
			ai++
		} else { // c.Rank != Ace
			d[di] = c
			di++
		}
	}
	Deck := decktet.NewDecktet(d)
	Aces := decktet.NewDecktet(a)

	fmt.Printf("Deck:\n%s\n\nAces:\n%s\n\n", Deck, Aces)
}

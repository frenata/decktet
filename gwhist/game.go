package main

import (
	"github.com/frenata/deck"
	"github.com/frenata/decktet"
)

var (
	gdeck *deck.Deck
	gaces *deck.Deck
)

// game status struct, holds the two decks, a dummy player, the actual player, and the bid status
type game struct {
	aces  *deck.Deck
	cards *deck.Deck
	*dummy
	player
	bid [8]bool // a true bid means it has been bid and made previously
}

// split up the standard decktet on import
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
	gdeck = decktet.NewDecktet(d)
	gaces = decktet.NewDecktet(a)

	//fmt.Printf("Deck:\n%s\n\nAces:\n%s\n\n", Deck, Aces)
}

// start a new game of Gongor Whist
func newGame(d dummy, p player) *game {
	g := &game{aces: gaces, cards: gdeck}
	g.dummy = &d
	g.player = p
	return g
}

// when shuffling, shuffle both decks
func (g *game) shuffle() {
	g.aces.Shuffle()
	g.cards.Shuffle()
}

// deal out 7 cards to each player
func (g *game) deal() {
	for i := 0; i < 7; i++ {
		g.cards.Deal(g.dummy)
		g.cards.Deal(g.player)
	}
}

// trump is always the top card on the aces stack
func (g *game) trump() *decktet.DecktetCard {
	return g.aces.Cards()[0].(*decktet.DecktetCard)
}

package main

import (
	"github.com/frenata/deck"
	"github.com/frenata/decktet"
)

var (
	gdeck *deck.Deck
	gaces *deck.Deck
)

type game struct {
	aces  *deck.Deck
	cards *deck.Deck
	*dummy
	player
	bid [8]bool
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
	gdeck = decktet.NewDecktet(d)
	gaces = decktet.NewDecktet(a)

	//fmt.Printf("Deck:\n%s\n\nAces:\n%s\n\n", Deck, Aces)
}

func New(d dummy, p player) *game {
	g := &game{aces: gaces, cards: gdeck}
	g.dummy = &d
	g.player = p
	return g
}

func (g *game) shuffle() {
	g.aces.Shuffle()
	g.cards.Shuffle()
}

func (g *game) deal() {
	for i := 0; i < 7; i++ {
		g.cards.Deal(g.dummy)
		g.cards.Deal(g.player)
	}
}

func (g *game) trump() *decktet.DecktetCard {
	return g.aces.Cards()[0].(*decktet.DecktetCard)
}

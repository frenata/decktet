package main

import (
	"math/rand"

	"github.com/frenata/deck"
	"github.com/frenata/decktet"
)

// a dummy is the dummy hand in gongor whist, simply holding a stack of cards and mindlessly flipping them over
type dummy struct {
	name  string
	cards []*decktet.DecktetCard
}

// implementing the deck.Player interface
func (d *dummy) String() string      { return d.name + ": " + decktet.ShortPrint(d.cards) }
func (d *dummy) AddCard(c deck.Card) { d.cards = append(d.cards, c.(*decktet.DecktetCard)) }

// play dictates how the dummy will play cards
// the dummy plays strictly randomly.
func (d *dummy) play() *decktet.DecktetCard {
	i := rand.Intn(len(d.cards))
	c := d.cards[i]

	d.cards = append(d.cards[:i], d.cards[i+1:]...)
	return c
}

// newDummy makes a named dummy player
func newDummy(name string) dummy {
	d := dummy{name: name}

	d.cards = make([]*decktet.DecktetCard, 0, 7)
	return d
}

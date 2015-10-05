package main

import (
	"fmt"
	"math/rand"

	"github.com/frenata/deck"
	"github.com/frenata/decktet"
)

// implement various player for Gongor Whist
// - Dummy
// - AIPlayer
// - Human

// A player is a local interface that defines structs able to 'play' the game of gongor whist
type player interface {
	deck.Player
	play(g *game, dc *decktet.DecktetCard) *decktet.DecktetCard
	Name() string
}

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

// ai is a self-playing gongor whist player
type ai struct {
	name  string
	cards []*decktet.DecktetCard
}

// implementing the deck.Player interface
func (a *ai) String() string      { return a.name + ": " + decktet.ShortPrint(a.cards) }
func (a *ai) Name() string        { return a.name }
func (a *ai) AddCard(c deck.Card) { a.cards = append(a.cards, c.(*decktet.DecktetCard)) }

// play dicates how the ai player will select cards
// heavily in development
func (a *ai) play(g *game, dc *decktet.DecktetCard) *decktet.DecktetCard {
	var i int
	if a.canFollow(dc) {
		for yes := false; !yes; {
			i = a.randCard()
			yes = FollowSuit(dc, a.cards[i])
		}
	} else {
		fmt.Println("can't follow suit")
		i = a.randCard()
	}

	c := a.cards[i]
	a.cards = append(a.cards[:i], a.cards[i+1:]...)
	return c
}

// can the player follow suit from the led card?
func (a *ai) canFollow(dc *decktet.DecktetCard) bool {
	return FollowSuit(dc, a.cards...)
}

// select a random card from hand
func (a *ai) randCard() int {
	return rand.Intn(len(a.cards))
}

// create a new named ai player
func newAi(name string) player {
	a := &ai{name: name}

	a.cards = make([]*decktet.DecktetCard, 0, 7)
	return a
}

// a human cli player
type human struct {
	name  string
	cards []*decktet.DecktetCard
}

// implementing deck.Player
func (h *human) String() string      { return h.name }
func (h *human) Name() string        { return h.name }
func (h *human) AddCard(c deck.Card) { h.cards = append(h.cards, c.(*decktet.DecktetCard)) }

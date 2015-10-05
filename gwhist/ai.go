package main

import (
	"fmt"
	"math/rand"

	"github.com/frenata/deck"
	"github.com/frenata/decktet"
)

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

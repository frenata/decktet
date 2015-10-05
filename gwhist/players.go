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

type player interface {
	deck.Player
	Play(g *game, dc *decktet.DecktetCard) *decktet.DecktetCard
	Name() string
}

type dummy struct {
	name  string
	cards []*decktet.DecktetCard
}

func (d *dummy) String() string      { return d.name + ": " + decktet.ShortPrint(d.cards) }
func (d *dummy) AddCard(c deck.Card) { d.cards = append(d.cards, c.(*decktet.DecktetCard)) }
func (d *dummy) Play() *decktet.DecktetCard {
	i := rand.Intn(len(d.cards))
	c := d.cards[i]

	d.cards = append(d.cards[:i], d.cards[i+1:]...)
	return c
}
func newDummy(name string) dummy {
	d := dummy{name: name}

	d.cards = make([]*decktet.DecktetCard, 0, 7)
	return d
}

type ai struct {
	name  string
	cards []*decktet.DecktetCard
}

func (a *ai) String() string      { return a.name + ": " + decktet.ShortPrint(a.cards) }
func (a *ai) Name() string        { return a.name }
func (a *ai) AddCard(c deck.Card) { a.cards = append(a.cards, c.(*decktet.DecktetCard)) }
func (a *ai) Play(g *game, dc *decktet.DecktetCard) *decktet.DecktetCard {
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

func (a *ai) canFollow(dc *decktet.DecktetCard) bool {
	return FollowSuit(dc, a.cards...)
}

func (a *ai) randCard() int {
	return rand.Intn(len(a.cards))
}

func newAi(name string) player {
	a := &ai{name: name}

	a.cards = make([]*decktet.DecktetCard, 0, 7)
	return a
}

type human struct {
	name  string
	cards []*decktet.DecktetCard
}

func (h *human) String() string      { return h.name }
func (h *human) Name() string        { return h.name }
func (h *human) AddCard(c deck.Card) { h.cards = append(h.cards, c.(*decktet.DecktetCard)) }

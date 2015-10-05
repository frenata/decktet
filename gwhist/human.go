package main

import (
	"github.com/frenata/deck"
	"github.com/frenata/decktet"
)

// a human cli player
type human struct {
	name  string
	cards []*decktet.DecktetCard
}

// implementing deck.Player
func (h *human) String() string      { return h.name }
func (h *human) Name() string        { return h.name }
func (h *human) AddCard(c deck.Card) { h.cards = append(h.cards, c.(*decktet.DecktetCard)) }

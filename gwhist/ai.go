package main

import (
	"errors"
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

// create a new named ai player
func newAi(name string) player {
	a := &ai{name: name}

	a.cards = make([]*decktet.DecktetCard, 0, 7)
	return a
}

// implementing the deck.Player interface
func (a *ai) String() string      { return a.name + ": " + decktet.ShortPrint(a.cards) }
func (a *ai) Name() string        { return a.name }
func (a *ai) AddCard(c deck.Card) { a.cards = append(a.cards, c.(*decktet.DecktetCard)) }

// play dicates how the ai player will select cards
// heavily in development
func (a *ai) play(g *game, dc *decktet.DecktetCard) *decktet.DecktetCard {
	/*var i int
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
	*/

	plays := a.legalPlays(dc)
	win, loss := a.findWinners(dc, g.trump(), plays)
	if len(win) > 1 {
		fmt.Println(win)
	}
	c := decktet.Min(win)
	if c == nil {
		c = decktet.Min(loss)
	}

	if a.pop(c) {
		return c
	}

	return nil
}

func (a *ai) pop(c *decktet.DecktetCard) bool {
	for i, v := range a.cards {
		if c == v {
			//fmt.Println(s)
			a.cards = append(a.cards[:i], a.cards[i+1:]...)
			return true
		}
	}
	return false
}

// can the player follow suit from the led card?
func (a *ai) canFollow(dc *decktet.DecktetCard) bool {
	return FollowSuit(dc, a.cards...)
}

// select a random card from hand
func randCard(cards []*decktet.DecktetCard) (*decktet.DecktetCard, error) {
	if len(cards) == 0 {
		return nil, errors.New("no card available")
	}
	return cards[rand.Intn(len(cards))], nil
}

// Functions for choosing which card to play

// finds all legally playable cards
func (a *ai) legalPlays(dc *decktet.DecktetCard) (results []*decktet.DecktetCard) {
	if a.canFollow(dc) {
		//fmt.Printf("before checking legal plays, a.cards\n%s\n", a.cards)
		for _, c := range a.cards {
			if decktet.SuitMatch(dc, c) {
				results = append(results, c)
			}
		}
		return results
	}
	return a.cards
}

// outputs a set of winning cards and losing cards
func (a *ai) findWinners(dc, ace *decktet.DecktetCard, cards []*decktet.DecktetCard) (win, loss []*decktet.DecktetCard) {
	for _, c := range cards {
		if testWin(dc, c, ace) {
			win = append(win, c)
		}
		loss = append(loss, c)
	}
	return win, loss
}

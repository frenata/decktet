package main

import (
	"errors"
	"fmt"
	"math"
	"math/rand"

	"github.com/frenata/deck"
	"github.com/frenata/decktet"
)

// ai is a self-playing gongor whist player
type ai struct {
	name  string
	cards []*decktet.DecktetCard
	bid   int
	score int
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

	var c *decktet.DecktetCard
	tilt := a.tilt(g.trump())
	if tilt < 0 {
		fmt.Printf("want to win! tilt is %f\n", tilt)
		c = decktet.Min(win)
		if c == nil && tilt < -0.25 {
			fmt.Printf("New Trump: %s\n", g.flipTrump())
			plays = a.legalPlays(dc)
			win, loss = a.findWinners(dc, g.trump(), plays)
			c = decktet.Min(win)
		}
		if len(a.cards) <= a.need()+1 {
			//if true {
			for c == nil {
				fmt.Println("still no win")
				newt := g.flipTrump()
				if newt == nil {
					break
				}
				fmt.Printf("New Trump: %s\n", newt)
				plays = a.legalPlays(dc)
				win, loss = a.findWinners(dc, g.trump(), plays)
				c = decktet.Min(win)
			}
		}
		if c == nil {
			c = decktet.Min(loss)
		}
	} else {
		fmt.Printf("want to lose! tilt is %f\n", tilt)
		c = decktet.Max(loss)
		if c == nil && tilt > 0.15 {
			fmt.Printf("New Trump: %s\n", g.flipTrump())
			plays = a.legalPlays(dc)
			win, loss = a.findWinners(dc, g.trump(), plays)
			c = decktet.Max(loss)
		}
		if a.score == a.bid {
			for c == nil {
				fmt.Println("still no loss")
				newt := g.flipTrump()
				if newt == nil {
					break
				}
				fmt.Printf("New Trump: %s\n", newt)
				plays = a.legalPlays(dc)
				win, loss = a.findWinners(dc, g.trump(), plays)
				c = decktet.Max(loss)
			}
		}
		if c == nil {
			c = decktet.Max(win)
		}
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
		} else {
			loss = append(loss, c)
		}
	}
	return win, loss
}

// Functions for bidding, evalutating

func (a *ai) bestbid(g *game) int {
	var bid int
	guess := a.guess(g.trump())

	c := math.Ceil(guess)
	f := math.Floor(guess)

	if math.Abs(c-guess) < math.Abs(guess-f) {
		bid = int(c)
	} else {
		bid = int(f)
	}

	inc := true
	if bid <= 3 {
		inc = false
	}

	for g.bid[bid] {
		//fmt.Println(g.bid)
		fmt.Println("bid already taken, finding another")
		fmt.Println(bid)
		switch {
		case bid == 7:
			inc = false
			bid--
		case bid == 0:
			inc = true
			bid++
		case inc:
			bid++
			//case !inc:
		default:
			bid--
			/*default:
			panic("can't bid")*/
		}
	}

	a.bid = bid
	return bid
}

func (a *ai) guess(ace *decktet.DecktetCard) (p float64) {
	for _, c := range a.cards {
		var value float64
		switch {
		case c.Rank() == decktet.Crown:
			value = .95
		case c.Rank() == decktet.Nine:
			value = .75
		case c.Rank() == decktet.Eight:
			value = .55
		case c.Rank() == decktet.Seven:
			value = .35
		case c.Rank() == decktet.Six:
			value = .15
		default:
			value = .05
		}

		if decktet.SuitMatch(c, ace) {
			value = value + value
			if value > 1 {
				value = 1
			}
			if value < .25 {
				value = .25
			}
		}
		//		fmt.Printf("%s estimated value is %f\n", c, value)
		p += value
	}
	return p
}

func (a *ai) need() int {
	return a.bid - a.score
}

func (a *ai) tilt(ace *decktet.DecktetCard) float64 {
	g := a.guess(ace)
	return g - float64(a.need())
}

package main

import (
	"fmt"

	"github.com/frenata/gaga"
	. "github.com/frenata/gaga/decktet"
)

type AdamanPlayer struct {
	name      string
	resources []*DecktetCard
	capital   []*DecktetCard
	palace    []*DecktetCard
	discard   []*DecktetCard
}

func NewAdamanPlayer() *AdamanPlayer {
	p := new(AdamanPlayer)
	p.resources = make([]*DecktetCard, 0, 5)
	p.capital = make([]*DecktetCard, 0, 5)
	p.palace = make([]*DecktetCard, 0, 5)
	p.discard = make([]*DecktetCard, 0)
	return p
}

func (p *AdamanPlayer) AddCard(c gaga.Card) {
	dc := c.(*DecktetCard)
	if len(p.resources) < 5 {
		p.addResource(dc)
	} else if len(p.capital) < 5 {
		p.capital = append(p.capital, dc)
	}
	//return errors.New("shouldn't have drawn a card!")
}

func (p *AdamanPlayer) addResource(card *DecktetCard) {
	if isPerson(card) {
		p.palace = append(p.palace, card)
	} else {
		p.resources = append(p.resources, card)
	}
}

func (p *AdamanPlayer) String() string {
	pal := fmt.Sprintf("\nPalace:\n%v\n", printRow(p.palace))
	c := fmt.Sprintf("Capital:\n%v\n", printRow(p.capital))
	r := fmt.Sprintf("Resources:\n%v\n", printRow(p.resources))

	return pal + c + r
}

func (p *AdamanPlayer) dealAll(d *gaga.Deck) {
	for len(p.capital) != 5 || len(p.resources) != 5 {
		d.Deal(p)
	}
}

func printRow(row []*DecktetCard) string {
	var s string
	for _, card := range row {
		if isPerson(card) {
			s += fmt.Sprintf("*%v  ", ShortPrintCard(card))
		} else {
			s += fmt.Sprintf(" %v  ", ShortPrintCard(card))
		}
	}
	return s
}

func isPerson(card *DecktetCard) bool {
	for _, cat := range card.Cats() {
		if cat == Person {
			return true
		}
	}
	return false
}

// func evalutate possible claim
// target card. list of cards, can it be claimed, and what's the overage

// func find combinations
// target card, a big list of cards, figure out all combinations that could apply
// first all cards that match at least one suit
// find all possible combinations (or is it permutations?) of sublist that matches suit,
// then call our evaluate function
// return the "best" / most efficient

// func check targets
// given the gamestate
// find combos and evaluate for all targets
// store the overages into a map of some kind

// func decide
// take the map from check targets
// pick the best card to claim

// various versions of decide
//

// suns, moons

func main() {
	player := NewAdamanPlayer()
	deck := NewDecktet(BasicDeck)

	deck.Shuffle(-1)
	fmt.Println(countSuits(deck.Shuffled, true))
	fmt.Println(countSuits(deck.Shuffled, false))
	player.dealAll(deck)

	fmt.Println(player)
	fmt.Println(countSuits(deck.Shuffled, true))
	fmt.Println(countSuits(deck.Shuffled, false))
}

func rankToInt(c *DecktetCard) int {
	r := c.Rank()
	if r == Crown {
		return 10
	} else {
		return int(r)
	}
}

func countSuits(cards []gaga.Card, onlyPersons bool) map[string]int {
	m := make(map[string]int)
	for _, g := range cards {
		c := g.(*DecktetCard)
		if isPerson(c) || !onlyPersons {
			for _, s := range c.Suits() {
				m[string(s)] += rankToInt(c)
			}
		}
	}
	return m
}

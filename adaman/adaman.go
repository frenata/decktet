package main

import (
	"fmt"

	"github.com/frenata/gaga"
	. "github.com/frenata/gaga/decktet"
)

type AdamanPlayer struct {
	name          string
	resources     []*DecktetCard
	capital       []*DecktetCard
	palace        []*DecktetCard
	discard       []*DecktetCard
	capitalClaims map[*DecktetCard][]*DecktetCard
	palaceClaims  map[*DecktetCard][]*DecktetCard
}

func NewAdamanPlayer() *AdamanPlayer {
	p := new(AdamanPlayer)
	p.resources = make([]*DecktetCard, 0, 5)
	p.capital = make([]*DecktetCard, 0, 5)
	p.palace = make([]*DecktetCard, 0, 5)
	p.discard = make([]*DecktetCard, 0)

	p.capitalClaims = make(map[*DecktetCard][]*DecktetCard)
	p.palaceClaims = make(map[*DecktetCard][]*DecktetCard)
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
func evalClaim(target *DecktetCard, cards []*DecktetCard) int {
	// list should already match suits
	var targetV, listV int
	targetV = rankToInt(target)
	for _, c := range cards {
		listV += rankToInt(c)
	}
	return listV - targetV
}

// func find combinations
// target card, a big list of cards, figure out all combinations that could apply
// first all cards that match at least one suit
// find all possible combinations (or is it permutations?) of sublist that matches suit,
// then call our evaluate function
// return the "best" / most efficient
func findCombos(target *DecktetCard, cards []*DecktetCard) (claim []*DecktetCard) {
	var matching []*DecktetCard
	for _, c := range cards { // only consider cards that match a suit
		if SuitMatch(target, c) {
			matching = append(matching, c)
		}
	}

	combos := gaga.CardCombinations(toCardSlice(matching))

	var best []*DecktetCard
	var lowscore int = 100
	for _, c := range combos {
		//fmt.Println(c)
		dc := toDecktetSlice(c)
		score := evalClaim(target, dc)
		//fmt.Println(score)
		if score < 0 {
			continue
		} else if score < lowscore {
			best = dc
			lowscore = score
		}
	}

	return best
}

// func check targets
// given the gamestate
// find combos and evaluate for all targets
// store the overages into a map of some kind
func (p *AdamanPlayer) checkTargets() {
	for _, c := range p.capital {
		p.capitalClaims[c] = findCombos(c, p.resources)
	}

	for _, c := range p.palace {
		p.palaceClaims[c] = findCombos(c, p.resources)
	}

	//fmt.Println("Capital Claims:\n", p.capitalClaims)
	//fmt.Println("Palace Claims:\n", p.palaceClaims)
}

// func decide
// take the map from check targets
// pick the best card to claim
func (p *AdamanPlayer) decideEfficient() {
	var target *DecktetCard
	var claim []*DecktetCard
	var lowscore int = 100

	p.checkTargets()
	for k, v := range p.capitalClaims {
		score := evalClaim(k, v)
		if score < lowscore && score >= 0 {
			target = k
			claim = v
			lowscore = score
		}
	}

	fmt.Printf("Decided: claim %v with %v for %v overpayment.\n", ShortPrintCard(target), ShortPrint(claim), lowscore)
	p.Play(target, claim)
}

func (p *AdamanPlayer) Play(target *DecktetCard, claim []*DecktetCard) {
	// find where the target is
	// remove the target from capital/palace row
	// remove claim from map
	// remove claim cards from resource row
	// discard all used cards
	// make/add a count of how many cards have been claimed
}

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

	//fmt.Println(evalClaim(player.capital[0], player.resources))
	//fmt.Println(findCombos(player.capital[0], player.resources))
	player.decideEfficient()

	/*
		fmt.Println("combos of 3")
		fmt.Println(gaga.Combination(3))
		fmt.Println("combos of 4")
		fmt.Println(gaga.Combination(4))
		fmt.Println("combos of 5")
		fmt.Println(gaga.Combination(5))

		combos := gaga.CardCombinations(toCardSlice(player.resources))

		for _, c := range combos {
			fmt.Println(ShortPrint(toDecktetSlice(c)))
		}
	*/
}

func toCardSlice(dc []*DecktetCard) (c []gaga.Card) {
	c = make([]gaga.Card, len(dc))
	for i := range dc {
		c[i] = dc[i]
	}
	return c
}
func toDecktetSlice(c []gaga.Card) (dc []*DecktetCard) {
	dc = make([]*DecktetCard, len(c))
	for i := range c {
		dc[i] = c[i].(*DecktetCard)
	}
	return dc
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

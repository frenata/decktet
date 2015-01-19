package adaman

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
		if len(d.Shuffled) == 0 {
			return
		}
		d.Deal(p)
		//fmt.Printf("left in (shuffled) deck: %v\n", len(d.Shuffled))
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
func (p *AdamanPlayer) decideEfficient() bool {
	var target *DecktetCard
	var claim []*DecktetCard
	var lowscore int = 100

	p.checkTargets()
	for k, v := range p.palaceClaims {
		score := evalClaim(k, v)
		if score < lowscore && score >= 0 {
			target = k
			claim = v
			lowscore = score
		}
	}
	for k, v := range p.capitalClaims {
		score := evalClaim(k, v)
		if score < lowscore && score >= 0 {
			target = k
			claim = v
			lowscore = score
		}
	}

	if target == nil {
		return false
	} else {
		fmt.Printf("Decided: claim %v with %v for %v overpayment.\n", ShortPrintCard(target), ShortPrint(claim), lowscore)
		p.claim(target, claim)
		return true
	}
}

func (p *AdamanPlayer) pop(c *DecktetCard, ps *[]*DecktetCard) bool {
	s := *ps
	for i, v := range s {
		if c == v {
			s = append(s[:i], s[i+1:]...)
			return true
		}
	}
	return false
}

func (p *AdamanPlayer) claim(target *DecktetCard, claim []*DecktetCard) {
	// remove the target from capital/palace row
	for i, v := range p.palace {
		if target == v {
			p.palace = append(p.palace[:i], p.palace[i+1:]...)
			delete(p.palaceClaims, target) // clean map
			break
		}
	}

	for i, v := range p.capital {
		if target == v {
			p.capital = append(p.capital[:i], p.capital[i+1:]...)
			delete(p.capitalClaims, target)
			break
		}
	}

	// remove claim cards from resource row
	for _, c := range claim {
		for i, v := range p.resources {
			if c == v {
				p.resources = append(p.resources[:i], p.resources[i+1:]...)
				break
			}
		}
	}

	// discard all used cards
	// make/add a count of how many cards have been claimed
	if isPerson(target) {
		p.discard = append(p.discard, target)
	}

}

func (p *AdamanPlayer) isGameOver() (result string) {
	if len(p.palace) > 5 {
		fmt.Println("palace overfull")
		return "total loss"
	} else if len(p.discard) == 11 {
		fmt.Println("all people captured")
		return "win"
	} else {
		return ""
	}
}

func (p *AdamanPlayer) Play(deck *gaga.Deck) int {
	var result string
	var round int
	for result == "" {
		round++
		fmt.Printf("Round %v!\n", round)
		if len(deck.Shuffled) > 0 {
			p.dealAll(deck)
		}
		fmt.Println(p)
		//fmt.Println(countSuits(deck.Shuffled, true))
		//fmt.Println(countSuits(deck.Shuffled, false))
		result = p.isGameOver()
		if result != "" {
			break
		}
		if !p.decideEfficient() {
			result = "loss"
			fmt.Println("no more decisions possible")
		}
	}

	fmt.Println("Final table", p)
	fmt.Println("Captured: ", ShortPrint(p.discard))
	return p.score(result)
}

func (p *AdamanPlayer) score(result string) int {
	var totalP, totalR int

	for _, c := range p.discard {
		totalP += rankToInt(c)
	}

	for _, c := range p.resources {
		totalR += rankToInt(c)
	}

	switch result {
	case "total loss":
		return 0
	case "loss":
		return totalP
	case "win":
		return totalP + totalR
	default:
		panic("bad end condition")
		return -1
	}
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

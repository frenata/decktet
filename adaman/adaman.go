package main

import (
	"fmt"

	"github.com/frenata/deck"
	. "github.com/frenata/decktet"
)

type AdamanPlayer struct {
	name string
	*deck.Deck
	resources     []*DecktetCard
	capital       []*DecktetCard
	palace        []*DecktetCard
	discard       []*DecktetCard
	capitalClaims map[*DecktetCard][]*DecktetCard
	palaceClaims  map[*DecktetCard][]*DecktetCard
}

func NewAdamanPlayer() *AdamanPlayer {
	p := new(AdamanPlayer)
	p.Deck = BasicDeck()

	p.start()

	return p
}

func (p *AdamanPlayer) start() {
	p.resources = make([]*DecktetCard, 0, 5)
	p.capital = make([]*DecktetCard, 0, 5)
	p.palace = make([]*DecktetCard, 0, 5)
	p.discard = make([]*DecktetCard, 0)
	p.capitalClaims = make(map[*DecktetCard][]*DecktetCard)
	p.palaceClaims = make(map[*DecktetCard][]*DecktetCard)
}

func (p *AdamanPlayer) cleanup() {
	p.Discard(toCardSlice(p.resources)...)
	p.Discard(toCardSlice(p.capital)...)
	p.Discard(toCardSlice(p.palace)...)
	p.Discard(toCardSlice(p.discard)...)
	p.start()
}

//func (p *AdamanPlayer) Shuffle(seed int) {
//	p.Deck.Shuffle(seed)
//}

func (p *AdamanPlayer) AddCard(c deck.Card) {
	dc := c.(*DecktetCard)
	//fmt.Printf("Card dealt: %s\n", dc)
	if len(p.resources) < 5 {
		p.addResource(dc)
	} else if len(p.capital) < 5 {
		p.capital = append(p.capital, dc)
	}
	//return errors.New("shouldn't have drawn a card!")
	//fmt.Println("shouldn't have drawn a card!")
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
	k := fmt.Sprintf("Kills:\n%v\n", printRow(p.discard))

	return pal + c + r + k
}

func (p *AdamanPlayer) dealAll() {
	for len(p.capital) < 5 || len(p.resources) < 5 {
		//for len(p.capital) != 5 || len(p.resources) != 5 {
		if len(p.Cards()) == 0 {
			return
		}
		p.Deal(p)
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
func (p *AdamanPlayer) evalClaim(target *DecktetCard, cards []*DecktetCard) int {
	// list should already match suits
	var targetV, listV int
	targetV = rankToInt(target)
	for _, c := range cards {
		listV += rankToInt(c)
	}
	return listV - targetV
}

func (p *AdamanPlayer) adjClaim(target *DecktetCard, cards []*DecktetCard) float64 {

	var targetI int
	var targetF, listF float64

	targetI = rankToInt(target)
	targetF = float64(targetI) * p.cardValue(target, true)

	for _, c := range cards {
		v := rankToInt(c)
		f := p.cardValue(c, false)
		listF = listF + (float64(v) * f)
	}
	return listF - targetF
}

func (p *AdamanPlayer) cardValue(card *DecktetCard, target bool) (total float64) {
	count := countSuits(p.Cards(), true)
	values := make(map[string]float64)

	for k, v := range count {
		values[k] = float64(v) / 55
	}

	for _, s := range card.Suits() {
		total += values[string(s)]
	}
	/*var bestTarget float64 = 100
	for _, s := range card.Suits() {
		tmp := values[string(s)]
		if target {
			if tmp < bestTarget {
				bestTarget = tmp
			}
		} else {
			if tmp > best {
				best = tmp
			}
		}
	}

	if best == 0 && bestTarget != 100 {
		best = bestTarget
	}*/
	return total
}

// func find combinations
// target card, a big list of cards, figure out all combinations that could apply
// first all cards that match at least one suit
// find all possible combinations (or is it permutations?) of sublist that matches suit,
// then call our evaluate function
// return the "best" / most efficient
func (p *AdamanPlayer) findCombos(target *DecktetCard, cards []*DecktetCard) (claim []*DecktetCard) {
	var matching []*DecktetCard
	for _, c := range cards { // only consider cards that match a suit
		if SuitMatch(target, c) {
			matching = append(matching, c)
		}
	}

	combos := CardCombinations(toCardSlice(matching))

	var best []*DecktetCard
	var lowscore float64 = 100
	for _, c := range combos {
		dc := toDecktetSlice(c)
		// score := p.evalClaim(target,dc)
		score := p.adjClaim(target, dc)
		if p.evalClaim(target, dc) < 0 {
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
	p.capitalClaims = make(map[*DecktetCard][]*DecktetCard)
	p.palaceClaims = make(map[*DecktetCard][]*DecktetCard)

	for _, c := range p.capital {
		if maybe := p.findCombos(c, p.resources); maybe != nil {
			//fmt.Printf("p.target: %s, p.claim %s\n", c, maybe)
			p.capitalClaims[c] = maybe
		}
	}

	for _, c := range p.palace {
		if maybe := p.findCombos(c, p.resources); maybe != nil {
			//fmt.Printf("p.target: %s, p.claim %s\n", c, maybe)
			p.palaceClaims[c] = maybe
		}
	}

}

// func decide
// take the map from check targets
// pick the best card to claim
/*func (p *AdamanPlayer) decideEfficient() bool {
	var target *DecktetCard
	var claim []*DecktetCard
	var lowscore int = 100

	p.checkTargets()
	for k, v := range p.palaceClaims {
		score := p.evalClaim(k, v)
		if score < lowscore && score >= 0 {
			target = k
			claim = v
			lowscore = score
		}
	}
	for k, v := range p.capitalClaims {
		score := p.evalClaim(k, v)
		if score < lowscore && score >= 0 {
			target = k
			claim = v
			lowscore = score
		}
	}

	if target == nil {
		return false
	} else {
		//fmt.Printf("Decided: claim %v with %v for %v overpayment.\n", ShortPrintCard(target), ShortPrint(claim), lowscore)
		p.claim(target, claim)
		return true
	}
}*/

func (p *AdamanPlayer) decideAdj() bool {
	var target *DecktetCard
	var claim []*DecktetCard
	var lowscore float64 = 100

	p.checkTargets()
	for k, v := range p.palaceClaims {
		score := float64(p.evalClaim(k, v))
		// push this to antoher decide function
		if !HasSuit(k, Suns) && !HasSuit(k, Moons) {
			score += 5
		}
		if score < lowscore && score >= 0 {
			target = k
			claim = v
			lowscore = score
		}
	}
	for k, v := range p.capitalClaims {
		score := float64(p.evalClaim(k, v))
		if !HasSuit(k, Suns) && !HasSuit(k, Moons) {
			score += 5
		}
		if score < lowscore && score >= 0 {
			target = k
			claim = v
			lowscore = score
		}
	}

	if target == nil {
		return false
	} else {
		//fmt.Printf("Decided: claim %v with %v for %v overpayment.\n", ShortPrintCard(target), ShortPrint(claim), lowscore)
		p.claim(target, claim)
		return true
	}
}

// dyanmic scoring
// assign value to each suit: x/55 where x = the number of personality "total ranks" left in the shuffle
// card value = both it's suit values added (OR: the most valuable suit?) * rank
// claim value = target value - all cards in claim stack value (= surplus value claimed)
// will have to check this both places (initial check and then decision) or add data to map/struct

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

	failure := true
	// remove claim cards from resource row
	for _, c := range claim {
		for i, v := range p.resources {
			if c == v {
				//fmt.Printf("claim: %s, resource: %s\n", c, v)
				p.resources = append(p.resources[:i], p.resources[i+1:]...)
				//fmt.Printf("resource: %s\n", p.resources)
				failure = false
				break
			}
		}
	}

	if failure {
		fmt.Printf("Faulty Claim! %s\n", claim)
	}

	// discard all used cards
	//TODO
	p.Discard(toCardSlice(claim)...)

	// make/add a count of how many cards have been claimed
	if isPerson(target) {
		p.discard = append(p.discard, target)
	} else { // target is a resource (claimed from capital), add it to resource row
		p.resources = append(p.resources, target)
	}

}

func (p *AdamanPlayer) isGameOver() (result string) {
	if len(p.palace) > 5 {
		//fmt.Println("palace overfull")
		return "total loss"
	} else if len(p.discard) == 11 {
		//fmt.Println("all people captured")
		return "win"
	} else {
		return ""
	}
}

func (p *AdamanPlayer) Play() int {
	var result string
	var round int
	for result == "" {
		round++
		//fmt.Printf("Round %v!\n", round)
		if len(p.Cards()) > 0 {
			p.dealAll()
		}
		//fmt.Println(p)
		//fmt.Printf("cards: %d, discards: %d\n", len(p.Cards()), len(p.Discards()))
		//fmt.Println(countSuits(deck.Shuffled, true))
		//fmt.Println(countSuits(deck.Shuffled, false))
		result = p.isGameOver()
		if result != "" {
			break
		}
		if !p.decideAdj() {
			result = "loss"
			//fmt.Println("no more decisions possible")
		}
	}

	//fmt.Println("Final table", p)
	//fmt.Println("Captured: ", ShortPrint(p.discard))
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

func toCardSlice(dc []*DecktetCard) (c []deck.Card) {
	c = make([]deck.Card, len(dc))
	for i := range dc {
		c[i] = dc[i]
	}
	return c
}
func toDecktetSlice(c []deck.Card) (dc []*DecktetCard) {
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

func countSuits(cards []deck.Card, onlyPersons bool) map[string]int {
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

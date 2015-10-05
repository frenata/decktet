package decktet

import (
	"fmt"
	"strconv"

	"github.com/frenata/deck"
)

/*
type Card interface {
	String() string
}

type Player interface {
	AddCard(c Card)
	String() string
}
*/

type DecktetCard struct {
	name  string
	rank  rank
	suits []suit
	cats  []category
}

func (d *DecktetCard) String() string {
	s := d.rank.String() + " of "

	switch len(d.suits) {
	case 1:
		return s + fmt.Sprint(d.suits[0])
	case 2:
		return s + fmt.Sprintf("%v, %v", d.suits[0], d.suits[1])
	case 3:
		return s + fmt.Sprintf("%v, %v, %v", d.suits[0], d.suits[1], d.suits[2])
	default:
		return "E"
	}
}

func ShortPrint(cards []*DecktetCard) string {
	var s string
	for _, c := range cards {
		s += ShortPrintCard(c) + " "
	}
	return s
}

func ShortPrintCard(card *DecktetCard) string {
	s := shortRank(card.rank)
	for _, suit := range card.suits {
		s += shortSuit(suit)
	}
	if len(card.suits) == 1 {
		s += " "
	}
	return s
}

func SuitMatch(one, two *DecktetCard) bool {
	for _, s1 := range one.suits {
		for _, s2 := range two.suits {
			if s1 == s2 {
				return true
			}
		}
	}
	return false
}

func PopDecktetCard(c *DecktetCard, s []*DecktetCard) bool {
	for i, v := range s {
		if c == v {
			//fmt.Println(s)
			s = append(s[:i], s[i+1:]...)
			return true
		}
	}
	return false
}

func HasSuit(c *DecktetCard, s suit) bool {
	for _, x := range c.Suits() {
		if x == s {
			return true
		}
	}
	return false
}

func (d *DecktetCard) Cats() []category {
	return d.cats
}
func (d *DecktetCard) Suits() []suit {
	return d.suits
}
func (d *DecktetCard) Rank() rank {
	return d.rank
}
func (d *DecktetCard) Name() string {
	return d.name
}

type rank int
type suit string
type category string

func shortRank(r rank) string {
	switch r {
	case Ace:
		return "A"
	case Pawn:
		return "P"
	case Court:
		return "C"
	case Crown:
		return "X"
	default:
		return strconv.Itoa(int(r))
	}
}

func shortSuit(s suit) string {
	switch s {
	case Wyrms:
		return "Y"
	default:
		return string(s[:1])
	}
}

func (r rank) String() string {
	switch r {
	case Ace:
		return "Ace"
	case Two:
		return "Two"
	case Three:
		return "Three"
	case Four:
		return "Four"
	case Five:
		return "Five"
	case Six:
		return "Six"
	case Seven:
		return "Seven"
	case Eight:
		return "Eight"
	case Nine:
		return "Nine"
	case Pawn:
		return "Pawn"
	case Court:
		return "Court"
	case Crown:
		return "Crown"
	default:
		return ""
	}
}

const (
	_        = iota
	Ace rank = iota
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Pawn
	Court
	Crown

	Suns   suit = "Suns"
	Moons  suit = "Moons"
	Waves  suit = "Waves"
	Leaves suit = "Leaves"
	Knots  suit = "Knots"
	Wyrms  suit = "Wyrms"

	Person category = "Person"
	Place  category = "Place"
	Event  category = "Event"
)

var BasicRanks = [10]rank{Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Crown}

func BasicDeck() *deck.Deck {
	bd := *basicDeck
	return &bd
}

var basicDeck = NewDecktet([]*DecktetCard{
	{"", Ace, []suit{Suns}, nil},
	{"", Ace, []suit{Moons}, nil},
	{"", Ace, []suit{Waves}, nil},
	{"", Ace, []suit{Leaves}, nil},
	{"", Ace, []suit{Knots}, nil},
	{"", Ace, []suit{Wyrms}, nil},
	{"Author", Two, []suit{Moons, Knots}, []category{Person}},
	{"Origin", Two, []suit{Waves, Leaves}, []category{Event, Place}},
	{"Desert", Two, []suit{Suns, Wyrms}, []category{Place}},
	{"Journey", Three, []suit{Moons, Waves}, []category{Event}},
	{"Savage", Three, []suit{Leaves, Wyrms}, []category{Person}},
	{"Painter", Three, []suit{Suns, Knots}, []category{Person}},
	{"Mountain", Four, []suit{Moons, Suns}, []category{Place}},
	{"Battle", Four, []suit{Wyrms, Knots}, []category{Event}},
	{"Sailor", Four, []suit{Waves, Leaves}, []category{Person}},
	{"Discovery", Five, []suit{Suns, Waves}, []category{Event}},
	{"Soldier", Five, []suit{Wyrms, Knots}, []category{Person}},
	{"Forest", Five, []suit{Moons, Leaves}, []category{Place}},
	{"Penitent", Six, []suit{Suns, Wyrms}, []category{Person}},
	{"Lunatic", Six, []suit{Moons, Waves}, []category{Person}},
	{"Market", Six, []suit{Leaves, Knots}, []category{Event, Place}},
	{"Cave", Seven, []suit{Waves, Wyrms}, []category{Place}},
	{"Castle", Seven, []suit{Suns, Knots}, []category{Place}},
	{"Chance Meeting", Seven, []suit{Moons, Leaves}, []category{Event}},
	{"Betrayal", Eight, []suit{Wyrms, Knots}, []category{Event}},
	{"Mill", Eight, []suit{Waves, Leaves}, []category{Place}},
	{"Diplomat", Eight, []suit{Moons, Suns}, []category{Person}},
	{"Merchant", Nine, []suit{Leaves, Knots}, []category{Person}},
	{"Darkness", Nine, []suit{Waves, Wyrms}, []category{Place}},
	{"Pact", Nine, []suit{Moons, Suns}, []category{Event}},
	{"Calamity", Crown, []suit{Wyrms}, []category{Event}},
	{"Huntress", Crown, []suit{Moons}, []category{Person}},
	{"Bard", Crown, []suit{Suns}, []category{Person}},
	{"Sea", Crown, []suit{Waves}, []category{Place}},
	{"Windfall", Crown, []suit{Knots}, []category{Event}},
	{"End", Crown, []suit{Leaves}, []category{Event, Place}},
})

func NewDecktet(dc []*DecktetCard) *deck.Deck {
	c := make([]deck.Card, len(dc))

	for i := range dc {
		c[i] = dc[i]
	}

	return deck.New(c)
}

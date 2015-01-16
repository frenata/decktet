package decktet

import (
	"fmt"
	"strconv"

	"github.com/frenata/gaga"
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
	name   string
	rank   rank
	suits  []suit
	cats   []category
	played []gaga.Player
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

var BasicDeck = []*DecktetCard{
	{"", Ace, []suit{Suns}, nil, nil},
	{"", Ace, []suit{Moons}, nil, nil},
	{"", Ace, []suit{Waves}, nil, nil},
	{"", Ace, []suit{Leaves}, nil, nil},
	{"", Ace, []suit{Knots}, nil, nil},
	{"", Ace, []suit{Wyrms}, nil, nil},
	{"Author", Two, []suit{Moons, Knots}, []category{Person}, nil},
	{"Origin", Two, []suit{Waves, Leaves}, []category{Event, Place}, nil},
	{"Desert", Two, []suit{Suns, Wyrms}, []category{Place}, nil},
	{"Journey", Three, []suit{Moons, Waves}, []category{Event}, nil},
	{"Savage", Three, []suit{Leaves, Wyrms}, []category{Person}, nil},
	{"Painter", Three, []suit{Suns, Knots}, []category{Person}, nil},
	{"Mountain", Four, []suit{Moons, Suns}, []category{Place}, nil},
	{"Battle", Four, []suit{Wyrms, Knots}, []category{Event}, nil},
	{"Sailor", Four, []suit{Waves, Leaves}, []category{Person}, nil},
	{"Discovery", Five, []suit{Suns, Waves}, []category{Event}, nil},
	{"Soldier", Five, []suit{Wyrms, Knots}, []category{Person}, nil},
	{"Forest", Five, []suit{Moons, Leaves}, []category{Place}, nil},
	{"Penitent", Six, []suit{Suns, Wyrms}, []category{Person}, nil},
	{"Lunatic", Six, []suit{Moons, Waves}, []category{Person}, nil},
	{"Market", Six, []suit{Leaves, Knots}, []category{Event, Place}, nil},
	{"Cave", Seven, []suit{Waves, Wyrms}, []category{Place}, nil},
	{"Castle", Seven, []suit{Suns, Knots}, []category{Place}, nil},
	{"Chance Meeting", Seven, []suit{Moons, Leaves}, []category{Event}, nil},
	{"Betrayal", Eight, []suit{Wyrms, Knots}, []category{Event}, nil},
	{"Mill", Eight, []suit{Waves, Leaves}, []category{Place}, nil},
	{"Diplomat", Eight, []suit{Moons, Suns}, []category{Person}, nil},
	{"Merchant", Nine, []suit{Leaves, Knots}, []category{Person}, nil},
	{"Darkness", Nine, []suit{Waves, Wyrms}, []category{Place}, nil},
	{"Pact", Nine, []suit{Moons, Suns}, []category{Event}, nil},
	{"Calamity", Crown, []suit{Wyrms}, []category{Event}, nil},
	{"Huntress", Crown, []suit{Moons}, []category{Person}, nil},
	{"Bard", Crown, []suit{Suns}, []category{Person}, nil},
	{"Sea", Crown, []suit{Waves}, []category{Place}, nil},
	{"Windfall", Crown, []suit{Knots}, []category{Event}, nil},
	{"End", Crown, []suit{Leaves}, []category{Event, Place}, nil},
}

func NewDecktet(dc []*DecktetCard) *gaga.Deck {
	c := make([]gaga.Card, len(dc))

	for i := range dc {
		c[i] = dc[i]
	}

	return gaga.NewDeck(c)
}

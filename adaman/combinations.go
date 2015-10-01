package main

import (
	"math"

	"github.com/frenata/deck"
)

// CardCombinations returns a list of all possible combinations of given a slice of Cards.
func CardCombinations(cards []deck.Card) [][]deck.Card {
	var results [][]deck.Card
	set := Combination(len(cards))

	for _, s := range set {
		row := []deck.Card{}
		for _, v := range s {
			row = append(row, cards[v])
		}
		results = append(results, row)
	}
	return results
}

// Combination returns a list of all possible integer combinations given an integer.
func Combination(n int) [][]int {
	var slice []int
	var results [][]int
	var b byte
	for i := 0; i < int(math.Pow(2, float64(n))); i++ {
		b = byte(i)
		slice = []int{}
		for j := 0; j < n; j++ { //int(math.Pow(2, float64(n))); j++ {
			if b>>uint(j)&1 == 1 {
				slice = append(slice, j)
			}
		}
		if len(slice) != 0 {
			results = append(results, slice)
		}
	}
	return results
}

package main

import (
	"log"

	"github.com/agnivade/levenshtein"
)

func search(needle string, haystack []string) {
	res := make(map[string]int, len(haystack))
	min := 100 // stupid reasonable min
	for _, item := range haystack {
		dist := levenshtein.ComputeDistance(needle, item)
		if dist < min {
			min = dist
		}
		res[item] = dist
	}

	for result, distance := range res {
		if distance == min {
			log.Printf("%-10s: %d", result, min)
		}
	}
}

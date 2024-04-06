package main

import (
	"fmt"
	"time"
)

func perft(depth int) int {
	if depth == 0 {
		return 1
	}
	moves := generateLegalMoves()
	positionsCount := 0
	for _, move := range moves {
		makeMove(&move)
		positionsCount += perft(depth - 1)
		unmakeMove()
	}
	return positionsCount
}

func runPerft(depth int) {
	for i := 1; i <= depth; i++ {
		start := time.Now()
		moves := perft(i)
		t := time.Now()
		elapsed := t.Sub(start)
		fmt.Println("depth: ", i, "moves: ", moves, "time: ", elapsed)
	}
}

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

func runPerft() {
	//http://www.rocechess.ch/perft.html
	//starting position: rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1
	//good testposition: r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1
	position = newPosition("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
	for i := 1; i < 5; i++ {
		start := time.Now()
		moves := perft(i)
		t := time.Now()
		elapsed := t.Sub(start)
		fmt.Println("depth: ", i, "moves: ", moves, "time: ", elapsed)
	}
}

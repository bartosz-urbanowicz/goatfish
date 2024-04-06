package main

// http://www.rocechess.ch/perft.html
// starting position: rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1
// good startposition: r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1
// en passant: rnbqkbnr/ppp1p1pp/8/3pPp2/8/8/PPPP1PPP/RNBQKBNR w KQkq f6 0 1
// promotion: n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1
var position = newPosition("n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1")

func main() {
	runGame()
}

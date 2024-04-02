package main

var startingPositionFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
var position = newPosition(startingPositionFEN)
// var position = newPosition("8/1P2k3/8/8/2K5/5p2/8/8 w - - 0 1")

func main() {
	runGame()
}

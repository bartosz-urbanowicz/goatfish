package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseCoordinate(coord string) int {
	fileMap := map[byte]int{
		'a': 0,
		'b': 1,
		'c': 2,
		'd': 3,
		'e': 4,
		'f': 5,
		'g': 6,
		'h': 7,
	}
	fileChar := coord[0]
	fileNumber := fileMap[fileChar]
	rankNumber, err := strconv.Atoi(string(coord[1]))
	if err != nil {
		panic(err)
	}
	return 20 + ((8 - rankNumber) * 10) + fileNumber + 1
}

func printBoard(board [120]byte) {
	pieceMap := map[byte]rune{
		1:  9817,
		2:  9816,
		3:  9815,
		4:  9814,
		5:  9813,
		6:  9812,
		0:  183,
		9:  9823,
		10: 9822,
		11: 9821,
		12: 9820,
		13: 9819,
		14: 9818,
	}
	for i := 2; i < 10; i++ {
		fmt.Print(8 - (i - 2))
		fmt.Print(" ")
		for j := 1; j < 9; j++ {
			fmt.Print(string(pieceMap[board[(i*10)+j]]), " ")
		}
		fmt.Print("\n")
	}
	fmt.Print(" ")
	for i := 1; i < 9; i++ {
		fmt.Printf("%c", (i + 96))
		fmt.Print(" ")
	}
	fmt.Print("\n")
} 

func makeMove(move *Move, legalMoves []Move) {
	var isLegal bool = false
	for _, legalMove := range legalMoves {
		if move.startField == legalMove.startField && move.targetField == legalMove.targetField {
			isLegal = true
		}
	}
	if isLegal {
		piece := position.board[move.startField]
		position.board[move.startField] = 0
		position.board[move.targetField] = piece
		position.sideToMove = !position.sideToMove
		position.fullmoveCounter++
	} else {
		fmt.Println("this move is illegal")
	}
}

func parseMove(move string) *Move {
	move = strings.TrimSuffix(move, "\n")
	move = strings.TrimSuffix(move, "\r")
	startField := parseCoordinate(move[:2])
	targetField := parseCoordinate(move[2:])
	return NewMove(startField, targetField)
}

func getPlayerMove() *Move {
	fmt.Print("enter move: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return parseMove(input)
}

var starting_position_fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq e3 0 1"
var position = newPosition(starting_position_fen)

func main() {
	for {
		legalMoves := generateMoves()
		printBoard(position.board)
		playerMove := getPlayerMove()
		makeMove(playerMove, legalMoves)
	}
}

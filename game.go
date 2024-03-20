package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"slices"
	"strconv"
)

var (
	pieceMap = map[byte]rune{
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
)


func runGame() {
	scanner := bufio.NewScanner(os.Stdin)
	printBoard()
	fmt.Print("> ")
	for scanner.Scan() {
		input := scanner.Text()
		if input == "q" {
			break
		} else {
			legalMoves := generateMoves()
			if strings.HasPrefix(input, "show") {
				field := parseCoordinate(strings.Split(input, " ")[1])
				printLegalMovesFromField(field, legalMoves)
			} else {
				move := parseMove(input)
				makeMove(move, legalMoves)
				printBoard()
			}
			fmt.Print("> ")
		}
	}
}

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

func printBoard() {
	for i := 2; i < 10; i++ {
		//printing rank number
		fmt.Print(8 - (i - 2))
		fmt.Print(" ")
		for j := 1; j < 9; j++ {
			fmt.Print(string(pieceMap[position.board[(i*10)+j]]), " ")
		}
		fmt.Print("\n")
	}
	//printing file letters
	fmt.Print("  ")
	for i := 1; i < 9; i++ {
		fmt.Printf("%c", (i + 96))
		fmt.Print(" ")
	}
	fmt.Print("\n")
}

func printLegalMovesFromField(field int, legalMoves []Move) {
	fmt.Println("printing moves form field", field)
	fieldsToMove := []int{}
	for _, move := range legalMoves {
		if move.startField == field {
			fieldsToMove = append(fieldsToMove, move.targetField)
		}
	}
	for i := 2; i < 10; i++ {
		//printing rank number
		fmt.Print(8 - (i - 2))
		fmt.Print(" ")
		for j := 1; j < 9; j++ {
			currField := (i * 10) + j
			if slices.Contains(fieldsToMove, currField) {
				if position.board[currField] == 0 {
					fmt.Print("#", " ")
				} else {
					fmt.Print("X ")
				}
			} else {
				fmt.Print(string(pieceMap[position.board[currField]]), " ")
			}
		}
		fmt.Print("\n")
	}
	//printing file letters
	fmt.Print("  ")
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
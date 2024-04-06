package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
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
	// TODO deeper perft than 9
	regex := regexp.MustCompile("^(?:[a-h][1-8][a-h][1-8]|q|show [a-h][1-8]|unmake|perft [1-9])$")
	printBoard()
	fmt.Print("> ")
	for scanner.Scan() {
		input := scanner.Text()
		if !regex.MatchString(input) {
			fmt.Print("invalid input\n")
		} else if input == "q" {
			break
		} else if input == "unmake" {
			unmakeMove()
			printBoard()
		} else if strings.HasPrefix(input, "perft") {
			depth, err := strconv.Atoi(strings.Split(input, " ")[1])
			if err != nil {
				panic(err)
			}
			runPerft(depth)
		} else {
			validMoves := generateMoves()
			if strings.HasPrefix(input, "show") {
				field := parseCoordinate(strings.Split(input, " ")[1])
				printLegalMovesFromField(field, validMoves)
			} else {
				move := parseMove(input)
				if isValid(move, validMoves) && isLegal(move) {
					makeMove(move)
				} else {
					fmt.Println("invalid move")
				}
				printBoard()
			}
		}
		fmt.Print("> ")
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
	// fmt.Println("(black to move: ", position.blackToMove, ")")
	// fmt.Println("(castling rights: ", position.castlingRights, ")")
	// fmt.Println("(en passant target: ", position.enPassantTarget, ")")
	// fmt.Println("(halfmove clock: ", position.halfmoveClock, ")")
	// fmt.Println("(fullmove counter: ", position.fullmoveCounter, ")")
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
	fieldsToMove := []int{}
	for _, move := range legalMoves {
		if move.startField == field {
			if isLegal(&move) {
				fieldsToMove = append(fieldsToMove, move.targetField)
			}
		}
	}
	for i := 2; i < 10; i++ {
		//printing rank number
		fmt.Print(8 - (i - 2))
		fmt.Print(" ")
		for j := 1; j < 9; j++ {
			currField := (i * 10) + j
			if slices.Contains(fieldsToMove, currField) {
				if position.board[currField] == 0 && currField != position.enPassantTarget {
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

func parseMove(move string) *Move {
	move = strings.TrimSuffix(move, "\n")
	move = strings.TrimSuffix(move, "\r")
	startField := parseCoordinate(move[:2])
	targetField := parseCoordinate(move[2:])
	if isType(position.board[startField], "pawn") {
		if math.Abs(float64(targetField-startField)) == 20 {
			return NewMove(startField, targetField, firstPawnMove)
		} else if targetField == position.enPassantTarget {
			return NewMove(startField, targetField, enPassant)
		} else if targetField/10 == 2 || targetField/10 == 9 {
			// TODO minor promotions
			return NewMove(startField, targetField, promotionQueen)
		}
	} else {
		castleShortMoves := []string{"e8g8", "e1g1"}
		castleLongMoves := []string{"e8c8", "e1c1"}
		if slices.Contains(castleShortMoves, move) {
			return NewMove(startField, targetField, castleShort)
		} else if slices.Contains(castleLongMoves, move) {
			return NewMove(startField, targetField, castleLong)
		}
	}
	return NewMove(startField, targetField, normal)
}

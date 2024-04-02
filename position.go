package main

import (
	"strings"
	"unicode"
	"strconv"
)

type Position struct {
	board             [120]byte
	blackToMove      bool
	castlingRights   [4]bool
	enPassantTarget int
	halfmoveClock    byte
	fullmoveCounter  uint16
}

func newPosition(fen string) *Position {
	// bits 0-2: piece type
	// bit 3: color
	pieceMap := map[rune]byte{
		'P': 1,
		'N': 2,
		'B': 3,
		'R': 4,
		'Q': 5,
		'K': 6,
		0:   0,
		'p': 9,
		'n': 10,
		'b': 11,
		'r': 12,
		'q': 13,
		'k': 14,
	}

	//board layout
	var board [120]byte
	fields := strings.Split(fen, " ")
	ranks := strings.Split(fields[0], "/")
	currPosition := 0
	for i := 0; i < 20; i++ {
		board[currPosition] = 255
		currPosition++
	}
	for _, rank := range ranks {
		//sentinel fields at left edge
		board[currPosition] = 255
		currPosition++
		for _, char := range rank {
			if unicode.IsDigit(char) {
				// fmt.Println(int(char-'0'), "didgit", curr_position)
				// subtracting the value of rune '0' from any rune '0' through '9' will give you an integer 0 through 9
				for k := 0; k < int(char-'0'); k++ {
					board[currPosition] = 0
					currPosition++
				}
			} else {
				// fmt.Println(string(char), "letter", curr_position)
				board[currPosition] = pieceMap[char]
				currPosition++
			}
		}
		//sentinel fields at right edge
		board[currPosition] = 255
		currPosition++
	}
	for i := 0; i < 20; i++ {
		board[currPosition] = 255
		currPosition++
	}

	//side to move
	blackToMove := fields[1] != "w"

	//castling rights (indexes: white short - 0, white long - 1, black short - 2, black long - 3)
	castlingRights := [4]bool{false, false, false, false}
	if fields[2] != "-" {
		if strings.Contains(fields[2], "K") {
			castlingRights[0] = true
		}
		if strings.Contains(fields[2], "Q") {
			castlingRights[1] = true
		}
		if strings.Contains(fields[2], "k") {
			castlingRights[2] = true
		}
		if strings.Contains(fields[2], "q") {
			castlingRights[3] = true
		}
	}

	//en passant target field
	var enPassantTarget int
	if fields[3] == "-" {
		enPassantTarget = -1
	} else {
		rank, err := strconv.Atoi(string(fields[3][1]))
		if err != nil {
			panic(err)
		}
		//96 becauese we subtract 95 to get a number from ascii value and subtract 1 because of sentinel fields
		enPassantTarget = ((10 - int(rank)) * 10) + (int(fields[3][0]) - 96)
	}

	//halfmove clock
	halfmoveClock, err := strconv.Atoi(fields[4])
	if err != nil {
		panic(err)
	}

	//fullmove counter
	fullmoveCounter, err := strconv.Atoi(fields[5])
	if err != nil {
		panic(err)
	}

	p := Position{board: board,
		blackToMove:      blackToMove,
		castlingRights:   castlingRights,
		enPassantTarget: enPassantTarget,
		halfmoveClock:    byte(halfmoveClock),
		fullmoveCounter:  uint16(fullmoveCounter),
	}
	return &p
}
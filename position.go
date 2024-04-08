package main

import (
	"strconv"
	"strings"
	"unicode"
)

type UnmakeInfo struct {
	move               *Move
	targetFieldContent byte
	castlingRights     [4]bool
	enPassantTarget    int
	halfmoveClock      byte
}

type Position struct {
	board           [120]byte
	blackToMove     bool
	castlingRights  [4]bool
	enPassantTarget int
	halfmoveClock   byte
	fullmoveCounter uint16
	unmakeHistory   []UnmakeInfo
}

var (
	// bits 0-2: piece type
	// bit 3: color
	FENPieceMap = map[rune]byte{
		'P': whitePawn,
		'N': whiteKnight,
		'B': whiteBishop,
		'R': whiteRook,
		'Q': whiteQueen,
		'K': whiteKing,
		'0': empty,
		'p': blackPawn,
		'n': blackKnight,
		'b': blackBishop,
		'r': blackRook,
		'q': blackQueen,
		'k': blackKing,
	}
)

func boardFromFEN(boardString string) [120]byte {
	var board [120]byte
	ranks := strings.Split(boardString, "/")
	currPosition := 0
	for i := 0; i < 20; i++ {
		board[currPosition] = border
		currPosition++
	}
	for _, rank := range ranks {
		//sentinel fields at left edge
		board[currPosition] = border
		currPosition++
		for _, char := range rank {
			if unicode.IsDigit(char) {
				// subtracting the value of rune '0' from any rune '0' through '9' will give you an integer 0 through 9
				for k := 0; k < int(char-'0'); k++ {
					board[currPosition] = 0
					currPosition++
				}
			} else {
				// fmt.Println(string(char), "letter", curr_position)
				board[currPosition] = FENPieceMap[char]
				currPosition++
			}
		}
		//sentinel fields at right edge
		board[currPosition] = border
		currPosition++
	}
	for i := 0; i < 20; i++ {
		board[currPosition] = border
		currPosition++
	}
	return board
}

func castlingRightsFromFEN(castleRightsString string) [4]bool {
	//indexes: white short - 0, white long - 1, black short - 2, black long - 3
	castlingRights := [4]bool{false, false, false, false}
	if castleRightsString != "-" {
		if strings.Contains(castleRightsString, "K") {
			castlingRights[0] = true
		}
		if strings.Contains(castleRightsString, "Q") {
			castlingRights[1] = true
		}
		if strings.Contains(castleRightsString, "k") {
			castlingRights[2] = true
		}
		if strings.Contains(castleRightsString, "q") {
			castlingRights[3] = true
		}
	}
	return castlingRights
}

func enPassantTargetFromFEN(enPassantTargetString string) int {
	var enPassantTarget int
	if enPassantTargetString == "-" {
		enPassantTarget = -1
	} else {
		rank, err := strconv.Atoi(string(enPassantTargetString[1]))
		if err != nil {
			panic(err)
		}
		//96 becauese we subtract 95 to get a number from ascii value and subtract 1 because of sentinel fields
		enPassantTarget = ((10 - int(rank)) * 10) + (int(enPassantTargetString[0]) - 96)
	}
	return enPassantTarget
}

func newPosition(fen string) *Position {
	fields := strings.Split(fen, " ")

	board := boardFromFEN(fields[0])
	blackToMove := fields[1] != "w"
	castlingRights := castlingRightsFromFEN(fields[2])
	enPassantTarget := enPassantTargetFromFEN(fields[3])

	halfmoveClock, err := strconv.Atoi(fields[4])
	if err != nil {
		panic(err)
	}

	fullmoveCounter, err := strconv.Atoi(fields[5])
	if err != nil {
		panic(err)
	}

	p := Position{board: board,
		blackToMove:     blackToMove,
		castlingRights:  castlingRights,
		enPassantTarget: enPassantTarget,
		halfmoveClock:   byte(halfmoveClock),
		fullmoveCounter: uint16(fullmoveCounter),
	}
	return &p
}

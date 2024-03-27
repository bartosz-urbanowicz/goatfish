package main

import (
	"fmt"
	"slices"
)

var (
	blackPawnStartingSquares = []int{31, 32, 33, 34, 35, 36, 37, 38}
	whitePawnStartingSquares = []int{81, 82, 83, 84, 85, 86, 87, 88}
	offsets                  = map[byte]int{
		0: -10,
		1: 1,
		2: 10,
		3: -1,
		4: -11,
		5: -9,
		6: 11,
		7: 9,
	}
)

const (
	normal    byte = 0
	enPassant byte = 1
	firstPawnMove byte = 2
)

type Move struct {
	startField  int
	targetField int
	moveType    byte
}

func NewMove(startField int, targetField int, moveType byte) *Move {
	m := new(Move)
	m.startField = startField
	m.targetField = targetField
	m.moveType = moveType
	return m
}

func generateMoves() []Move {
	fmt.Println("generating moves")
	var moves []Move
	for field, piece := range position.board {
		if piece != 0 {
			if isBlack(piece) == position.sideToMove {
				if isRayPiece(piece) {
					moves = append(moves, generateRayMoves(field, piece)...)
				} else if isType(piece, "pawn") {
					moves = append(moves, generatePawnMoves(field, piece)...)
				} else if isType(piece, "king") {
					moves = append(moves, generateKingMoves(field)...)
				} else if isType(piece, "knight") {
					moves = append(moves, generateKnightMoves(field)...)
				}
			}
		}
	}
	return moves
}

func possibleOffsetsInDirection(currOffset int, startField int, direction byte, possibleOffsets []int) []int {
	field := position.board[startField+currOffset]
	if field == 255 {
		return possibleOffsets
	} else if field != 0 {
		if isBlack(field) == position.sideToMove {
			return possibleOffsets
		} else {
			return append(possibleOffsets, currOffset)
		}
	} else {
		return possibleOffsetsInDirection(currOffset+offsets[direction], startField, direction, append(possibleOffsets, currOffset))
	}
}

func generateRayMoves(field int, piece byte) []Move {
	// 1 => up, 2 => right, 3 => down, 4 => left, 5 => left-up, 6 => right-up, 7 => right-down, 8 => left-down
	moves := []Move{}
	var directions []byte
	if isType(piece, "bishop") {
		directions = []byte{4, 5, 6, 7}
	} else if isType(piece, "rook") {
		directions = []byte{0, 1, 2, 3}
	} else {
		directions = []byte{0, 1, 2, 3, 4, 5, 6, 7}
	}
	for _, dir := range directions {
		possibleOffsets := possibleOffsetsInDirection(offsets[dir], field, dir, []int{})
		for _, offset := range possibleOffsets {
			moves = append(moves, *NewMove(field, field+offset, normal))
		}
	}
	return moves
}

func generatePawnMoves(field int, piece byte) []Move {
	moves := []Move{}
	var possiblePushOffsets []int
	var possibleTakeOffsets [2]int
	if isBlack(piece) {
		possibleTakeOffsets = [2]int{offsets[6], offsets[7]}
		if slices.Contains(blackPawnStartingSquares, field) {
			possiblePushOffsets = []int{offsets[2], 2 * offsets[2]}
		} else {
			possiblePushOffsets = []int{offsets[2]}
		}
	} else {
		possibleTakeOffsets = [2]int{offsets[4], offsets[5]}
		if slices.Contains(whitePawnStartingSquares, field) {
			possiblePushOffsets = []int{offsets[0], 2 * offsets[0]}
		} else {
			possiblePushOffsets = []int{offsets[0]}
		}
	}
	for i, offset := range possiblePushOffsets {
		piece := position.board[field+offset]
		if piece == 0 {
			//if this is the second then the move generates an en passant target square
			if i == 1 {
				moves = append(moves, *NewMove(field, field+offset, firstPawnMove))
			} else {
				moves = append(moves, *NewMove(field, field+offset, normal))
			}
		} else {
			//if the first field we check is blocked with a piece we cant move to the second field
			break
		}
	}
	for _, offset := range possibleTakeOffsets {
		piece := position.board[field+offset]
		if (isPiece(piece) && isBlack(piece) != position.sideToMove) {
			fmt.Println("adding normal take move", field + offset)
			moves = append(moves, *NewMove(field, field+offset, normal))
		} else if field + offset == position.enPassantTarget {
			fmt.Println("adding en passant move", field + offset)
			moves = append(moves, *NewMove(field, field+offset, enPassant))
		}
	}
	return moves
}

func generateKingMoves(field int) []Move {
	moves := []Move{}
	for i := 0; i < 8; i++ {
		piece := position.board[field+offsets[byte(i)]]
		if piece == 0 || isBlack(piece) != position.sideToMove {
			moves = append(moves, *NewMove(field, field+offsets[byte(i)], normal))
		}
	}
	return moves
}

func generateKnightMoves(field int) []Move {
	moves := []Move{}
	knightOffsets := []int{-21, -19, -8, 12, 21, 19, 8, -12}
	for _, offset := range knightOffsets {
		piece := position.board[field+offset]
		if piece == 0 || isBlack(piece) != position.sideToMove {
			moves = append(moves, *NewMove(field, field+offset, normal))
		}
	}
	return moves
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
		position.enPassantTarget = -1
		if move.moveType == firstPawnMove {
			if position.sideToMove {
				position.enPassantTarget = move.targetField - 10
			} else {
				position.enPassantTarget = move.targetField + 10
			}
			
		}
		if move.moveType == enPassant {
			fmt.Println("google en passant")
			if position.sideToMove {
				position.board[move.targetField - 10] = 0
			} else {
				position.board[move.targetField + 10] = 0
			}
		}
		position.sideToMove = !position.sideToMove
		position.fullmoveCounter++
	} else {
		fmt.Println("this move is illegal")
	}
}

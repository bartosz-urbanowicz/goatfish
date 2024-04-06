package main

import (
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
	blackCastleShortFields = []int{26, 27}
	blackCastleLongFields  = []int{22, 23, 24}
	whiteCastleShortFields = []int{96, 97}
	whiteCastleLongFields  = []int{92, 93, 94}
)

const (
	normal          byte = 0
	enPassant       byte = 1
	firstPawnMove   byte = 2
	castleShort     byte = 3
	castleLong      byte = 4
	promotionQueen  byte = 5
	promotionRook   byte = 6
	promotionKnight byte = 7
	promotionBishop byte = 8
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
	var moves []Move
	for field, piece := range position.board {
		if piece != 0 {
			if isBlack(piece) == position.blackToMove {
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
	moves = append(moves, generateCastleMoves()...)
	return moves
}

func generateLegalMoves() []Move {
	moves := generateMoves()
	legalMoves := []Move{}
	for _, move := range moves {
		if isLegal(&move) {
			legalMoves = append(legalMoves, move)
		}
	}
	return legalMoves
}

func checkFieldsEmpty(fields []int) bool {
	for _, field := range fields {
		if position.board[field] != 0 {
			return false
		}
	}
	return true
}

func generateCastleMoves() []Move {
	moves := []Move{}
	if position.blackToMove {
		if checkFieldsEmpty(blackCastleShortFields) && position.castlingRights[2] {
			moves = append(moves, *NewMove(25, 27, castleShort))
		}
		if checkFieldsEmpty(blackCastleLongFields) && position.castlingRights[3] {
			moves = append(moves, *NewMove(25, 23, castleLong))
		}
	} else {
		if checkFieldsEmpty(whiteCastleShortFields) && position.castlingRights[0] {
			moves = append(moves, *NewMove(95, 97, castleShort))
		}
		if checkFieldsEmpty(whiteCastleLongFields) && position.castlingRights[1] {
			moves = append(moves, *NewMove(95, 93, castleLong))
		}
	}
	return moves
}

func possibleOffsetsInDirection(currOffset int, startField int, direction byte, possibleOffsets []int) []int {
	field := position.board[startField+currOffset]
	if field == 255 {
		return possibleOffsets
	} else if field != 0 {
		if isBlack(field) == position.blackToMove {
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

func generatePromotionMoves(field int, offset int) []Move {
	moves := []Move{}
	promotionTypes := [4]byte{promotionQueen, promotionRook, promotionKnight, promotionBishop}
	for _, promotionType := range promotionTypes {
		moves = append(moves, *NewMove(field, field+offset, promotionType))
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
			//if this is the second iteration then the move generates an en passant target square
			if i == 1 {
				moves = append(moves, *NewMove(field, field+offset, firstPawnMove))
			} else {
				if (field+offset)/10 == 2 || (field+offset)/10 == 9 {
					//promotion moves
					moves = append(moves, generatePromotionMoves(field, offset)...)
				} else {
					moves = append(moves, *NewMove(field, field+offset, normal))
				}
			}
		} else {
			//if the first field we check is blocked with a piece we cant move to the second field
			break
		}
	}
	for _, offset := range possibleTakeOffsets {
		piece := position.board[field+offset]
		if isPiece(piece) && isBlack(piece) != position.blackToMove {
			if (field+offset)/10 == 2 || (field+offset)/10 == 9 {
				//promotion moves
				moves = append(moves, generatePromotionMoves(field, offset)...)
			} else {
				moves = append(moves, *NewMove(field, field+offset, normal))
			}
		} else if field+offset == position.enPassantTarget {
			moves = append(moves, *NewMove(field, field+offset, enPassant))
		}
	}
	return moves
}

func generateKingMoves(field int) []Move {
	moves := []Move{}
	for i := 0; i < 8; i++ {
		piece := position.board[field+offsets[byte(i)]]
		if (piece == 0 || isBlack(piece) != position.blackToMove) && piece != 255 {
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
		if (piece == 0 || isBlack(piece) != position.blackToMove) && piece != 255 {
			moves = append(moves, *NewMove(field, field+offset, normal))
		}
	}
	return moves
}

func handleCastlingRights(piece byte, startField int, targetField int) {
	if isType(position.board[targetField], "rook") {
		switch targetField {
		case 21:
			position.castlingRights[3] = false
		case 28:
			position.castlingRights[2] = false
		case 91:
			position.castlingRights[1] = false
		case 98:
			position.castlingRights[0] = false
		}
	} else if isType(piece, "king") {
		if position.blackToMove {
			position.castlingRights[2] = false
			position.castlingRights[3] = false
		} else {
			position.castlingRights[0] = false
			position.castlingRights[1] = false
		}
	} else if isType(piece, "rook") {
		switch startField {
		case 21:
			position.castlingRights[3] = false
		case 28:
			position.castlingRights[2] = false
		case 91:
			position.castlingRights[1] = false
		case 98:
			position.castlingRights[0] = false
		}
	}
}

func checkKingSafe() bool {
	position.blackToMove = !position.blackToMove
	enemyMoves := generateMoves()
	for _, enemyMove := range enemyMoves {
		attackedPiece := position.board[enemyMove.targetField]
		if isType(attackedPiece, "king") {
			position.blackToMove = !position.blackToMove
			return false
		}
	}
	position.blackToMove = !position.blackToMove
	return true
}

func checkCastleLegal(move *Move) bool {
	var fields []int
	if move.moveType == castleShort {
		if position.blackToMove {
			fields = blackCastleShortFields
		} else {
			fields = whiteCastleShortFields
		}
	} else if move.moveType == castleLong {
		//we do not use castleFields because king doesnt pass the b file
		if position.blackToMove {
			fields = []int{23, 24}
		} else {
			fields = []int{93, 94}
		}
	}

	for _, field := range fields {
		if position.blackToMove {
			position.board[field] = 14
		} else {
			position.board[field] = 6
		}
	}
	kingSafe := checkKingSafe()
	for _, field := range fields {
		position.board[field] = 0
	}
	if kingSafe {
		return true
	} else {
		return false
	}
}

func isLegal(move *Move) bool {
	if move.moveType == castleShort || move.moveType == castleLong {
		if !checkCastleLegal(move) {
			return false
		}
	}
	makeMove(move)
	position.blackToMove = !position.blackToMove
	kingSafe := checkKingSafe()
	position.blackToMove = !position.blackToMove
	unmakeMove()
	if kingSafe {
		return true
	} else {
		return false
	}
}

func isValid(move *Move, validMoves []Move) bool {
	for _, validMove := range validMoves {
		if move.startField == validMove.startField && move.targetField == validMove.targetField {
			return true
		}
	}
	return false
}

func makeMove(move *Move) {
	unmakeInfo := UnmakeInfo{
		move:               move,
		targetFieldContent: position.board[move.targetField],
		castlingRights:     position.castlingRights,
		enPassantTarget:    position.enPassantTarget,
		halfmoveClock:      position.halfmoveClock,
	}
	position.unmakeHistory = append(position.unmakeHistory, unmakeInfo)
	piece := position.board[move.startField]
	if isType(piece, "pawn") || position.board[move.targetField] != 0 {
		position.halfmoveClock = 0
	} else {
		position.halfmoveClock++
	}
	handleCastlingRights(piece, move.startField, move.targetField)
	position.board[move.startField] = 0
	position.board[move.targetField] = piece
	position.enPassantTarget = -1
	switch move.moveType {
	case firstPawnMove:
		if position.blackToMove {
			position.enPassantTarget = move.targetField - 10
		} else {
			position.enPassantTarget = move.targetField + 10
		}
	case enPassant:
		if position.blackToMove {
			position.board[move.targetField-10] = 0
		} else {
			position.board[move.targetField+10] = 0
		}
	case castleShort:
		if position.blackToMove {
			position.board[28] = 0
			position.board[26] = 12
		} else {
			position.board[98] = 0
			position.board[96] = 4
		}
	case castleLong:
		if position.blackToMove {
			position.board[21] = 0
			position.board[24] = 12
		} else {
			position.board[91] = 0
			position.board[94] = 4
		}
	case promotionQueen:
		if position.blackToMove {
			position.board[move.targetField] = 13
		} else {
			position.board[move.targetField] = 5
		}
	case promotionRook:
		if position.blackToMove {
			position.board[move.targetField] = 12
		} else {
			position.board[move.targetField] = 4
		}
	case promotionKnight:
		if position.blackToMove {
			position.board[move.targetField] = 10
		} else {
			position.board[move.targetField] = 2
		}
	case promotionBishop:
		if position.blackToMove {
			position.board[move.targetField] = 11
		} else {
			position.board[move.targetField] = 3
		}
	}
	if position.blackToMove {
		position.fullmoveCounter++
	}
	position.blackToMove = !position.blackToMove
}

func unmakeMove() {
	unmakeInfo := position.unmakeHistory[len(position.unmakeHistory)-1]
	position.unmakeHistory = position.unmakeHistory[0 : len(position.unmakeHistory)-1]
	move := unmakeInfo.move
	position.board[move.startField] = position.board[move.targetField]
	position.board[move.targetField] = unmakeInfo.targetFieldContent
	//we already switch the side here to make the conditions below more logical
	position.blackToMove = !position.blackToMove
	switch move.moveType {
	case enPassant:
		if position.blackToMove {
			position.board[move.targetField-10] = 1
		} else {
			position.board[move.targetField+10] = 9
		}
	case castleShort:
		if position.blackToMove {
			position.board[26] = 0
			position.board[28] = 12
		} else {
			position.board[96] = 0
			position.board[98] = 4
		}
	case castleLong:
		if position.blackToMove {
			position.board[24] = 0
			position.board[21] = 12
		} else {
			position.board[94] = 0
			position.board[91] = 4
		}
	case promotionQueen, promotionRook, promotionKnight, promotionBishop:
		if position.blackToMove {
			position.board[move.startField] = 9
		} else {
			position.board[move.startField] = 1
		}
	}
	if position.blackToMove {
		position.fullmoveCounter--
	}
	position.castlingRights = unmakeInfo.castlingRights
	position.enPassantTarget = unmakeInfo.enPassantTarget
	position.halfmoveClock = unmakeInfo.halfmoveClock
}

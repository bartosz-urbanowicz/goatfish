package main

import (
	"slices"
)

var (
	blackPawnStartingSquares     = []int{31, 32, 33, 34, 35, 36, 37, 38}
	whitePawnStartingSquares     = []int{81, 82, 83, 84, 85, 86, 87, 88}
	blackCastleShortFields       = []int{26, 27}
	blackCastleLongFields        = []int{22, 23, 24}
	whiteCastleShortFields       = []int{96, 97}
	whiteCastleLongFields        = []int{92, 93, 94}
	up                       int = -10
	right                    int = 1
	down                     int = 10
	left                     int = -1
	leftUp                   int = -11
	rightUp                  int = -9
	rightDown                int = 11
	leftDown                 int = 9
)

func generateMoves() []Move {
	var moves []Move
	var pieces map[int]bool
	if position.blackToMove {
		pieces = position.blackPieces
	} else {
		pieces = position.whitePieces
	}
	for field, _ := range pieces {
		piece := position.board[field]
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
		if position.board[field] != empty {
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

func possibleOffsetsInDirection(currOffset int, startField int, direction int, possibleOffsets []int) []int {
	field := position.board[startField+currOffset]
	if field == border {
		return possibleOffsets
	} else if field != empty {
		if isBlack(field) == position.blackToMove {
			return possibleOffsets
		} else {
			return append(possibleOffsets, currOffset)
		}
	} else {
		return possibleOffsetsInDirection(currOffset+direction, startField, direction, append(possibleOffsets, currOffset))
	}
}

func generateRayMoves(field int, piece byte) []Move {

	moves := []Move{}
	var directions []int
	if isType(piece, "bishop") {
		directions = []int{leftUp, rightUp, rightDown, leftDown}
	} else if isType(piece, "rook") {
		directions = []int{up, right, down, left}
	} else {
		directions = []int{up, right, down, left, leftUp, rightUp, rightDown, leftDown}
	}
	for _, dir := range directions {
		possibleOffsets := possibleOffsetsInDirection(dir, field, dir, []int{})
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
		possibleTakeOffsets = [2]int{rightDown, leftDown}
		if slices.Contains(blackPawnStartingSquares, field) {
			possiblePushOffsets = []int{down, 2 * down}
		} else {
			possiblePushOffsets = []int{down}
		}
	} else {
		possibleTakeOffsets = [2]int{leftUp, rightUp}
		if slices.Contains(whitePawnStartingSquares, field) {
			possiblePushOffsets = []int{up, 2 * up}
		} else {
			possiblePushOffsets = []int{up}
		}
	}
	for i, offset := range possiblePushOffsets {
		piece := position.board[field+offset]
		if piece == empty {
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
	offsets := []int{up, right, down, left, leftUp, rightUp, rightDown, leftDown}
	for _, offset := range offsets {
		piece := position.board[field+offset]
		if (piece == empty || isBlack(piece) != position.blackToMove) && piece != border {
			moves = append(moves, *NewMove(field, field+offset, normal))
		}
	}
	return moves
}

func generateKnightMoves(field int) []Move {
	moves := []Move{}
	knightOffsets := []int{-21, -19, -8, 12, 21, 19, 8, -12}
	for _, offset := range knightOffsets {
		piece := position.board[field+offset]
		if (piece == empty || isBlack(piece) != position.blackToMove) && piece != border {
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
			position.board[field] = blackKing
		} else {
			position.board[field] = whiteKing
		}
	}
	kingSafe := checkKingSafe()
	for _, field := range fields {
		position.board[field] = empty
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

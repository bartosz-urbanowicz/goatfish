package main

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
	position.board[move.startField] = empty
	if position.board[move.targetField] != empty {
		if position.blackToMove {
			delete(position.whitePieces, move.targetField)
		} else {
			delete(position.blackPieces, move.targetField)
		}
	}
	//if the field was'nt in the set (that happens with en passant) nothing happens
	if position.blackToMove {
		delete(position.blackPieces, move.startField)
		position.blackPieces[move.targetField] = true
	} else {
		delete(position.whitePieces, move.startField)
		position.whitePieces[move.targetField] = true
	}
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
			position.board[move.targetField-10] = empty
			delete(position.whitePieces, move.targetField-10)
		} else {
			position.board[move.targetField+10] = empty
			delete(position.blackPieces, move.targetField+10)
		}
	case castleShort:
		if position.blackToMove {
			position.board[28] = empty

			position.board[26] = blackRook
			position.blackPieces[26] = true
		} else {
			position.board[98] = empty
			delete(position.whitePieces, 98)
			position.board[96] = whiteRook
			position.whitePieces[96] = true
		}
	case castleLong:
		if position.blackToMove {
			position.board[21] = empty
			delete(position.blackPieces, 21)
			position.board[24] = blackRook
			position.blackPieces[24] = true
		} else {
			position.board[91] = empty
			delete(position.whitePieces, 91)
			position.board[94] = whiteRook
			position.whitePieces[94] = true
		}
	case promotionQueen:
		if position.blackToMove {
			position.board[move.targetField] = blackQueen
		} else {
			position.board[move.targetField] = whiteQueen
		}
	case promotionRook:
		if position.blackToMove {
			position.board[move.targetField] = blackRook
		} else {
			position.board[move.targetField] = whiteRook
		}
	case promotionKnight:
		if position.blackToMove {
			position.board[move.targetField] = blackKnight
		} else {
			position.board[move.targetField] = whiteKnight
		}
	case promotionBishop:
		if position.blackToMove {
			position.board[move.targetField] = blackBishop
		} else {
			position.board[move.targetField] = whiteBishop
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
	if position.blackToMove {
		delete(position.blackPieces, move.targetField)
		position.blackPieces[move.startField] = true
		if unmakeInfo.targetFieldContent != 0 {
			position.whitePieces[move.targetField] = true
		}
	} else {
		delete(position.whitePieces, move.targetField)
		position.whitePieces[move.startField] = true
		if unmakeInfo.targetFieldContent != empty {
			position.blackPieces[move.targetField] = true
		}
	}
	switch move.moveType {
	case enPassant:
		if position.blackToMove {
			position.board[move.targetField-10] = whitePawn
			position.whitePieces[move.targetField-10] = true
		} else {
			position.board[move.targetField+10] = blackPawn
			position.blackPieces[move.targetField+10] = true
		}
	case castleShort:
		if position.blackToMove {
			position.board[26] = empty
			delete(position.blackPieces, 26)
			position.board[28] = blackRook
			position.blackPieces[28] = true
		} else {
			position.board[96] = empty
			delete(position.whitePieces, 96)
			position.board[98] = whiteRook
			position.whitePieces[98] = true
		}
	case castleLong:
		if position.blackToMove {
			position.board[24] = empty
			delete(position.blackPieces, 24)
			position.board[21] = blackRook
			position.blackPieces[21] = true
		} else {
			position.board[94] = empty
			delete(position.whitePieces, 94)
			position.board[91] = whiteRook
			position.whitePieces[91] = true
		}
	case promotionQueen, promotionRook, promotionKnight, promotionBishop:
		if position.blackToMove {
			position.board[move.startField] = blackPawn
		} else {
			position.board[move.startField] = whitePawn
		}
	}
	if position.blackToMove {
		position.fullmoveCounter--
	}
	position.castlingRights = unmakeInfo.castlingRights
	position.enPassantTarget = unmakeInfo.enPassantTarget
	position.halfmoveClock = unmakeInfo.halfmoveClock
}

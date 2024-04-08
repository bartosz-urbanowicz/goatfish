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
		} else {
			position.board[move.targetField+10] = empty
		}
	case castleShort:
		if position.blackToMove {
			position.board[28] = empty
			position.board[26] = blackRook
		} else {
			position.board[98] = empty
			position.board[96] = whiteRook
		}
	case castleLong:
		if position.blackToMove {
			position.board[21] = empty
			position.board[24] = blackRook
		} else {
			position.board[91] = empty
			position.board[94] = whiteRook
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
	switch move.moveType {
	case enPassant:
		if position.blackToMove {
			position.board[move.targetField-10] = 1
		} else {
			position.board[move.targetField+10] = 9
		}
	case castleShort:
		if position.blackToMove {
			position.board[26] = empty
			position.board[28] = blackRook
		} else {
			position.board[96] = empty
			position.board[98] = whiteRook
		}
	case castleLong:
		if position.blackToMove {
			position.board[24] = empty
			position.board[21] = blackRook
		} else {
			position.board[94] = empty
			position.board[91] = whiteRook
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

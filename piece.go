package main

var (
	empty       byte = 0
	border      byte = 255
	whitePawn   byte = 1
	whiteKnight byte = 2
	whiteBishop byte = 3
	whiteRook   byte = 4
	whiteQueen  byte = 5
	whiteKing   byte = 6
	blackPawn   byte = 9
	blackKnight byte = 10
	blackBishop byte = 11
	blackRook   byte = 12
	blackQueen  byte = 13
	blackKing   byte = 14
)

func isPiece(piece byte) bool {
	return piece != 0 && piece != 255
}

func isBlack(piece byte) bool {
	return (piece & 0b00001000) == 8
}

func isRayPiece(piece byte) bool {
	rayTypes := [3]string{"bishop", "rook", "queen"}
	for _, rayType := range rayTypes {
		if isType(piece, rayType) {
			return true
		}
	}
	return false
}

func isType(piece byte, pieceType string) bool {
	pieceMap := map[string]byte{
		"pawn":   1,
		"knight": 2,
		"bishop": 3,
		"rook":   4,
		"queen":  5,
		"king":   6,
		"empty":  0,
	}
	// we check only last 3 bits because color doesnt matter
	piece = piece & 7
	return piece == pieceMap[pieceType]
}

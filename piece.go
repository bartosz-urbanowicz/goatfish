package main

func isPiece(piece byte) bool {
	return piece != 0 && piece != 255
}

func isBlack(piece byte) bool {
	return (piece & 0b00001000) == 8
}

func isRayPiece(piece byte) bool {
	rayTypes := [3]string{"bishop", "rook",  "queen"}
	for _, rayType := range rayTypes {
		if isType(piece, rayType) {
			return true
		}
	}
	return false
}

func isType(piece byte, pieceType string) bool {
	pieceMap := map[string]byte{
		"pawn": 1,
		"knight": 2,
		"bishop": 3,
		"rook": 4,
		"queen": 5,
		"king": 6,
		"empty": 0,
	}
	// we check only last 3 bits because color doesnt matter
	piece = piece & 7
	return piece == pieceMap[pieceType]
}
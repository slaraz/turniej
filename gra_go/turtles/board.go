package turtles

var gameBoards = map[string][]Field{}

type Field struct {
	Pawns []Color
}

func getCleanBoard() []Field {
	return make([]Field, 9)
}

func CreateGameBoard(uuid string) []Field {
	gameBoards[uuid] = getCleanBoard()
	return gameBoards[uuid]
}

func DeleteBoard(uuid string) {
	delete(gameBoards, uuid)
}
func MovePawn(board []Field, pawn Color, move int) ([]Field, error) {

	fieldNumber, pawnNumber := findPawn(pawn, board)

	newIndex := fieldNumber + move
	if newIndex < 0 || newIndex >= len(board) {
		return nil, ErrInvalidMove
	}
	if fieldNumber == -1 || pawnNumber == -1 {
		board[newIndex].Pawns = append(board[newIndex].Pawns, pawn)
		return board, nil

	}
	pawnsToSwap := board[fieldNumber].Pawns[pawnNumber:]
	board[fieldNumber].Pawns = board[fieldNumber].Pawns[:pawnNumber]
	board[newIndex].Pawns = append(board[newIndex].Pawns, pawnsToSwap...)
	return board, nil
}
func findPawn(pawn Color, board []Field) (fieldNumber int, pawnNumber int) {
	for fieldNumber, field := range board {
		for pawnNumber, p := range field.Pawns {
			if p == pawn {
				return fieldNumber, pawnNumber
			}
		}
	}
	return -1, -1
}

func CheckIfGameOver(board []Field) (bool, Color) {
	if len(board[len(board)-1].Pawns) == 0 {
		return false, ""
	}
	return true, board[len(board)-1].Pawns[0]
}

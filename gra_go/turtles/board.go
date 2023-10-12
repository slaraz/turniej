package turtles

var gameBoards = map[string][]Field{}

type Field struct {
	Pawns []Color `json:"pawns"`
}

func getCleanBoard() []Field {
	return make([]Field, NUMBER_OF_FIELDS_ON_THE_BOARD)
}

func CreateGameBoard(uuid string) []Field {
	gameBoards[uuid] = getCleanBoard()
	return gameBoards[uuid]
}

//	func DeleteBoard(uuid string) {
//		delete(gameBoards, uuid)
//	}
func MovePawn(board []Field, pawn Color, move int) ([]Field, error) {

	fieldNumber, pawnNumber := findPawn(pawn, board)

	newIndex := fieldNumber + move
	if newIndex >= len(board) {
		newIndex = len(board) - 1
	}
	if fieldNumber < 0 && newIndex < 0 {
		return board, nil
	}
	if fieldNumber == -1 || pawnNumber == -1 {
		board[newIndex].Pawns = append(board[newIndex].Pawns, pawn)
		return board, nil

	}
	pawnsToSwap := board[fieldNumber].Pawns[pawnNumber:]
	board[fieldNumber].Pawns = board[fieldNumber].Pawns[:pawnNumber]
	if newIndex >= 0 {
		board[newIndex].Pawns = append(board[newIndex].Pawns, pawnsToSwap...)
	}
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
	numberOfPawns := len(board[len(board)-1].Pawns)
	return true, board[len(board)-1].Pawns[numberOfPawns-1]
}

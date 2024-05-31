package gotack

// Board 是棋盘的接口，需要实现该棋盘的全部方法，才可以调用博弈树。
type Board interface {
	GetAllMoves(isMaxPlayer bool) []Move
	Move(move Move)
	UndoMove(move Move)
	IsGameOver() bool
	Evaluate(isMaxPlayer bool) int
	Clone() Board
}

// Move 是一个表示棋盘上一步棋动作的接口。
type Move interface {
}

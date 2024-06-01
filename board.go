package gotack

// Board 是棋盘的接口，需要实现该棋盘的全部方法，才可以调用博弈树。
type Board interface {
	// Print 打印棋盘的当前状态。
	Print()

	// GetAllMoves 获取当前状态下所有可能的走法（步法生成）。
	// 参数:
	//   - isMaxPlayer bool: 标识当前是最大化玩家还是最小化玩家。
	// 返回值:
	//   - []Move: 当前状态下所有可能的走法。
	GetAllMoves(isMaxPlayer bool) []Move

	// Move 应用一个走法并返回新的棋盘状态。
	// 参数:
	//   - move Move: 要应用的走法。
	Move(move Move)

	// UndoMove 撤销一个走法并返回新的棋盘状态。
	// 参数:
	//   - move Move: 要撤销的走法。
	UndoMove(move Move)

	// IsGameOver 检查游戏是否结束。
	// 返回值:
	//   - bool: 若游戏结束返回true，否则返回false。
	IsGameOver() bool

	// Hash 生成棋盘状态的哈希值。
	// 返回值:
	//   - uint64: 棋盘状态的哈希值。
	Hash() uint64
	Clone() Board
}

// Move 是一个表示棋盘上一步棋动作的接口。
type Move interface {
	// String 返回一个表示棋盘上一步棋动作的字符串。
	String() string
}

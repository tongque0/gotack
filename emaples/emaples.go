package emaples

import (
	"fmt"

	"github.com/tongque0/gotack"
)

func main() {
	// 假设有一个实现了 gotack.Board 和 gotack.Move 接口的棋盘
	var board gotack.Board // 你的棋盘实现
	evaluator := gotack.NewEvaluator(gotack.AlphaBeta, 3, true, func(board gotack.Board, isMaxPlayer bool) int { return 0 })

	// 获取最佳移动
	bestMove := evaluator.GetBestMove(board)
	fmt.Println("Best Move:", bestMove)
}

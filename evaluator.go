package gotack

import "fmt"

type GameTreeType int

const (
	Minimax   GameTreeType = iota // 使用 Minimax 算法
	AlphaBeta                     // 使用 Alpha-Beta 剪枝算法
	// 可以添加更多的算法类型
)

type Evaluator struct {
	TreeType     GameTreeType
	Depth        int
	IsMaxPlayer  bool
	EvaluateFunc func(board Board, isMaxPlayer bool) int
}

// NewEvaluator 创建并初始化一个 Evaluator，接受博弈树类型、搜索深度、最大玩家标识和评估函数。
func NewEvaluator(treeType GameTreeType, depth int, isMaxPlayer bool, evalFunc func(board Board, isMaxPlayer bool) int) *Evaluator {
	if treeType != AlphaBeta {
		fmt.Println("Unsupported tree type")
		return nil
	}
	return &Evaluator{
		TreeType:     treeType,
		Depth:        depth,
		IsMaxPlayer:  isMaxPlayer,
		EvaluateFunc: evalFunc,
	}
}

// GetBestMove 方法，返回最近一次评估中找到的最佳移动
func (e *Evaluator) GetBestMove(board Board) Move {
	var bestMove Move
	switch e.TreeType {
	case Minimax:
		// 假设 Minimax 方法返回一个 EvaluatedMove 结构
		// result := e.Minimax(board, e.Depth, e.IsMaxPlayer)
		// bestMove = result.Move

	case AlphaBeta:
		// 假设 AlphaBeta 方法返回一个 EvaluatedMove 结构
		_, bestMove := e.alphaBeta(board, e.Depth, -1000000, 1000000, e.IsMaxPlayer)
		return bestMove

	// 添加其他算法类型的 case 分支
	// 例如：
	// case SomeOtherAlgorithm:
	//     result := e.SomeOtherAlgorithm(board, e.Depth, e.IsMaxPlayer)
	//     bestMove = result.Move

	default:
		fmt.Println("Unsupported tree type")
	}
	return bestMove
}


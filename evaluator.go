package gotack

import (
	"fmt"
	"math"
)

type GameTreeType int

const (
	AlphaBeta GameTreeType = iota // 使用 Alpha-Beta  算法
	PVS                           // 使用 PVS 剪枝算法
	// 可以添加更多的算法类型
)

type Evaluator struct {
	TreeType     GameTreeType
	EvaluateFunc func(opts *EvalOptions) float64
	EvalOptions  *EvalOptions
	Board        Board
	Depth        int
	BestMoves    []Move
}

// NewEvaluator 创建并初始化一个 Evaluator 对象。
// 此函数接收以下参数：
//   - treeType: 博弈树的类型。当前支持 AlphaBeta 和 PVS 类型。
//   - opts: EvalOptions 结构，包含用于评估的配置选项，如棋盘状态、搜索深度等。
//   - evalFunc: 评估函数，它接受一个 EvalOptions 指针并返回一个表示局面评估值的浮点数。
//     此函数用于根据 EvalOptions 中的配置来评估棋盘状态。
//
// 返回值：
// - 返回一个指向 Evaluator 的指针。如果不支持指定的博弈树类型，或其他配置错误，则返回 nil。
//
// 示例用法：
//
//	evalFunc := func(opts *EvalOptions) float64 {
//	    // 实现具体的评估逻辑，可能使用 opts.Board 和 opts.Depth
//	    return 0.0 // 示意返回一个评估值
//	}
//	opts := NewEvaluatorOptions(WithDepth(5), WithBoard(someBoard), WithIsMaxPlayer(true))
//	evaluator := NewEvaluator(AlphaBeta, opts, evalFunc)
//
// 注意：评估函数和 EvalOptions 应当正确配合，确保所有必要的配置都被设置。
func NewEvaluator(treeType GameTreeType, opts *EvalOptions, evalFunc func(opts *EvalOptions) float64) *Evaluator {
	if treeType != AlphaBeta && treeType != PVS {
		fmt.Println("Unsupported game tree type.")
		return nil
	}
	return &Evaluator{
		TreeType:     treeType,
		Depth:        opts.Depth,
		Board:        opts.Board,
		EvaluateFunc: evalFunc,
		EvalOptions:  opts,
	}
}

// GetBestMove 返回最近一次评估中找到的最佳移动。
// 此方法通过在指定的棋盘状态上运行博弈树搜索算法来确定最佳移动。
//
// 参数:
// - board: Board 接口，代表当前棋盘的状态，需要由调用者提供。
//
// 返回值:
// - []Move: 从当前棋盘状态中评估得到的最佳移动。
//
// 方法逻辑:
// 根据 Evaluator 结构中的 TreeType 字段，选择适当的博弈树搜索算法。
//   - Minimax: 如果树类型为 Minimax，假设 Minimax 方法返回一个 EvaluatedMove 结构，
//     则调用 Minimax 方法并从中提取最佳移动。
//   - AlphaBeta: 如果树类型为 AlphaBeta，调用 alphaBeta 方法，并直接返回从该方法获得的最佳移动。
//     alpha 和 beta 的初始值分别设为 -1000000 和 1000000。
//   - Other Algorithms: 可以添加其他类型的博弈树算法，每种类型需根据其特定逻辑执行并返回最佳移动。
//   - default: 如果树类型不被支持，打印错误消息并返回一个默认的 Move 结构。
//
// 示例用法:
//
//	board := // 初始化或获取当前棋盘状态
//	evaluator := // 创建并初始化 Evaluator 实例
//	bestMove := evaluator.GetBestMove(board) // 使用 bestMove 进行下一步操作
func (e *Evaluator) GetBestMove() []Move {
	var bestMoves []Move
	var value float64
	switch e.TreeType {
	case AlphaBeta:
		value, bestMoves = e.alphaBeta(e.EvalOptions.Depth, -math.MaxFloat64, math.MaxFloat64, e.EvalOptions.IsMaxPlayer, e.EvalOptions)
	case PVS:
		value, bestMoves = e.pvs(e.EvalOptions.Depth, -math.MaxFloat64, math.MaxFloat64, e.EvalOptions.IsMaxPlayer, e.EvalOptions)
	default:
		fmt.Println("Unsupported tree type")
		return []Move{}
	}
	if e.EvalOptions.IsDetail {
		// 使用基本的 ASCII 字符格式化输出详细信息
		fmt.Println("+-----------------+----------------------------------+")
		fmt.Printf("| %-15s | %-32v |\n", "Algorithm", e.TreeType)
		fmt.Println("+-----------------+----------------------------------+")
		fmt.Printf("| %-15s | %-32d |\n", "Depth", e.EvalOptions.Depth)
		fmt.Println("+-----------------+----------------------------------+")
		fmt.Printf("| %-15s | %-32d |\n", "Step", e.EvalOptions.Step)
		fmt.Println("+-----------------+----------------------------------+")
		fmt.Printf("| %-15s | %-32v |\n", "IsMaxPlayer", e.EvalOptions.IsMaxPlayer)
		fmt.Println("+-----------------+----------------------------------+")
		fmt.Printf("| %-15s | %-32v |\n", "IsDetail", e.EvalOptions.IsDetail)
		fmt.Println("+-----------------+----------------------------------+")
		fmt.Printf("| %-15s | %-32f |\n", "Best Eval Value", value)
		fmt.Println("+-----------------+----------------------------------+")
		fmt.Print("| Best Moves     | ")

		for i, move := range bestMoves {
			if i > 2 { // 仅显示前三个
				break
			}
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(move) // 假设 Move 类型有 String() 方法实现
		}
		fmt.Println(" |")
		fmt.Println("+-----------------+----------------------------------+")

		// 打印 Extra 映射中的额外信息
		if len(e.EvalOptions.Extra) > 0 {
			fmt.Println("+-----------------+----------------------------------+")
			fmt.Println("| Extra Info      | Details                          |")
			fmt.Println("+-----------------+----------------------------------+")
			for key, value := range e.EvalOptions.Extra {
				fmt.Printf("| %-15s | %-32v |\n", key, value)
			}
			fmt.Println("+-----------------+----------------------------------+")
		}
	}
	return bestMoves
}

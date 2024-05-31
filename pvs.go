package gotack

import "math"

func (e *Evaluator) pvs(board Board, depth int, alpha, beta float64, isMaximizingPlayer bool, opts ...interface{}) (float64, Move) {
	if depth == 0 || board.IsGameOver() {
		return e.EvaluateFunc(board, isMaximizingPlayer, opts...), nil
	}

	var bestMove Move
	if isMaximizingPlayer {
		maxEval := math.Inf(-1)
		firstMove := true
		for _, move := range board.GetAllMoves(isMaximizingPlayer) {
			board.Move(move)
			var eval float64
			if firstMove {
				eval, _ = e.pvs(board, depth-1, alpha, beta, false, opts...)
				firstMove = false
			} else {
				// 试探性地用一个更小的窗口进行搜索
				eval, _ = e.pvs(board, depth-1, alpha, alpha+1, false, opts...)
				if eval > alpha && eval < beta { // 如果落在窗口之内，进行完整的重新搜索
					eval, _ = e.pvs(board, depth-1, alpha, beta, false, opts...)
				}
			}
			board.UndoMove(move)

			if eval > maxEval {
				maxEval = eval
				bestMove = move
			}
			alpha = math.Max(alpha, eval)
			if beta <= alpha {
				break
			}
		}
		return maxEval, bestMove
	} else {
		minEval := math.Inf(1)
		firstMove := true
		for _, move := range board.GetAllMoves(isMaximizingPlayer) {
			board.Move(move)
			var eval float64
			if firstMove {
				eval, _ = e.pvs(board, depth-1, alpha, beta, true, opts...)
				firstMove = false
			} else {
				// 试探性地用一个更小的窗口进行搜索
				eval, _ = e.pvs(board, depth-1, beta-1, beta, true, opts...)
				if eval < beta && eval > alpha { // 如果落在窗口之内，进行完整的重新搜索
					eval, _ = e.pvs(board, depth-1, alpha, beta, true, opts...)
				}
			}
			board.UndoMove(move)

			if eval < minEval {
				minEval = eval
				bestMove = move
			}
			beta = math.Min(beta, eval)
			if beta <= alpha {
				break
			}
		}
		return minEval, bestMove
	}
}

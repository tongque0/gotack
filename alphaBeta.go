package gotack

import "math"

func (e *Evaluator) alphaBeta(board Board, depth int, alpha, beta float64, isMaximizingPlayer bool, opts ...interface{}) (float64, Move) {
	if depth == 0 || board.IsGameOver() {
		return e.EvaluateFunc(board, isMaximizingPlayer, opts...), nil
	}

	var bestMove Move
	if isMaximizingPlayer {
		maxEval := math.Inf(-1)
		for _, move := range board.GetAllMoves(isMaximizingPlayer) {
			board.Move(move)
			eval, _ := e.alphaBeta(board, depth-1, alpha, beta, false, opts...)
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
		for _, move := range board.GetAllMoves(isMaximizingPlayer) {
			board.Move(move)
			eval, _ := e.alphaBeta(board, depth-1, alpha, beta, true, opts...)
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

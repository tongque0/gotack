package gotack

import "math"

func (e *Evaluator) alphaBeta(depth int, alpha, beta float64, isMaximizingPlayer bool, opts *EvalOptions) (float64, []Move) {
	if depth == 0 || e.Board.IsGameOver() {
		opts.Extra["depth"] = e.Depth - depth
		return e.Board.EvaluateFunc(*opts), nil
	}

	var bestMoves []Move
	var eval float64
	if isMaximizingPlayer {
		maxEval := math.Inf(-1)
		for _, move := range e.Board.GetAllMoves(isMaximizingPlayer) {
			e.Board.Move(move)
			eval, _ = e.alphaBeta(depth-1, alpha, beta, false, opts)
			e.Board.UndoMove(move)

			if eval > maxEval {
				maxEval = eval
				bestMoves = []Move{move}
			} else if eval == maxEval {
				bestMoves = append(bestMoves, move)
			}
			alpha = math.Max(alpha, eval)
			if beta <= alpha {
				break
			}
		}
		if depth == e.Depth {
			e.BestMoves = bestMoves
		}
		return maxEval, bestMoves
	} else {
		minEval := math.Inf(1)
		for _, move := range e.Board.GetAllMoves(isMaximizingPlayer) {
			e.Board.Move(move)
			eval, _ = e.alphaBeta(depth-1, alpha, beta, true, opts)
			e.Board.UndoMove(move)

			if eval < minEval {
				minEval = eval
				bestMoves = []Move{move}
			} else if eval == minEval {
				bestMoves = append(bestMoves, move)
			}
			beta = math.Min(beta, eval)
			if beta <= alpha {
				break
			}
		}
		if depth == e.Depth {
			e.BestMoves = bestMoves // 只在顶层更新 BestMoves
		}
		return minEval, bestMoves
	}
}

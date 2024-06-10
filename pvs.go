package gotack

import "math"

func (e *Evaluator) pvs(depth int, alpha, beta float64, isMaximizingPlayer bool, opts *EvalOptions) (float64, []Move) {
	if depth == 0 || e.Board.IsGameOver() {
		opts.Extra["depth"] = e.Depth - depth
		return e.Board.EvaluateFunc(*opts), nil
	}

	var bestMoves []Move
	var eval float64
	firstMove := true

	moves := e.Board.GetAllMoves(isMaximizingPlayer)
	if isMaximizingPlayer {
		maxEval := math.Inf(-1)
		for _, move := range moves {
			e.Board.Move(move)
			if firstMove {
				eval, _ = e.pvs(depth-1, alpha, beta, false, opts)
				firstMove = false
			} else {
				// Use a null window search initially
				eval, _ = e.pvs(depth-1, alpha, alpha+1, false, opts)
				// If the result is promising but not proven, re-search
				if eval > alpha && eval < beta {
					eval, _ = e.pvs(depth-1, alpha, beta, false, opts)
				}
			}
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
			e.BestMoves = bestMoves // Only update the BestMoves at the root call
		}
		return maxEval, bestMoves
	} else {
		minEval := math.Inf(1)
		for _, move := range moves {
			e.Board.Move(move)
			if firstMove {
				eval, _ = e.pvs(depth-1, alpha, beta, true, opts)
				firstMove = false
			} else {
				eval, _ = e.pvs(depth-1, beta-1, beta, true, opts)
				if eval < beta && eval > alpha {
					eval, _ = e.pvs(depth-1, alpha, beta, true, opts)
				}
			}
			e.Board.UndoMove(move)

			if eval < minEval {
				minEval = eval
				bestMoves = []Move{move}
			} else if eval == minEval {
				bestMoves = append(bestMoves, move)
			}
			beta = math.Min(beta, eval)
			if alpha >= beta {
				break
			}
		}
		if depth == e.Depth {
			e.BestMoves = bestMoves // Only update the BestMoves at the root call
		}
		return minEval, bestMoves
	}
}

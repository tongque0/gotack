package gotack

import (
	"math"
	"sync"
)

func (e *Evaluator) alphaBeta(depth int, alpha, beta float64, isMaximizingPlayer bool, opts *EvalOptions) (float64, []Move) {
	if depth == 0 || e.Board.IsGameOver() {
		opts.Extra["depth"] = e.Depth - depth
		return e.Board.EvaluateFunc(*opts), nil
	}
	threadNum := opts.GetExtraOption("ThreadNum", 1)

	return e.parallelSearch(depth, alpha, beta, isMaximizingPlayer, opts, threadNum)
}

func (e *Evaluator) parallelSearch(depth int, alpha, beta float64, isMaximizingPlayer bool, opts *EvalOptions, threadNum int) (float64, []Move) {
	var bestMoves []Move
	var bestValue float64
	if isMaximizingPlayer {
		bestValue = math.Inf(-1)
	} else {
		bestValue = math.Inf(1)
	}

	moves := e.Board.GetAllMoves(isMaximizingPlayer)
	results := make(chan struct {
		eval float64
		move Move
	}, len(moves))

	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, move := range moves {
		wg.Add(1)
		go func(move Move) {
			defer wg.Done()
			e.Board.Move(move)
			eval, _ := e.alphaBeta(depth-1, alpha, beta, !isMaximizingPlayer, opts)
			e.Board.UndoMove(move)
			results <- struct {
				eval float64
				move Move
			}{eval, move}
		}(move)

		if len(results) == threadNum {
			res := <-results
			e.processResult(res, &bestValue, &bestMoves, &alpha, &beta, isMaximizingPlayer, &mu)
		}
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		e.processResult(res, &bestValue, &bestMoves, &alpha, &beta, isMaximizingPlayer, &mu)
	}

	if depth == e.Depth {
		e.BestMoves = bestMoves
	}

	return bestValue, bestMoves
}

func (e *Evaluator) processResult(res struct {
	eval float64
	move Move
}, bestValue *float64, bestMoves *[]Move, alpha, beta *float64, isMaximizingPlayer bool, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()
	if isMaximizingPlayer && res.eval > *bestValue || !isMaximizingPlayer && res.eval < *bestValue {
		*bestValue = res.eval
		*bestMoves = []Move{res.move}
	} else if res.eval == *bestValue {
		*bestMoves = append(*bestMoves, res.move)
	}
	if isMaximizingPlayer {
		*alpha = math.Max(*alpha, res.eval)
		if *beta <= *alpha {
			return
		}
	} else {
		*beta = math.Min(*beta, res.eval)
		if *beta <= *alpha {
			return
		}
	}
}

package gotack

import (
	"math"
	"math/rand"
	"sort"
	"time"
)

type Node struct {
	State           Board
	Parent          *Node
	Children        []*Node
	Visits          int
	TotalReward     float64
	IsMaxPlayer     bool
	Move            Move
	UntriedMoves    []Move
	SimulationCount int // 模拟次数计数器
	ExpandedCount   int // 已扩展的节点数
}

// UCT uses Monte Carlo Tree Search algorithm to evaluate the current board state and return the best move.
func (e *Evaluator) uct(opts *EvalOptions) (float64, []Move) {
	root := &Node{State: e.Board, IsMaxPlayer: opts.IsMaxPlayer}

	startTime := time.Now()
	timeLimit := time.Duration(opts.TimeLimit) * time.Second
	iterations := opts.Iterations
	if timeLimit == 0 {
		timeLimit = time.Duration(math.MaxInt64)
	}
	if iterations == 0 {
		iterations = math.MaxInt32
	}

	if timeLimit == 0 && iterations == 0 {
		timeLimit = 10 * time.Second
	}

	// Configuration for expansion and simulation
	simulationThreshold := getOptionInt(opts.Extra, "SimThresh", 1)
	expandThreshold := getOptionInt(opts.Extra, "ExpandThresh", 1000)
	expandStep := getOptionInt(opts.Extra, "ExpandStep", 5)
	expandTopN := getOptionInt(opts.Extra, "ExpandTopN", 250)
	aheadStep := getOptionInt(opts.Extra, "AheadStep", 0)

	for i := 0; i < iterations; i++ {
		if time.Since(startTime) >= timeLimit {
			break
		}

		node := e.selectNode(root)
		if !node.State.IsGameOver() {
			result := e.simulate(node, aheadStep)
			e.backpropagate(node, result)
			node.SimulationCount++
			if node.SimulationCount >= simulationThreshold {
				e.expandNode(node, expandThreshold, expandStep, expandTopN)
				node.SimulationCount = 0
			}
		}
	}

	return e.selectBestMove(root)
}

// getOptionInt 从配置映射中提取整数值，如果未找到或类型不匹配，则返回默认值。
// options 是传递给函数的配置映射，key 是要检索的配置项，defaultValue 是找不到时的返回值。
// 返回配置项的整数值或在未找到时返回默认值。
func getOptionInt(options map[string]interface{}, key string, defaultValue int) int {
	if val, ok := options[key]; ok {
		if num, ok := val.(int); ok {
			return num
		}
	}
	return defaultValue
}

// UCTValue 计算并返回节点的UCT值，用于在树搜索中选择节点。
// totalVisits 是到达当前节点路径上的所有访问总次数。
// 返回节点的UCT评估值。
func (n *Node) UCTValue(totalVisits int) float64 {
	if n.Visits == 0 {
		return math.Inf(1)
	}
	avgReward := n.TotalReward / float64(n.Visits)
	exploration := math.Sqrt(2 * math.Log(float64(totalVisits)) / float64(n.Visits))
	return avgReward + exploration
}

// selectNode 根据UCT值递归选择最优子节点，直到达到叶节点。
// node 是当前考察的节点。
// 返回选中的叶节点。
func (e *Evaluator) selectNode(node *Node) *Node {
	for len(node.Children) > 0 {
		bestUCT := -math.MaxFloat64
		var bestChild *Node
		for _, child := range node.Children {
			uctValue := child.UCTValue(node.Visits)
			if uctValue > bestUCT {
				bestUCT = uctValue
				bestChild = child
			}
		}
		node = bestChild
	}
	return node
}

// expandNode 根据访问次数和扩展阈值动态地在树中扩展新的节点。
// node 是当前需要扩展的节点，expandThreshold 是节点访问次数的阈值，
// expandStep 是达到扩展阈值时应该扩展的节点数量，expandTopN 是节点可以扩展的最大子节点数。
func (e *Evaluator) expandNode(node *Node, expandThreshold int, expandStep int, expandTopN int) {
	// 首次初始化未尝试的移动列表
	if len(node.UntriedMoves) == 0 {
		allMoves := node.State.GetAllMoves(node.IsMaxPlayer)
		evaluateAndSortMoves(allMoves, node, e.EvalOptions)
		node.UntriedMoves = allMoves // 存储所有可尝试的移动
	}

	// 计算应该扩展的次数，即节点访问次数除以扩展阈值
	numExpansionsNeeded := node.Visits / expandThreshold
	// 计算当前已经扩展的次数
	currentExpansionsDone := node.ExpandedCount / expandStep

	if numExpansionsNeeded > currentExpansionsDone {
		// 计算新的扩展目标，即上次扩展后额外增加的扩展数量
		additionalExpansions := (numExpansionsNeeded - currentExpansionsDone) * expandStep
		// 更新目标扩展计数
		targetExpandCount := min(node.ExpandedCount+additionalExpansions, expandTopN)
		// 执行扩展操作，直到达到目标扩展计数或未尝试移动用尽
		for node.ExpandedCount < targetExpandCount && node.ExpandedCount < len(node.UntriedMoves) {
			move := node.UntriedMoves[node.ExpandedCount]
			newState := node.State.Clone()
			newState.Move(move)
			childNode := &Node{
				State:       newState,
				Parent:      node,
				IsMaxPlayer: !node.IsMaxPlayer,
				Move:        move,
			}
			node.Children = append(node.Children, childNode)
			node.ExpandedCount++
		}
	}
}
func evaluateAndSortMoves(moves []Move, node *Node, opts *EvalOptions) {
	moveEvaluations := make([]struct {
		move  Move
		value float64
	}, len(moves))

	for i, move := range moves {
		newState := node.State.Clone()
		newState.Move(move)
		moveEvaluations[i] = struct {
			move  Move
			value float64
		}{
			move:  move,
			value: node.State.EvaluateFunc(*opts),
		}
	}

	sort.Slice(moveEvaluations, func(i, j int) bool {
		if node.IsMaxPlayer {
			return moveEvaluations[i].value > moveEvaluations[j].value
		}
		return moveEvaluations[i].value < moveEvaluations[j].value
	})

	for i, eval := range moveEvaluations {
		moves[i] = eval.move
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (e *Evaluator) simulate(node *Node, aheadStep int) float64 {
	currentState := node.State.Clone()
	isMaxPlayer := node.IsMaxPlayer

	for steps := 0; steps < aheadStep && !currentState.IsGameOver(); steps++ {
		moves := currentState.GetAllMoves(isMaxPlayer)
		if len(moves) == 0 {
			break
		}
		moveIndex := rand.Intn(len(moves))
		currentState.Move(moves[moveIndex])
		isMaxPlayer = !isMaxPlayer
	}

	return e.evaluateGameState(currentState)
}

func (e *Evaluator) backpropagate(node *Node, result float64) {
	for node != nil {
		node.Visits++
		node.TotalReward += result
		node = node.Parent
	}
}

func (e *Evaluator) extractMoves(root, bestMove *Node) []Move {
	var moves []Move
	current := bestMove
	for current != nil && current != root {
		moves = append(moves, current.Move)
		current = current.Parent
	}
	return moves
}

func (e *Evaluator) evaluateGameState(state Board) float64 {
	return state.EvaluateFunc(*e.EvalOptions)
}

func (e *Evaluator) selectBestMove(root *Node) (float64, []Move) {
	var bestMove *Node
	maxVisits := -1
	for _, child := range root.Children {
		if child.Visits > maxVisits {
			bestMove = child
			maxVisits = child.Visits
		}
	}
	if bestMove != nil {
		return bestMove.TotalReward / float64(bestMove.Visits), e.extractMoves(root, bestMove)
	}
	return 0.0, nil
}

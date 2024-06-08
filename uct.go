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

// UCT 使用蒙特卡洛树搜索算法（UCT）来评估当前棋盘状态，并返回最佳移动。
// 该函数接收一个 EvalOptions 结构体，包含了评估器的配置选项，如搜索深度、时间限制等。
// 自定义参数 Extra 可以用于传递额外的配置信息，例如模拟阈值、初始扩展等。
// Extra["SimThresh"] 用于设置模拟阈值，当某个节点的模拟次数达到该阈值时，将扩展该节点。
// Extra["IES"] 用于设置初始扩展的节点数。
// Extra["ExpandStep"] 用于设置每次扩展的节点数。
// Extra["TopN"] 用于设置每次扩展的最大节点数。
// 返回值：
// - float64: 表示当前棋盘状态的评估值。
// - []Move: 表示最佳移动序列。
func (e *Evaluator) UCT(opts *EvalOptions) (float64, []Move) {
	root := &Node{State: e.Board, IsMaxPlayer: opts.IsMaxPlayer}

	// 设置时间限制和迭代次数
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

	var simulationThreshold int
	if val, ok := opts.Extra["SimThresh"]; ok {
		simulationThreshold, ok = val.(int)
		if !ok {
			simulationThreshold = 0 // 如果类型转换失败，则不采用该策略
		}
	} else {
		simulationThreshold = 0 // 不采用该策略
	}

	var initialExpand int
	if val, ok := opts.Extra["IES"]; ok {
		initialExpand, ok = val.(int)
		if !ok {
			initialExpand = 5 // 默认初始扩展5个位置
		}
	} else {
		initialExpand = 5 // 默认初始扩展5个位置
	}

	var expandStep int
	if val, ok := opts.Extra["ExpandStep"]; ok {
		expandStep, ok = val.(int)
		if !ok {
			expandStep = 5 // 默认每次扩展5个位置
		}
	} else {
		expandStep = 5 // 默认每次扩展5个位置
	}

	var topN int
	if val, ok := opts.Extra["TopN"]; ok {
		topN, ok = val.(int)
		if !ok {
			topN = 250 // 默认最大扩展250个位置
		}
	} else {
		topN = 250 // 默认最大扩展250个位置
	}

	// 开始 MCTS 迭代
	for i := 0; i < iterations; i++ {
		// 检查是否超出时间限制
		if time.Since(startTime) >= timeLimit {
			break
		}

		node := e.selectNode(root)
		if !node.State.IsGameOver() {
			result := e.simulate(node)
			e.backpropagate(node, result)
			node.SimulationCount++
			if simulationThreshold > 0 && node.SimulationCount >= simulationThreshold {
				e.expandNode(node, initialExpand, expandStep, topN)
				node.SimulationCount = 0
			}
		}
	}

	// 选择最佳移动
	var bestMove *Node
	maxVisits := -1
	for _, child := range root.Children {
		if child.Visits > maxVisits {
			bestMove = child
			maxVisits = child.Visits
		}
	}

	// 提取最佳移动序列
	if bestMove != nil {
		return bestMove.TotalReward / float64(bestMove.Visits), e.extractMoves(root, bestMove)
	}

	// 如果没有找到最佳移动，返回默认值
	return 0.0, nil
}

func (n *Node) UCTValue(totalVisits int) float64 {
	if n.Visits == 0 {
		return math.Inf(1)
	}
	avgReward := n.TotalReward / float64(n.Visits)
	exploration := math.Sqrt(2 * math.Log(float64(totalVisits)) / float64(n.Visits))
	return avgReward + exploration
}

func (e *Evaluator) selectNode(node *Node) *Node {
	for len(node.Children) > 0 {
		bestUCT := -1.0
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

func (e *Evaluator) expandNode(node *Node, initialExpand int, expandStep int, topN int) {
	// 如果还没有初始化未尝试过的动作，则初始化
	if len(node.UntriedMoves) == 0 {
		allMoves := node.State.GetAllMoves(node.IsMaxPlayer)

		// 计算每个未尝试过的动作的估值
		evaluateMove := func(move Move) float64 {
			newState := node.State.Clone()
			(*newState).Move(move)
			return e.evaluateGameState(*newState)
		}

		moveEvaluations := make([]struct {
			move  Move
			value float64
		}, len(allMoves))

		for i, move := range allMoves {
			moveEvaluations[i] = struct {
				move  Move
				value float64
			}{
				move:  move,
				value: evaluateMove(move),
			}
		}

		// 按估值排序，根据玩家类型决定升序还是降序
		if node.IsMaxPlayer {
			// 最大化玩家，按从高到低排序
			sort.Slice(moveEvaluations, func(i, j int) bool {
				return moveEvaluations[i].value > moveEvaluations[j].value
			})
		} else {
			// 最小化玩家，按从低到高排序
			sort.Slice(moveEvaluations, func(i, j int) bool {
				return moveEvaluations[i].value < moveEvaluations[j].value
			})
		}

		// 只保留前 topN 个位置
		if len(moveEvaluations) < topN {
			topN = len(moveEvaluations)
		}

		node.UntriedMoves = make([]Move, topN)
		for i := 0; i < topN; i++ {
			node.UntriedMoves[i] = moveEvaluations[i].move
		}
	}

	// 每次扩展 initialExpand 或 expandStep 个节点
	expandCount := initialExpand
	if node.ExpandedCount > 0 {
		expandCount = expandStep
	}

	end := node.ExpandedCount + expandCount
	if end > len(node.UntriedMoves) {
		end = len(node.UntriedMoves)
	}

	for i := node.ExpandedCount; i < end; i++ {
		move := node.UntriedMoves[i]
		newState := node.State.Clone()
		(*newState).Move(move)
		childNode := &Node{
			State:       *newState,
			Parent:      node,
			IsMaxPlayer: !node.IsMaxPlayer,
			Move:        move,
		}
		node.Children = append(node.Children, childNode)
	}

	node.ExpandedCount = end
}

func (e *Evaluator) simulate(node *Node) float64 {
	currentState := node.State.Clone()
	isMaxPlayer := node.IsMaxPlayer
	for !(*currentState).IsGameOver() {
		moves := (*currentState).GetAllMoves(isMaxPlayer)
		if len(moves) == 0 {
			break
		}
		move := moves[rand.Intn(len(moves))]
		(*currentState).Move(move)
		isMaxPlayer = !isMaxPlayer
	}
	return e.evaluateGameState(*currentState)
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
		moves = append([]Move{current.Move}, moves...)
		current = current.Parent
	}
	return moves
}

func (e *Evaluator) evaluateGameState(state Board) float64 {
	return state.EvaluateFunc()
}

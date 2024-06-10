package gotack

// EvalOptions 定义了评估器的配置选项，用于控制评估过程的各个方面。
type EvalOptions struct {
	// Board 表示当前游戏的棋盘状态，是评估过程中的基本输入。
	Board Board

	// Depth 表示评估时的当前需要搜索深度。这个深度在评估过程中可能会递减，不是搜索树的总深度。
	Depth int

	// Step 用于指定在评估过程的步长，可能影响某些评估逻辑的迭代或递归处理。
	Step int

	// IsDetail 控制是否输出详细的评估过程信息，有助于调试和详细的性能分析。
	IsDetail bool

	// IsMaxPlayer 指示当前评估的玩家是否是最大化玩家，通常在博弈树中用于区分玩家角色。
	IsMaxPlayer bool

	// Iterations 表示评估器在评估过程中的迭代次数，用于控制评估的精度和性能(MCTS)
	// 默认为 0，表示不限制迭代次数,直到时间限制到达。
	Iterations int

	// TimeLimit 表示评估器在评估过程中的时间限制，用于控制评估的性能和实时性。
	// 时间单位为秒，默认为 10 秒, 0 表示不限制时间(注意：迭代次数与时间限制不可同时为0)
	TimeLimit int

	// Thread 表示评估器在评估过程中的线程数，用于并行计算和提高评估的性能。
	ThreadNum int

	// Extra 提供了一个映射，用于存储评估过程中可能需要的任何额外信息或自定义数据。
	// 这使得 EvalOptions 可以灵活地适应各种额外的需求，而无需修改结构体定义。
	Extra map[string]interface{}
}

type EvalOption func(*EvalOptions)

// NewEvaluatorOptions 创建并初始化一个带有默认配置的 EvalOptions 实例。
// 可以通过传入不同的 EvalOption 配置函数来自定义配置项，例如 Depth 或 Board。
func NewEvaluatorOptions(opts ...EvalOption) *EvalOptions {
	opt := &EvalOptions{
		Depth:       1,
		Step:        1,
		Iterations:  0,
		TimeLimit:   10,
		IsDetail:    false,
		IsMaxPlayer: true,
		ThreadNum:   1,
		Extra:       make(map[string]interface{}),
	}
	for _, o := range opts {
		o(opt)
	}
	return opt
}

// GetExtraOption 从 EvalOptions 的 Extra 映射中获取指定键的值。
func (e *EvalOptions) GetExtraOption(key string, def int) int {
	if v, ok := e.Extra[key]; ok {
		if value, ok := v.(int); ok {
			return value
		}
	}
	return def
}

// WithBoard 配置 EvalOptions 的 Board 属性。
func WithBoard(board Board) EvalOption {
	return func(opts *EvalOptions) {
		opts.Board = board
	}
}

// WithDepth 配置 EvalOptions 的 Depth 属性。
// 注意这个深度指的是在评估过程中会递减的深度，不是固定的搜索树的总深度。
func WithDepth(depth int) EvalOption {
	return func(opts *EvalOptions) {
		opts.Depth = depth
	}
}

// WithStep 配置 EvalOptions 的 Step 属性。
func WithStep(step int) EvalOption {
	return func(opts *EvalOptions) {
		opts.Step = step
	}
}

// WithIsDetail 配置 EvalOptions 的 IsDetail 属性，决定是否展示详细信息。
func WithIsDetail(isDetail bool) EvalOption {
	return func(opts *EvalOptions) {
		opts.IsDetail = isDetail
	}
}

// WithIterations 配置 EvalOptions 的 Iterations 属性，用于控制评估的迭代次数。
func WithIterations(iterations int) EvalOption {
	return func(opts *EvalOptions) {
		opts.Iterations = iterations
	}
}

// WithTimeLimit 配置 EvalOptions 的 TimeLimit 属性，用于控制评估的时间限制。
func WithTimeLimit(timeLimit int) EvalOption {
	return func(opts *EvalOptions) {
		opts.TimeLimit = timeLimit
	}
}

// WithThreadNum 配置 EvalOptions 的 ThreadNum 属性，用于控制评估的线程数。
func WithThreadNum(threadNum int) EvalOption {
	return func(opts *EvalOptions) {
		opts.ThreadNum = threadNum
	}
}

// WithIsMaxPlayer 配置 EvalOptions 的 IsMaxPlayer 属性，指示当前评估的玩家是否是最大化玩家。
func WithIsMaxPlayer(isMaxPlayer bool) EvalOption {
	return func(opts *EvalOptions) {
		opts.IsMaxPlayer = isMaxPlayer
	}
}

// WithExtra 允许向 EvalOptions 的 Extra 映射中添加自定义键值对。
// 这可以用于存储评估过程中需要的任何额外信息。
func WithExtra(key string, value interface{}) EvalOption {
	return func(opts *EvalOptions) {
		opts.Extra[key] = value
	}
}

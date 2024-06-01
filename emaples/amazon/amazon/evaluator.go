package amazon

import (
	"github.com/tongque0/gotack"
)

func EvaluateFunc(board gotack.Board, isMaxPlayer bool, opts ...interface{}) float64 {
	// 尝试将board转换为*AmazonBoard类型
	amazonBoard, ok := board.(*AmazonBoard)
	if !ok {
		// fmt.Println("EvaluateFunc called with a board type that is not *AmazonBoard")
		return 0.0 // 或者处理这种情况的其他方式
	}

	// 解析opts，获取轮数
	var turn int
	if len(opts) > 0 {
		turnVal, ok := opts[0].(int) // 类型断言，将opts的第一个元素转换为整数
		if !ok {
			// fmt.Println("EvaluateFunc: Expected an integer for the turn, but got something else.")
			// 处理错误或者使用默认值
		} else {
			turn = turnVal
		}
	}
	// 假设我们使用AmazonBoard的TurnID属性和isMaxPlayer来调用CalculateEvaluationValue
	value := amazonBoard.CalculateEvaluationValue(turn, isMaxPlayer)

	// 返回评估分数
	return value
}

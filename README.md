# gotack

`gotack` 是一个 Go 语言编写的博弈树库，用于实现博弈理论中的Alpha-Beta等剪枝算法。它提供了灵活的接口来适配不同的棋盘游戏。

## 特性

- 支持 Minimax 算法。
- 支持 Alpha-Beta 剪枝算法。
- 易于集成到其他棋盘游戏项目。
- 提供清晰的接口来定义棋盘和移动。

## 快速开始

### 安装

要开始使用 `gotack`，请确保你已安装 Go 环境（版本 1.13 或更高），然后执行以下命令：
```bash
go get github.com/tongque0/gotack
```

### 示例代码

以下是如何在你的项目中使用 `gotack` 的一个简单示例：

```go
package main

import (
    "github.com/tongque0/gotack"
    "fmt"
)

func main() {
    // 假设有一个实现了 gotack.Board 和 gotack.Move 接口的棋盘
    var board gotack.Board // 你的棋盘实现
    evaluator := gotack.NewEvaluator(gotack.AlphaBeta, 3, true, board.Evaluate)

    // 获取最佳移动
    bestMove := evaluator.GetBestMove(board)
    fmt.Println("Best Move:", bestMove)
}
```
### 文档
更详细的 API 文档和更多示例，请查看 API 文档。

### 贡献
我们欢迎任何形式的贡献，无论是新功能的建议、问题报告还是拉取请求。

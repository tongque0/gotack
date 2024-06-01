package amazon

import "math"

// CalculateKingMoves 计算并返回两个棋盘，分别表示黑白棋king走法的棋盘
func (b *AmazonBoard) CalculateKingMoves() (KingmoveBlack, KingmoveWhite AmazonBoard) {
	// 初始化KingmoveBlack和KingmoveWhite为Empty
	for x := range KingmoveBlack {
		for y := range KingmoveBlack[x] {
			KingmoveBlack[x][y] = Empty
			KingmoveWhite[x][y] = Empty
		}
	}
	// 遍历棋盘
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			if b[x][y] == Black || b[x][y] == White {
				// 检查周围8个方向
				for _, d := range dir {
					newX, newY := x+d[0], y+d[1]
					if newX >= 0 && newX < 10 && newY >= 0 && newY < 10 && b[newX][newY] == Empty {
						if b[x][y] == Black {
							KingmoveBlack[newX][newY] = 1 // 标记为可移动
						} else {
							KingmoveWhite[newX][newY] = 1 // 标记为可移动
						}
					}
				}
			}
		}
	}
	return KingmoveBlack, KingmoveWhite
}

func (b *AmazonBoard) CalculateQueenMoves() (QueenmoveBlack, QueenmoveWhite AmazonBoard) {
	// 使用一个较大的数值来初始化棋盘，代表未被访问/不可达
	const maxSteps = 100
	for x := range QueenmoveBlack {
		for y := range QueenmoveBlack[x] {
			QueenmoveBlack[x][y] = maxSteps
			QueenmoveWhite[x][y] = maxSteps
		}
	}

	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			currentPiece := b[x][y]
			if currentPiece == Black || currentPiece == White {
				for _, d := range dir {
					steps := 1 // 从当前位置出发，所以步数从1开始计算
					newX, newY := x+d[0], y+d[1]
					for newX >= 0 && newX < 10 && newY >= 0 && newY < 10 && b[newX][newY] == Empty {
						// 更新到达该位置的最小步数
						if currentPiece == Black {
							if steps < QueenmoveBlack[newX][newY] {
								QueenmoveBlack[newX][newY] = steps
							}
						} else {
							if steps < QueenmoveWhite[newX][newY] {
								QueenmoveWhite[newX][newY] = steps
							}
						}
						newX += d[0]
						newY += d[1]
						steps++
					}
				}
			}
		}
	}
	return QueenmoveBlack, QueenmoveWhite
}

// CalculateTerritoryValue 计算并返回双方基于king走法的领土值
func (b *AmazonBoard) CalculateKingTerritory() (float64, float64) {
	kingMovesBlack, kingMovesWhite := b.CalculateKingMoves()
	var tkBlack, tkWhite float64

	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			if b[x][y] == Empty {
				blackSteps := kingMovesBlack[x][y]
				whiteSteps := kingMovesWhite[x][y]

				switch {
				case blackSteps == whiteSteps && blackSteps != 0:
					tkBlack += 0.5
					tkWhite += 0.5
				case blackSteps != 0 && (whiteSteps == 0 || blackSteps < whiteSteps):
					tkBlack += 1
				case whiteSteps != 0 && (blackSteps == 0 || whiteSteps < blackSteps):
					tkWhite += 1
				}
			}
		}
	}

	return tkBlack, tkWhite
}

// CalculateTerritoryValue 计算并返回双方基于女王走法的领土值
func (b *AmazonBoard) CalculateQueenTerritory() (float64, float64) {
	queenMovesBlack, queenMovesWhite := b.CalculateQueenMoves()
	var tqBlack, tqWhite float64

	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			if b[x][y] == Empty { // 仅考虑空格
				blackSteps := queenMovesBlack[x][y]
				whiteSteps := queenMovesWhite[x][y]

				// 比较双方的步数并计算领土值
				switch {
				case blackSteps == whiteSteps:
					if blackSteps != 0 { // 如果双方步数相同且都能到达
						tqBlack += 0.5
						tqWhite += 0.5
					}
				case blackSteps < whiteSteps:
					if whiteSteps == 0 {
						// 如果白方到达不了
						tqBlack += 2
					} else {
						// 如果黑方步数少
						tqBlack += 1
					}
				case blackSteps > whiteSteps:
					if blackSteps == 0 {
						// 如果黑方到达不了
						tqWhite += 2
					} else {
						// 如果白方步数少
						tqWhite += 1
					}
				}
			}
		}
	}

	return tqBlack, tqWhite
}

func (b *AmazonBoard) CalculateP1P2(queenMovesBlack, queenMovesWhite, kingMovesBlack, kingMovesWhite AmazonBoard) (float64, float64) {
	var p1, p2 float64

	// 计算P1，基于queen走法
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			if b[x][y] == Empty { // 只考虑空格
				blackSteps := float64(queenMovesBlack[x][y])
				whiteSteps := float64(queenMovesWhite[x][y])

				if blackSteps != 100 && whiteSteps != 100 {
					p1 += math.Pow(2.0, -blackSteps) - math.Pow(2.0, -whiteSteps)
				} else if blackSteps != 100 {
					p1 += math.Pow(2.0, -blackSteps)
				} else if whiteSteps != 100 {
					p1 -= math.Pow(2.0, -whiteSteps)
				}
			}
		}
	}
	p1 *= 2

	// 计算P2，基于king走法
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			if b[x][y] == Empty {
				blackControl := kingMovesBlack[x][y]
				whiteControl := kingMovesWhite[x][y]

				if blackControl == 1 && whiteControl == 1 {
					// 如果双方都可以控制这个格子，使用min和max函数处理差值
					diff := float64(whiteControl - blackControl)
					p2 += math.Min(1.0, math.Max(-1.0, diff/6.0))
				} else if blackControl == 1 {
					p2 += 1
				} else if whiteControl == 1 {
					p2 -= 1
				}
				// 注意，如果都不能到达该格子，则不计分
			}
		}
	}

	return p1, p2
}

// 灵活度
func (b *AmazonBoard) CalculateMobility() float64 {
	queenMovesBlack, queenMovesWhite := b.CalculateQueenMoves()
	kingMovesBlack, kingMovesWhite := b.CalculateKingMoves()
	var mobilityBlack, mobilityWhite float64

	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			if b[x][y] == Empty {
				mobility := float64(b.calculateMobilityForCell(x, y)) // 确保使用float64进行计算
				// 为避免除以零的情况，确保kingMoves中的值不是0，如果是0，意味着棋子无法到达，不应该计入灵活度
				if queenMovesBlack[x][y] != 100 && kingMovesBlack[x][y] != 0 {
					mobilityBlack += mobility / float64(kingMovesBlack[x][y])
				}
				if queenMovesWhite[x][y] != 100 && kingMovesWhite[x][y] != 0 {
					mobilityWhite += mobility / float64(kingMovesWhite[x][y])
				}
			}
		}
	}

	// 为避免除以零的错误，对方的灵活度值加上一个小数，确保除数不为零
	if mobilityWhite == 0 {
		mobilityWhite += 0.00001
	}
	if mobilityBlack == 0 {
		mobilityBlack += 0.00001
	}

	// 返回双方灵活度值的比值
	return mobilityBlack / mobilityWhite
}

// CalculateMobilityForCell 计算单个空格的灵活度值
func (b *AmazonBoard) calculateMobilityForCell(x, y int) int {
	mobility := 0
	for _, d := range dir {
		newX, newY := x+d[0], y+d[1]
		if newX >= 0 && newX < 10 && newY >= 0 && newY < 10 && b[newX][newY] == Empty {
			mobility++
		}
	}
	return mobility
}

// CalculateEvaluationValue 作为 AmazonBoard 的方法
func (b *AmazonBoard) CalculateEvaluationValue(turnID int, isBlackTurn bool) float64 {
	// 首先，计算tq, tk, p1, p2, 和mobility的值
	tqBlack, tqWhite := b.CalculateQueenTerritory()
	tkBlack, tkWhite := b.CalculateKingTerritory()
	queenMovesBlack, queenMovesWhite := b.CalculateQueenMoves()
	kingMovesBlack, kingMovesWhite := b.CalculateKingMoves()
	p1, p2 := b.CalculateP1P2(queenMovesBlack, queenMovesWhite, kingMovesBlack, kingMovesWhite)
	mobility := b.CalculateMobility()

	// 根据turnID调整每个要素的权重
	var k1, k2, k3, k4, k5 float64
	if turnID < 17 {
		k1 = 2 * (32 + float64(turnID)/2)
		k2 = 1 * (32 - 1.8*float64(turnID)/2)
		k3 = 1 * (32 - 1.8*float64(turnID)/2)
		k4 = 2 * (32 - 2.0*float64(turnID)/2)
		k5 = 0.5 * (32 - 1.8*float64(turnID)/2)
	} else {
		k1 = 5
		k2, k3, k4, k5 = 0, 0, 0, 0
	}

	// 计算整个棋局的评估值
	value := k1*(tqBlack-tqWhite) + k2*(tkBlack-tkWhite) + k3*p1 + k4*p2 + k5*mobility

	return value
}

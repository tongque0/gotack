package amazon

import (
	"math/rand"
	"time"
)

var zobristTable [10][10][4]uint64

func init() { // Go会自动调用此init函数
	initZobristTable()
}
func initZobristTable() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			for k := 0; k < 4; k++ {
				zobristTable[i][j][k] = rand.Uint64()
			}
		}
	}
}

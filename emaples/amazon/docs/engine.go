package docs

//引擎构建文件
// package main

// import (
// 	"GoTack/ai"
// 	"GoTack/games/amazon"
// 	"bufio"
// 	"fmt"
// 	"os"
// 	"strings"
// )

// const INF = 0x3f3f3f3f

// var (
// 	line  string
// 	step  int
// 	board *amazon.AmazonBoard
// 	color int
// )

// func main() {
// 	sc := bufio.NewScanner(os.Stdin)
// 	for sc.Scan() {
// 		line = sc.Text()
// 		if line == "name?" {
// 			fmt.Println("name GoTack-Amazon")
// 		} else if line == "quit" {
// 			os.Exit(0)
// 		} else if strings.HasPrefix(line, "new") {
// 			step = 1
// 			words := strings.Split(line, " ")
// 			board = amazon.NewBoard()
// 			if words[1] == "black" {
// 				color = amazon.Black
// 				runSearch()
// 			} else {
// 				color = amazon.White
// 			}
// 		} else if strings.HasPrefix(line, "move") {
// 			words := strings.Split(line, " ")
// 			move := words[1]
// 			board[move[3]-'A'][move[2]-'A'] = board[move[1]-'A'][move[0]-'A']
// 			board[move[1]-'A'][move[0]-'A'] = amazon.Empty
// 			board[move[5]-'A'][move[4]-'A'] = amazon.Arrow
// 			if !board.IsGameOver() {
// 				runSearch()
// 			}
// 		} else if line == "end" {
// 			fmt.Print("游戏结束")
// 		}
// 	}
// }
// func runSearch() {
// 	var IsMaxPlayer = true
// 	if color == 2 {
// 		IsMaxPlayer = false
// 	}
// 	var e *ai.AlphaBetaEvaluator
// 	if step < 12 {
// 		e = ai.NewAlphaBetaEvaluator(2, IsMaxPlayer, amazon.EvaluateFunc)
// 	} else {
// 		e = ai.NewAlphaBetaEvaluator(4, IsMaxPlayer, amazon.EvaluateFunc)
// 	}

// 	move := e.GetBestMove(board, step)
// 	m, ok := move.(amazon.AmazonMove)
// 	if !ok {
// 		fmt.Println(m)
// 		return
// 	}
// 	board.Move(move)
// 	fmt.Printf("move %c%c%c%c%c%c\n", m.From.Y+'A', m.From.X+'A', m.To.Y+'A', m.To.X+'A', m.Put.Y+'A', m.Put.X+'A')
// 	step++
// }

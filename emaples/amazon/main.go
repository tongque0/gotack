package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"tamazon/amazon"

	"github.com/tongque0/gotack"
)

const INF = 0x3f3f3f3f

var (
	line  string
	step  int
	board *amazon.AmazonBoard
	color int
)

func main() {
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		line = sc.Text()
		if line == "name?" {
			fmt.Println("name Tack-Amazon")
		} else if line == "quit" {
			os.Exit(0)
		} else if strings.HasPrefix(line, "new") {
			step = 1
			words := strings.Split(line, " ")
			board = amazon.NewBoard()
			if words[1] == "black" {
				color = amazon.Black
				runSearch()
			} else {
				color = amazon.White
			}
		} else if strings.HasPrefix(line, "move") {
			words := strings.Split(line, " ")
			move := words[1]
			board[move[3]-'A'][move[2]-'A'] = board[move[1]-'A'][move[0]-'A']
			board[move[1]-'A'][move[0]-'A'] = amazon.Empty
			board[move[5]-'A'][move[4]-'A'] = amazon.Arrow
			if !board.IsGameOver() {
				runSearch()
			}
		} else if line == "end" {
			continue
		} else {
			continue
		}
	}
}
func runSearch() {
	var IsMaxPlayer = true
	var e *gotack.Evaluator
	if color == 2 {
		IsMaxPlayer = false
	}
	if step < 12 {
		e = gotack.NewEvaluator(gotack.AlphaBeta, 2, IsMaxPlayer, func(board gotack.Board, isMaxPlayer bool, opts ...interface{}) float64 {
			return amazon.EvaluateFunc(board.(*amazon.AmazonBoard), isMaxPlayer, step)
		})
	} else {
		e = gotack.NewEvaluator(gotack.AlphaBeta, 4, IsMaxPlayer, func(board gotack.Board, isMaxPlayer bool, opts ...interface{}) float64 {
			return amazon.EvaluateFunc(board.(*amazon.AmazonBoard), isMaxPlayer, step)
		})
	}
	move := e.GetBestMove(board)
	m, ok := move.(amazon.AmazonMove)
	if !ok {
		return
	}
	board.Move(move)
	fmt.Printf("move %c%c%c%c%c%c\n", m.From.Y+'A', m.From.X+'A', m.To.Y+'A', m.To.X+'A', m.Put.Y+'A', m.Put.X+'A')
	step++
}

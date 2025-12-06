package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	boardSize   = 15
	emptySymbol = '.'
	humanSymbol = 'X'
	aiSymbol    = 'O'
)

type cell struct {
	row int
	col int
}

func main() {
	rand.Seed(time.Now().UnixNano())

	board := newBoard()
	fmt.Println("欢迎来到五子棋！你执子为 X，AI 执子为 O。")
	fmt.Println("输入格式: 行 列 (例如: 8 8)，范围 1-15。")

	reader := bufio.NewReader(os.Stdin)
	for {
		printBoard(board)
		humanMove(board, reader)
		if checkWin(board, humanSymbol) {
			printBoard(board)
			fmt.Println("你赢了！")
			return
		}
		if isDraw(board) {
			printBoard(board)
			fmt.Println("平局！")
			return
		}

		aiMove(board)
		if checkWin(board, aiSymbol) {
			printBoard(board)
			fmt.Println("AI 获胜！")
			return
		}
		if isDraw(board) {
			printBoard(board)
			fmt.Println("平局！")
			return
		}
	}
}

func newBoard() [][]rune {
	board := make([][]rune, boardSize)
	for i := range board {
		board[i] = make([]rune, boardSize)
		for j := range board[i] {
			board[i][j] = emptySymbol
		}
	}
	return board
}

func printBoard(board [][]rune) {
	fmt.Print("    ")
	for c := 1; c <= boardSize; c++ {
		fmt.Printf("%2d ", c)
	}
	fmt.Println()
	fmt.Println(strings.Repeat("---", boardSize+1))
	for r := 0; r < boardSize; r++ {
		fmt.Printf("%2d |", r+1)
		for c := 0; c < boardSize; c++ {
			fmt.Printf(" %c", board[r][c])
			if c < boardSize-1 {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func humanMove(board [][]rune, reader *bufio.Reader) {
	for {
		fmt.Print("你的回合，输入行 列: ")
		var r, c int
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("读取输入失败，请重试。")
			reader.Reset(os.Stdin)
			continue
		}
		line = strings.TrimSpace(line)
		if _, err := fmt.Sscanf(line, "%d %d", &r, &c); err != nil {
			fmt.Println("格式错误，请按 行 列 输入。")
			continue
		}
		r--
		c--
		if !isValidMove(board, r, c) {
			fmt.Println("位置无效或已被占用，请重新输入。")
			continue
		}
		board[r][c] = humanSymbol
		return
	}
}

func aiMove(board [][]rune) {
	move := chooseAIMove(board)
	board[move.row][move.col] = aiSymbol
	fmt.Printf("AI 落子: 行 %d 列 %d\n", move.row+1, move.col+1)
}

func chooseAIMove(board [][]rune) cell {
	emptyCells := availableCells(board)
	if len(emptyCells) == 0 {
		return cell{0, 0}
	}

	// 1. 自己能直接获胜
	if winCell, ok := findWinningMove(board, aiSymbol); ok {
		return winCell
	}

	// 2. 阻挡玩家即将获胜
	if blockCell, ok := findWinningMove(board, humanSymbol); ok {
		return blockCell
	}

	// 3. 优先中心
	center := boardSize / 2
	if board[center][center] == emptySymbol {
		return cell{center, center}
	}

	// 4. 选择靠近现有棋子的格子
	neighbors := prioritizedNeighbors(board)
	if len(neighbors) > 0 {
		return neighbors[rand.Intn(len(neighbors))]
	}

	// 5. 随机选择
	return emptyCells[rand.Intn(len(emptyCells))]
}

func findWinningMove(board [][]rune, symbol rune) (cell, bool) {
	for _, c := range availableCells(board) {
		board[c.row][c.col] = symbol
		win := checkWin(board, symbol)
		board[c.row][c.col] = emptySymbol
		if win {
			return c, true
		}
	}
	return cell{}, false
}

func prioritizedNeighbors(board [][]rune) []cell {
	candidates := []cell{}
	for r := 0; r < boardSize; r++ {
		for c := 0; c < boardSize; c++ {
			if board[r][c] != emptySymbol {
				continue
			}
			if hasNeighbor(board, r, c) {
				candidates = append(candidates, cell{r, c})
			}
		}
	}
	return candidates
}

func hasNeighbor(board [][]rune, r, c int) bool {
	for dr := -1; dr <= 1; dr++ {
		for dc := -1; dc <= 1; dc++ {
			if dr == 0 && dc == 0 {
				continue
			}
			nr := r + dr
			nc := c + dc
			if nr >= 0 && nr < boardSize && nc >= 0 && nc < boardSize {
				if board[nr][nc] != emptySymbol {
					return true
				}
			}
		}
	}
	return false
}

func availableCells(board [][]rune) []cell {
	cells := []cell{}
	for r := 0; r < boardSize; r++ {
		for c := 0; c < boardSize; c++ {
			if board[r][c] == emptySymbol {
				cells = append(cells, cell{r, c})
			}
		}
	}
	return cells
}

func isValidMove(board [][]rune, r, c int) bool {
	return r >= 0 && r < boardSize && c >= 0 && c < boardSize && board[r][c] == emptySymbol
}

func isDraw(board [][]rune) bool {
	for r := 0; r < boardSize; r++ {
		for c := 0; c < boardSize; c++ {
			if board[r][c] == emptySymbol {
				return false
			}
		}
	}
	return true
}

func checkWin(board [][]rune, symbol rune) bool {
	directions := []struct{ dr, dc int }{{0, 1}, {1, 0}, {1, 1}, {1, -1}}
	for r := 0; r < boardSize; r++ {
		for c := 0; c < boardSize; c++ {
			if board[r][c] != symbol {
				continue
			}
			for _, d := range directions {
				if countConsecutive(board, r, c, d.dr, d.dc, symbol) >= 5 {
					return true
				}
			}
		}
	}
	return false
}

func countConsecutive(board [][]rune, r, c, dr, dc int, symbol rune) int {
	count := 0
	cr, cc := r, c
	for cr >= 0 && cr < boardSize && cc >= 0 && cc < boardSize && board[cr][cc] == symbol {
		count++
		cr += dr
		cc += dc
	}
	return count
}

package main

import (
	"fmt"
)

const edge = 30

type Game struct {
	board [edge][edge]bool
}

func countNeighbours(g Game, x int, y int) int {
	count := 0

	if x > 0 && g.board[x-1][y] {
		count = count + 1
	}
	if x < edge && g.board[x+1][y] {
		count = count + 1
	}

	if y > 0 && g.board[x][y-1] {
		count = count + 1
	}
	if y < edge && g.board[x][y+1] {
		count = count + 1
	}

	if x > 0 && y > 0 && g.board[x-1][y-1] {
		count = count + 1
	}
	if x < edge && y > 0 && g.board[x+1][y-1] {
		count = count + 1
	}
	if x > 0 && y < edge && g.board[x-1][y+1] {
		count = count + 1
	}
	if x < edge && y < edge && g.board[x+1][y+1] {
		count = count + 1
	}

	return count
}

func iterate(g Game) Game {
	ret := Game{[edge][edge]bool{}}

	for i := 0; i < edge; i++ {
		for j := 0; j < edge; j++ {
			if countNeighbours(g, i, j) == 3 {
				ret.board[i][j] = true
			} else if g.board[i][j] && countNeighbours(g, i, j) == 2 {
				ret.board[i][j] = true
			} else {
				ret.board[i][j] = false
			}
		}
	}

	return ret
}

func main() {
	c := make(chan struct{}, 0)
	fmt.Println("Hello, WebAssembly!")
	// js.Global().Set("iterate", js.FuncOf(iterate))
	<-c
}

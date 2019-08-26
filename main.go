package main

import (
	"fmt"
	"syscall/js"
)

const edge = 10
const backgroundColor = "#111111"
const color = "#888888"

type Game struct {
	board [edge][edge]bool
}

func countNeighbours(g Game, x int, y int) int {
	count := 0

	if x > 0 && g.board[x-1][y] {
		count = count + 1
	}
	if x < edge-1 && g.board[x+1][y] {
		count = count + 1
	}

	if y > 0 && g.board[x][y-1] {
		count = count + 1
	}
	if y < edge-1 && g.board[x][y+1] {
		count = count + 1
	}

	if x > 0 && y > 0 && g.board[x-1][y-1] {
		count = count + 1
	}
	if x < edge-1 && y > 0 && g.board[x+1][y-1] {
		count = count + 1
	}
	if x > 0 && y < edge-1 && g.board[x-1][y+1] {
		count = count + 1
	}
	if x < edge-1 && y < edge-1 && g.board[x+1][y+1] {
		count = count + 1
	}

	return count
}

func update(g Game) Game {
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

func render(g Game) {
	canvases := js.Global().Get("document").Call("querySelectorAll", "[data-conways]")
	for i := 0; i < canvases.Length(); i++ {
		canvases.Index(i).Get("style").Set("background-color", backgroundColor)
		ctx := canvases.Index(i).Call("getContext", "2d")

		for x := 0; x < edge; x++ {
			for y := 0; y < edge; y++ {
				if g.board[x][y] {
					ctx.Call("rect", x*edge, y*edge, edge, edge)
					ctx.Set("fillStyle", color)
					ctx.Call("fill")
				}
			}
		}
	}
}

func main() {
	c := make(chan struct{}, 0)
	fmt.Println("Hello, WebAssembly!")

	game := Game{[edge][edge]bool{}}
	game.board[0][0] = true
	game.board[5][1] = true

	render(game)

	<-c
}

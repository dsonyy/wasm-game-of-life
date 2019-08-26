package main

import (
	"fmt"
	"math/rand"
	"syscall/js"
)

const backgroundColor = "#111111"
const color = "#ffffff"

type Game struct {
	board  [][]int
	width  int
	height int
	cell   int
}

// func countNeighbours(g Game, x int, y int) int {
// 	count := 0

// 	if x > 0 && g.board[x-1][y] {
// 		count = count + 1
// 	}
// 	if x < edge-1 && g.board[x+1][y] {
// 		count = count + 1
// 	}

// 	if y > 0 && g.board[x][y-1] {
// 		count = count + 1
// 	}
// 	if y < edge-1 && g.board[x][y+1] {
// 		count = count + 1
// 	}

// 	if x > 0 && y > 0 && g.board[x-1][y-1] {
// 		count = count + 1
// 	}
// 	if x < edge-1 && y > 0 && g.board[x+1][y-1] {
// 		count = count + 1
// 	}
// 	if x > 0 && y < edge-1 && g.board[x-1][y+1] {
// 		count = count + 1
// 	}
// 	if x < edge-1 && y < edge-1 && g.board[x+1][y+1] {
// 		count = count + 1
// 	}

// 	return count
// }

// func update(g Game) Game {
// 	ret := Game{[][]int{}}

// for i := 0; i < edge; i++ {
// 	for j := 0; j < edge; j++ {
// 		if countNeighbours(g, i, j) == 3 {
// 			ret.board[i][j] = true
// 		} else if g.board[i][j] && countNeighbours(g, i, j) == 2 {
// 			ret.board[i][j] = true
// 		} else {
// 			ret.board[i][j] = false
// 		}
// 	}
// }

// 	return ret
// }

func render(g Game) {
	//canvases := js.Global().Get("document").Call("querySelectorAll", "[data-conways]")

	// for i := 0; i < canvases.Length(); i++ {
	// 	canvases.Index(i).Get("style").Set("background-color", backgroundColor)
	// 	ctx := canvases.Index(i).Call("getContext", "2d")

	// 	for x := 0; x < edge; x++ {
	// 		for y := 0; y < edge; y++ {
	// 			if g.board[x][y] {
	// 				ctx.Call("rect", x*edge, y*edge, edge, edge)
	// 				ctx.Set("fillStyle", color)
	// 				ctx.Call("fill")
	// 			}
	// 		}
	// 	}
	// }
}

var games map[string]Game

func initCanvases() {
	canvases := js.Global().Get("document").Call("querySelectorAll", "[data-conways]")
	games = make(map[string]Game)

	for i := 0; i < canvases.Length(); i++ {
		cell := 10
		width := canvases.Index(i).Get("width").Int() / cell
		height := canvases.Index(i).Get("height").Int() / cell
		id := canvases.Index(i).Get("id").String()

		game := Game{[][]int{}, width, height, cell}
		for y := 0; y < int(height); y++ {
			game.board = append(game.board, []int{})
		}
		games[id] = game
	}
}

func fillCanvases(percentage int) {
	canvases := js.Global().Get("document").Call("querySelectorAll", "[data-conways]")

	for i := 0; i < canvases.Length(); i++ {
		id := canvases.Index(i).Get("id").String()
		for y := 0; y < int(games[id].height); y++ {
			for x := 0; x < int(games[id].width); x++ {
				if rand.Intn(100) < percentage {
					games[id].board[y] = append(games[id].board[y], x)
				}
			}
		}
	}

}

func updateCanvases() {

}

func renderCanvases() {
	canvases := js.Global().Get("document").Call("querySelectorAll", "[data-conways]")

	for i := 0; i < canvases.Length(); i++ {
		context := canvases.Index(i).Call("getContext", "2d")
		id := canvases.Index(i).Get("id").String()

		for y := 0; y < games[id].height; y++ {
			for i := 0; i < len(games[id].board[y]); i++ {
				edge := games[id].cell
				yPos := y * edge
				xPos := games[id].board[y][i] * edge

				context.Call("rect", xPos, yPos, edge, edge)
				context.Set("fillStyle", color)
				context.Call("fill")
			}
		}
	}
}

func main() {
	c := make(chan struct{}, 0)
	fmt.Println("Hello, WebAssembly!")

	initCanvases()
	fmt.Println(games)
	fillCanvases(10)
	fmt.Println(games)
	renderCanvases()

	<-c
}

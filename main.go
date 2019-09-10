package main

import (
	"fmt"
	"log"
	"math/rand"
	"syscall/js"
	"time"
)

const backgroundColor = "#555"
const color = "#eee"

type Game struct {
	board  [][]bool
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

var games map[int]Game

func initCanvases() {
	canvases := js.Global().Get("document").Call("querySelectorAll", "[data-conways]")
	games = make(map[int]Game)

	for i := 0; i < canvases.Length(); i++ {
		cell := 3
		width := canvases.Index(i).Get("width").Int() / cell
		height := canvases.Index(i).Get("height").Int() / cell

		game := Game{[][]bool{}, width, height, cell}
		for y := 0; y < height; y++ {
			game.board = append(game.board, []bool{})
			for x := 0; x < width; x++ {
				game.board[y] = append(game.board[y], false)
			}
		}

		games[i] = game
	}
}

func fillCanvases(percentage int) {
	canvases := js.Global().Get("document").Call("querySelectorAll", "[data-conways]")

	for i := 0; i < canvases.Length(); i++ {

		for y := 0; y < games[i].height; y++ {
			for x := 0; x < games[i].width; x++ {
				games[i].board[y][x] = percentage > rand.Intn(100)
			}
		}
	}

}

func countNeighbours(i int, x int, y int) int {
	count := 0

	if x > 0 && games[i].board[x-1][y] {
		count = count + 1
	}
	if x < games[i].width-1 && games[i].board[x+1][y] {
		count = count + 1
	}

	if y > 0 && games[i].board[x][y-1] {
		count = count + 1
	}

	if y < games[i].height-1 && games[i].board[x][y+1] {
		count = count + 1
	}

	if x > 0 && y > 0 && games[i].board[x-1][y-1] {
		count = count + 1
	}
	if x < games[i].width-1 && y > 0 && games[i].board[x+1][y-1] {
		count = count + 1
	}
	if x > 0 && y < games[i].height-1 && games[i].board[x-1][y+1] {
		count = count + 1
	}
	if x < games[i].width-1 && y < games[i].height-1 && games[i].board[x+1][y+1] {
		count = count + 1
	}

	return count
}

func updateCanvases() {

	for i := range games {

		updated := Game{[][]bool{}, games[i].width, games[i].height, games[i].cell}
		for y := 0; y < updated.height; y++ {
			updated.board = append(updated.board, []bool{})
			for x := 0; x < updated.width; x++ {
				updated.board[y] = append(updated.board[y], false)
			}
		}


		start := time.Now()
		for y := 0; y < games[i].height; y++ {
			for x := 0; x < games[i].width; x++ {
				neighbours := countNeighbours(i, x, y)
				if neighbours == 3 {
					updated.board[x][y] = true
				} else if games[i].board[x][y] && neighbours == 2 {
					updated.board[x][y] = true
				}
			}
		}
		games[i] = updated
		elapsed := time.Since(start)
		log.Printf("--> %s", elapsed)
	}

}

func renderCanvases() {
	canvases := js.Global().Get("document").Call("querySelectorAll", "[data-conways]")

	for i := 0; i < canvases.Length(); i++ {
		context := canvases.Index(i).Call("getContext", "2d")

		context.Set("fillStyle", backgroundColor)
		context.Call("fillRect", 0, 0, games[i].width*games[i].cell, games[i].height*games[i].cell)

		for y := 0; y < games[i].height; y++ {
			for x := 0; x < games[i].width; x++ {
				if games[i].board[x][y] == true {

					context.Set("fillStyle", color)
					context.Call("fillRect", x*games[i].cell, y*games[i].cell, games[i].cell, games[i].cell)
				}
			}
		}
	}
}

func loop(this js.Value, args []js.Value) interface{} {

	log.Println("AAA")
	updateCanvases()

	js.Global().Get("window").Call("requestAnimationFrame", js.FuncOf(loop))
	return js.Value{}
}

func main() {
	c := make(chan struct{}, 0)
	fmt.Println("Hello, WebAssembly!")

	initCanvases()
	fillCanvases(30)


	js.Global().Get("window").Call("requestAnimationFrame", js.FuncOf(loop))

	//for {
		// str := ""
		// for y := 0; y < games[0].height; y++ {
		// 	for x := 0; x < games[0].width; x++ {
		// 		if games[0].board[x][y] {
		// 			str += "[" + strconv.FormatInt(int64(countNeighbours(0, x, y)), 10) + "]"
		// 		} else {
		// 			str += " " + strconv.FormatInt(int64(countNeighbours(0, x, y)), 10) + " "
		// 		}
		// 	}
		// 	str += "\n"
		// }
		// fmt.Println(str)

		//updateCanvases()
		//renderCanvases()

		//time.Sleep(500 * time.Millisecond)
	//}

	//fmt.Println(games)
	//fillCanvases(10)

	//for i := int64(0); true; i++ {
	//renderCanvases()
	// rand.Seed(1)
	// fillCanvases(10)

	// renderCanvases()
	// rand.Seed(1)
	// fillCanvases(10)

	// renderCanvases()
	// rand.Seed(1)
	// fillCanvases(10)

	// renderCanvases()
	// rand.Seed(111)
	// fillCanvases(0)
	//}

	<-c
}

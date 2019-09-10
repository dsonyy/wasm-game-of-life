package main

import (
	"fmt"
	"log"
	"math/rand"
	"syscall/js"
)

const backgroundColor = "#555"
const color = "#eee"

type Game struct {
	board  [][]int8
	width  int
	height int
	cell   int
}

var games map[int]Game

func initCanvases() {
	canvases := js.Global().Get("document").Call("querySelectorAll", "[data-conways]")
	games = make(map[int]Game)

	for i := 0; i < canvases.Length(); i++ {
		cell := 3
		width := canvases.Index(i).Get("width").Int() / cell
		height := canvases.Index(i).Get("height").Int() / cell

		game := Game{[][]int8{}, width, height, cell}
		for y := 0; y < height; y++ {
			game.board = append(game.board, []int8{})
			for x := 0; x < width; x++ {
				game.board[y] = append(game.board[y], 0)
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
				if (percentage > rand.Intn(100)) {
					games[i].board[y][x] = 1
				}
			}
		}
	}

}

func countNeighbours(i int, x int, y int) int {
	count := 0

	if x > 0 && games[i].board[x-1][y] > 0 {
		count = count + 1
	}
	if x < games[i].width-1 && games[i].board[x+1][y] > 0 {
		count = count + 1
	}

	if y > 0 && games[i].board[x][y-1] > 0 {
		count = count + 1
	}

	if y < games[i].height-1 && games[i].board[x][y+1] > 0 {
		count = count + 1
	}

	if x > 0 && y > 0 && games[i].board[x-1][y-1] > 0 {
		count = count + 1
	}
	if x < games[i].width-1 && y > 0 && games[i].board[x+1][y-1] > 0 {
		count = count + 1
	}
	if x > 0 && y < games[i].height-1 && games[i].board[x-1][y+1] > 0 {
		count = count + 1
	}
	if x < games[i].width-1 && y < games[i].height-1 && games[i].board[x+1][y+1] > 0 {
		count = count + 1
	}

	return count
}


func updateCanvases() {

	for i := range games {

		updated := Game{[][]int8{}, games[i].width, games[i].height, games[i].cell}
		for y := 0; y < updated.height; y++ {
			updated.board = append(updated.board, []int8{})
			for x := 0; x < updated.width; x++ {
				updated.board[y] = append(updated.board[y], 0)
			}
		}


		// start := time.Now()
		for y := 0; y < games[i].height; y++ {
			for x := 0; x < games[i].width; x++ {
				neighbours := countNeighbours(i, x, y)
				if neighbours == 3 {
					updated.board[x][y] = 1
				} else if games[i].board[x][y] > 0 && neighbours == 2 {
					updated.board[x][y] = 1
				}
			}
		}
		games[i] = updated
		// elapsed := time.Since(start)
		// log.Printf("--> %s", elapsed)
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
				if games[i].board[x][y] > 0 {

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
	//renderCanvases()

	js.Global().Get("window").Call("requestAnimationFrame", js.FuncOf(loop))
	return js.Value{}
}

func main() {
	c := make(chan struct{}, 0)
	fmt.Println("Hello, WebAssembly!")

	initCanvases()
	fillCanvases(30)


	js.Global().Get("window").Call("requestAnimationFrame", js.FuncOf(loop))

	<-c
}

package main

import (
	"log"
	"math/rand"
	"syscall/js"
	"time"
)

type Game struct {
	board  [][]uint8
	width  int
	height int
	cell   int
	stopped bool
	color string
	backgroundColor string
}

func newGame(width int, height int, cell int) Game {
	return Game{[][]uint8{}, width, height, cell, false, "#eee", "#555"}
}

var games map[string]*Game

func initCanvases() {
	canvases := js.Global().Get("document").Call("querySelectorAll", "[data-conways]")
	games = make(map[string]*Game)

	for i := 0; i < canvases.Length(); i++ {
		cell := 3
		width := canvases.Index(i).Get("width").Int() / cell
		height := canvases.Index(i).Get("height").Int() / cell

		game := newGame(width, height, cell)
		for y := 0; y < height; y++ {
			game.board = append(game.board, []uint8{})
			for x := 0; x < width; x++ {
				game.board[y] = append(game.board[y], 0)
			}
		}

		games[canvases.Index(i).Get("id").String()] = &game
	}
}

func fillCanvases(percentage int) {
	canvases := js.Global().Get("document").Call("querySelectorAll", "[data-conways]")

	for i := 0; i < canvases.Length(); i++ {
		id := canvases.Index(i).Get("id").String()

		for y := 1; y < games[id].height - 1; y++ {
			for x := 1; x < games[id].width - 1; x++ {
				if (percentage < rand.Intn(100)) {
					games[id].board[y][x] += 100
					
					games[id].board[y][x - 1]++
					games[id].board[y][x + 1]++
					games[id].board[y - 1][x - 1]++
					games[id].board[y - 1][x]++
					games[id].board[y - 1][x + 1]++
					games[id].board[y + 1][x - 1]++
					games[id].board[y + 1][x]++
					games[id].board[y + 1][x + 1]++
				}
			}
		}
	}

}


func updateCanvases() {
	for i := range games {

		if games[i].stopped {
			continue
		}

		updated := newGame(games[i].width, games[i].height, games[i].cell)
		for y := 0; y < updated.height; y++ {
			updated.board = append(updated.board, []uint8{})
			for x := 0; x < updated.width; x++ {
				updated.board[y] = append(updated.board[y], 0)
			}
		}


		for y := 1; y < games[i].height - 1; y++ {
			for x := 1; x < games[i].width - 1; x++ { 
				if games[i].board[y][x] == 103 || games[i].board[y][x] == 102{
					updated.board[y][x] += 100
					
					updated.board[y][x - 1]++
					updated.board[y][x + 1]++
					updated.board[y - 1][x - 1]++
					updated.board[y - 1][x]++
					updated.board[y - 1][x + 1]++
					updated.board[y + 1][x - 1]++
					updated.board[y + 1][x]++
					updated.board[y + 1][x + 1]++

				} else if games[i].board[y][x] == 3 {
					updated.board[y][x] += 100
					
					updated.board[y][x - 1]++
					updated.board[y][x + 1]++
					updated.board[y - 1][x - 1]++
					updated.board[y - 1][x]++
					updated.board[y - 1][x + 1]++
					updated.board[y + 1][x - 1]++
					updated.board[y + 1][x]++
					updated.board[y + 1][x + 1]++	
				}
			}
		}
		
		games[i] = &updated

	}

}

func renderCanvases() {
	canvases := js.Global().Get("document").Call("querySelectorAll", "[data-conways]")

	for i := 0; i < canvases.Length(); i++ {
		id := canvases.Index(i).Get("id").String()
		start := time.Now()

		context := canvases.Index(i).Call("getContext", "2d")

		context.Set("fillStyle", games[id].backgroundColor)
		context.Call("fillRect", 0, 0, games[id].width*games[id].cell, games[id].height*games[id].cell)

		context.Set("fillStyle", games[id].color)
		for y := 0; y < games[id].height; y++ {
			for x := 0; x < games[id].width; x++ {
				if games[id].board[x][y] >= 100 {
					context.Call("fillRect", x*games[id].cell, y*games[id].cell, games[id].cell, games[id].cell)	 			
				}
			}
		}

		elapsed := time.Since(start)
		log.Printf("--> %s", elapsed)
	}
}

func stop(this js.Value, args []js.Value) interface{} {
	games[this.Get("id").String()].stopped = true
	return js.Value{}
}

func setColor(this js.Value, args []js.Value) interface{} {
	if len(args) > 0 && args[0].Type() == js.TypeString {
		games[this.Get("id").String()].color = args[0].String()
	}
	return js.Value{}
}

func setBackgroundColor(this js.Value, args []js.Value) interface{} {
	if len(args) > 0 && args[0].Type() == js.TypeString {
		games[this.Get("id").String()].backgroundColor = args[0].String()
	}
	return js.Value{}
}

func setMinInterval(this js.Value, args []js.Value) interface{} {
	return js.Value{} 
}


func resume(this js.Value, args []js.Value) interface{} {
	games[this.Get("id").String()].stopped = false
	return js.Value{}
}

func clear(this js.Value, args []js.Value) interface{} {
	return js.Value{} 
}

func spawn(this js.Value, args []js.Value) interface{} {
	return js.Value{} 
}

func kill(this js.Value, args []js.Value) interface{} {
	return js.Value{} 
}

func getCellSize(this js.Value, args []js.Value) interface{} {
	return js.Value{} 
}


func loop(this js.Value, args []js.Value) interface{} {

	updateCanvases()
	renderCanvases()

	js.Global().Get("window").Call("requestAnimationFrame", js.FuncOf(loop))
	return js.Value{}
}

func main() {
	c := make(chan struct{}, 0)
	log.Println("Hello, WebAssembly!")

	initCanvases()
	fillCanvases(70)

	canvases := js.Global().Get("document").Call("querySelectorAll", "[data-conways]")
	for i := 0; i < canvases.Length(); i++ {
		canvases.Index(i).Set("stop", js.FuncOf(stop))
		canvases.Index(i).Set("resume", js.FuncOf(resume))
		canvases.Index(i).Set("setColor", js.FuncOf(setColor))
	}

	js.Global().Get("window").Call("requestAnimationFrame", js.FuncOf(loop))

	<-c
}

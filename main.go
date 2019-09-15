package main

import (
	"log"
	"math/rand"
	"syscall/js"
	// "time"
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

func newGame(width int, height int, cell int, color string, backgroundColor string) Game {
	return Game{[][]uint8{}, width, height, cell, false, color, backgroundColor}
}

var games map[string]*Game

func initCanvases() {
	canvases := js.Global().Get("document").Call("querySelectorAll", "[data-conways]")
	games = make(map[string]*Game)

	for i := 0; i < canvases.Length(); i++ {
		cell := 3
		width := canvases.Index(i).Get("width").Int() / cell
		height := canvases.Index(i).Get("height").Int() / cell

		game := newGame(width, height, cell, "#eee", "#555")
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

		updated := newGame(games[i].width, games[i].height, games[i].cell, games[i].color, games[i].backgroundColor)
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
		// start := time.Now()

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

		// elapsed := time.Since(start)
		// log.Printf("--> %s", elapsed)
	}
}

func loop(this js.Value, args []js.Value) interface{} {

	updateCanvases()
	renderCanvases()

	js.Global().Get("window").Call("requestAnimationFrame", js.FuncOf(loop))
	return js.Value{}
}


func jsStop(this js.Value, args []js.Value) interface{} {
	games[this.Get("id").String()].stopped = true
	return js.Value{}
}

func jsResume(this js.Value, args []js.Value) interface{} {
	games[this.Get("id").String()].stopped = false
	return js.Value{}
}

func jsSetColor(this js.Value, args []js.Value) interface{} {
	if len(args) > 0 && args[0].Type() == js.TypeString {
		games[this.Get("id").String()].color = args[0].String()
	}
	return js.Value{}
}

func jsSetBackgroundColor(this js.Value, args []js.Value) interface{} {
	if len(args) > 0 && args[0].Type() == js.TypeString {
		games[this.Get("id").String()].backgroundColor = args[0].String()
	}
	return js.Value{}
}

func jsSetMinInterval(this js.Value, args []js.Value) interface{} {
	return js.Value{} 
}

func jsClear(this js.Value, args []js.Value) interface{} {
	id := this.Get("id").String()

	for y := 0; y < games[id].height; y++ {
		for x := 0; x < games[id].width; x++ {
			games[id].board[y][x] = 0
		}
	}

	return js.Value{} 
}

func jsSpawn(this js.Value, args []js.Value) interface{} {
	id := this.Get("id").String()
	
	if len(args) >= 2 && args[0].Type() == js.TypeNumber && args[1].Type() == js.TypeNumber && 
	args[0].Int() >= 0 && args[0].Int() < games[id].width && args[1].Int() >= 0 && args[1].Int() < games[id].height && games[id].board[args[0].Int()][args[1].Int()] < 100 {
		x := args[0].Int()
		y := args[1].Int()
		
		games[id].board[x][y] += 100

		games[id].board[y][x - 1]++
		games[id].board[y][x + 1]++
		games[id].board[y - 1][x - 1]++
		games[id].board[y - 1][x]++
		games[id].board[y - 1][x + 1]++
		games[id].board[y + 1][x - 1]++
		games[id].board[y + 1][x]++
		games[id].board[y + 1][x + 1]++
	}

	return js.Value{} 
}

func jsKill(this js.Value, args []js.Value) interface{} {
	return js.Value{} 
}


func jsGetWidthInPx(this js.Value, args []js.Value) interface{} {
	return games[this.Get("id").String()].width * games[this.Get("id").String()].cell
}

func jsGetHeightInPx(this js.Value, args []js.Value) interface{} {
	return games[this.Get("id").String()].height * games[this.Get("id").String()].cell
}

func jsGetWidthInCells(this js.Value, args []js.Value) interface{} {
	return games[this.Get("id").String()].width
}

func jsGetHeightInCells(this js.Value, args []js.Value) interface{} {
	return games[this.Get("id").String()].height
}

func jsGetColor(this js.Value, args []js.Value) interface{} {
	return games[this.Get("id").String()].color
}

func jsGetBackgroundColor(this js.Value, args []js.Value) interface{} {
	return games[this.Get("id").String()].backgroundColor
}

func jsGetCellSize(this js.Value, args []js.Value) interface{} {
	return games[this.Get("id").String()].cell
}

func jsIsStopped(this js.Value, args []js.Value) interface{} {
	return games[this.Get("id").String()].stopped
}


func main() {
	c := make(chan struct{}, 0)
	log.Println("Hello, WebAssembly!")

	initCanvases()
	fillCanvases(70)

	canvases := js.Global().Get("document").Call("querySelectorAll", "[data-conways]")
	for i := 0; i < canvases.Length(); i++ {

		canvases.Index(i).Set("getWidthInPx", js.FuncOf(jsGetWidthInPx))
		canvases.Index(i).Set("getHeightInPx", js.FuncOf(jsGetHeightInPx))
		canvases.Index(i).Set("getWidthInCells", js.FuncOf(jsGetWidthInCells))
		canvases.Index(i).Set("getHeightInCells", js.FuncOf(jsGetHeightInCells))
		canvases.Index(i).Set("getColor", js.FuncOf(jsGetColor))
		canvases.Index(i).Set("getBackgroundColor", js.FuncOf(jsGetBackgroundColor))
		canvases.Index(i).Set("getCellSize", js.FuncOf(jsGetCellSize))
		canvases.Index(i).Set("isStopped", js.FuncOf(jsIsStopped))

		canvases.Index(i).Set("setColor", js.FuncOf(jsSetColor))
		canvases.Index(i).Set("setBackgroundColor", js.FuncOf(jsSetBackgroundColor))

		canvases.Index(i).Set("stop", js.FuncOf(jsStop))
		canvases.Index(i).Set("resume", js.FuncOf(jsResume))

		canvases.Index(i).Set("clear", js.FuncOf(jsClear))
		canvases.Index(i).Set("spawn", js.FuncOf(jsSpawn))
		canvases.Index(i).Set("kill", js.FuncOf(jsKill))		
	}

	js.Global().Get("window").Call("requestAnimationFrame", js.FuncOf(loop))

	<-c
}

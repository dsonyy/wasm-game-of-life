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
	game := Game{[][]uint8{}, width, height, cell, false, color, backgroundColor}
	for x := 0; x < game.width; x++ {
		game.board = append(game.board, []uint8{})
		for y := 0; y < game.height; y++ {
			game.board[x] = append(game.board[x], 0)
		}
	}
	return game
}

var games map[string]*Game

// Initializes every canvas with data-conways selector.
func initCanvases() {
	canvases := js.Global().Get("document").Call("querySelectorAll", "[data-conways]")
	games = make(map[string]*Game)

	for i := 0; i < canvases.Length(); i++ {
		cell := 3
		width := 2 + canvases.Index(i).Get("width").Int() / cell
		height := 2 + canvases.Index(i).Get("height").Int() / cell

		game := newGame(width, height, cell, "#eee", "#555")
		for x := 0; x < width; x++ {
			game.board = append(game.board, []uint8{})
			for y := 0; y < height; y++ {
				game.board[x] = append(game.board[x], 0)
			}
		}

		games[canvases.Index(i).Get("id").String()] = &game
	}
}

// Creates alive cell with x, y coordinates in game[id]
func birth(id string, x int, y int) {
	games[id].board[x][y] += 100
					
	games[id].board[x][y - 1]++
	games[id].board[x][y + 1]++
	games[id].board[x - 1][y - 1]++
	games[id].board[x - 1][y]++
	games[id].board[x - 1][y + 1]++
	games[id].board[x + 1][y - 1]++
	games[id].board[x + 1][y]++
	games[id].board[x + 1][y + 1]++
}

// Kills cell with x, y coordinates in game[id]
func kill(id string, x int, y int) {
	games[id].board[x][y] -= 100
					
	games[id].board[x][y - 1]--
	games[id].board[x][y + 1]--
	games[id].board[x - 1][y - 1]--
	games[id].board[x - 1][y]--
	games[id].board[x - 1][y + 1]--
	games[id].board[x + 1][y - 1]--
	games[id].board[x + 1][y]--
	games[id].board[x + 1][y + 1]--
}

func fillCanvases(percentage int) {
	canvases := js.Global().Get("document").Call("querySelectorAll", "[data-conways]")

	for i := 0; i < canvases.Length(); i++ {
		id := canvases.Index(i).Get("id").String()

		for y := 1; y < games[id].height - 1; y++ {
			for x := 1; x < games[id].width - 1; x++ {
				if (percentage < rand.Intn(100)) {
					games[id].board[x][y] += 100
					games[id].board[x][y - 1]++
					games[id].board[x][y + 1]++
					games[id].board[x - 1][y - 1]++
					games[id].board[x - 1][y]++
					games[id].board[x - 1][y + 1]++
					games[id].board[x + 1][y - 1]++
					games[id].board[x + 1][y]++
					games[id].board[x + 1][y + 1]++
				}
			}
		}
	}

}

// Updates every canvas with data-conways selector.
func updateCanvases() {
	for i := range games {

		if games[i].stopped {
			continue
		}

		updated := newGame(games[i].width, games[i].height, games[i].cell, games[i].color, games[i].backgroundColor)

		for y := 1; y < games[i].height - 1; y++ {
			for x := 1; x < games[i].width - 1; x++ { 
				if games[i].board[x][y] == 103 || games[i].board[x][y] == 102 || games[i].board[x][y] == 3 {
					updated.board[x][y] += 100
					
					updated.board[x][y - 1]++
					updated.board[x][y + 1]++
					updated.board[x - 1][y - 1]++
					updated.board[x - 1][y]++
					updated.board[x - 1][y + 1]++
					updated.board[x + 1][y - 1]++
					updated.board[x + 1][y]++
					updated.board[x + 1][y + 1]++
				}
			}
		}
		
		games[i] = &updated
	}

}

// Renders every canvas with data-conways selector.
func renderCanvases() {
	canvases := js.Global().Get("document").Call("querySelectorAll", "[data-conways]")

	for i := 0; i < canvases.Length(); i++ {
		id := canvases.Index(i).Get("id").String()
		// start := time.Now()

		context := canvases.Index(i).Call("getContext", "2d")

		context.Set("fillStyle", games[id].backgroundColor)
		context.Call("fillRect", 0, 0, games[id].width*games[id].cell, games[id].height*games[id].cell)

		context.Set("fillStyle", games[id].color)
		for y := 1; y < games[id].height - 1; y++ {
			for x := 1; x < games[id].width - 1; x++ {
				if games[id].board[x][y] >= 100 {
					context.Call("fillRect", (x - 1)*games[id].cell, (y - 1)*games[id].cell, games[id].cell, games[id].cell)	 			
				}
			}
		}

		// elapsed := time.Since(start)
		// log.Printf("--> %s", elapsed)
	}
}

// Main game loop.
func loop(this js.Value, args []js.Value) interface{} {

	updateCanvases()
	renderCanvases()

	js.Global().Get("window").Call("requestAnimationFrame", js.FuncOf(loop))
	return js.Value{}
}


// Stops game.
func jsStop(this js.Value, args []js.Value) interface{} {
	games[this.Get("id").String()].stopped = true
	return js.Value{}
}

// Resumes game.
func jsResume(this js.Value, args []js.Value) interface{} {
	games[this.Get("id").String()].stopped = false
	return js.Value{}
}

// Sets game cells color. Recieves string argument with color.
func jsSetColor(this js.Value, args []js.Value) interface{} {
	if len(args) > 0 && args[0].Type() == js.TypeString {
		games[this.Get("id").String()].color = args[0].String()
	}
	return js.Value{}
}

// Sets game backgorund color. Recieves string argument with color.
func jsSetBackgroundColor(this js.Value, args []js.Value) interface{} {
	if len(args) > 0 && args[0].Type() == js.TypeString {
		games[this.Get("id").String()].backgroundColor = args[0].String()
	}
	return js.Value{}
}

func jsSetMinInterval(this js.Value, args []js.Value) interface{} {
	return js.Value{} 
}

// Fills game board with killed cells
func jsClear(this js.Value, args []js.Value) interface{} {
	id := this.Get("id").String()

	for y := 0; y < games[id].height; y++ {
		for x := 0; x < games[id].width; x++ {
			games[id].board[y][x] = 0
		}
	}

	log.Println(games[id])
	return js.Value{} 
}

// Creates alive cell. Recieves X and Y coordinates of cell.
func jsBirth(this js.Value, args []js.Value) interface{} {
	id := this.Get("id").String()

	if len(args) >= 2 && args[0].Type() == js.TypeNumber && args[1].Type() == js.TypeNumber {
		x := args[0].Int() + 1
		y := args[1].Int() + 1

		if x >= 1 && y >= 1 && x < games[id].width - 1 && y < games[id].height - 1 &&
		games[id].board[y][x] < 100 {
			birth(id, x, y)
		}
	}

	return js.Value{} 
}

// Kills cell. Recieves X and Y coordinates of cell.
func jsKill(this js.Value, args []js.Value) interface{} {
	id := this.Get("id").String()

	if len(args) >= 2 && args[0].Type() == js.TypeNumber && args[1].Type() == js.TypeNumber {
		x := args[0].Int() + 1
		y := args[1].Int() + 1

		if x >= 1 && y >= 1 && x < games[id].width - 1 && y < games[id].height - 1 &&
		games[id].board[y][x] >= 100 {
			kill(id, x, y)
		}
	}

	return js.Value{} 
}

// Returns number of cell neighbours. Recieves X and Y coordinates of cell.
func jsGetNeighbours(this js.Value, args []js.Value) interface{} {
	id := this.Get("id").String()

	if len(args) >= 2 && args[0].Type() == js.TypeNumber && args[1].Type() == js.TypeNumber {
		x := args[0].Int() + 1
		y := args[1].Int() + 1

		if x >= 1 && y >= 1 && x < games[id].width - 1 && y < games[id].height - 1 {
			if games[id].board[y][x] >= 100 {
				return games[id].board[y][x] - 100
			} else {
				return games[id].board[y][x]
			}
		}
	}

	return js.Value{} 
}

// Returns true if cell is alive, otherwise returns false. Recieves X and Y coordinates of cell.
func jsGet(this js.Value, args []js.Value) interface{} {
	id := this.Get("id").String()

	if len(args) >= 2 && args[0].Type() == js.TypeNumber && args[1].Type() == js.TypeNumber {
		x := args[0].Int() + 1
		y := args[1].Int() + 1

		if x >= 1 && y >= 1 && x < games[id].width - 1 && y < games[id].height - 1{
			if games[id].board[y][x] >= 100 {
				return true
			} else {
				return false
			}
		}
	}

	return js.Value{} 
}

// Returns used width of game board in pixels.
func jsGetWidthInPx(this js.Value, args []js.Value) interface{} {
	return (games[this.Get("id").String()].width - 2) * games[this.Get("id").String()].cell
}

// Returns used height of game board in pixels.
func jsGetHeightInPx(this js.Value, args []js.Value) interface{} {
	return (games[this.Get("id").String()].height - 2) * games[this.Get("id").String()].cell
}

// Returns width of game board in cells.
func jsGetWidthInCells(this js.Value, args []js.Value) interface{} {
	return games[this.Get("id").String()].width - 2
}

// Returns height of game board in cells.
func jsGetHeightInCells(this js.Value, args []js.Value) interface{} {
	return games[this.Get("id").String()].height - 2
}

// Returns game cells color.
func jsGetColor(this js.Value, args []js.Value) interface{} {
	return games[this.Get("id").String()].color
}

// Returns game background color.
func jsGetBackgroundColor(this js.Value, args []js.Value) interface{} {
	return games[this.Get("id").String()].backgroundColor
}

// Returns edge size of single game cell in pixels.
func jsGetCellSize(this js.Value, args []js.Value) interface{} {
	return games[this.Get("id").String()].cell
}

// Returns true if game is stopped. Otherwise returns false.
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
		canvases.Index(i).Set("birth", js.FuncOf(jsBirth))
		canvases.Index(i).Set("kill", js.FuncOf(jsKill))		
		canvases.Index(i).Set("get", js.FuncOf(jsGet))		
		canvases.Index(i).Set("getNeighbours", js.FuncOf(jsGetNeighbours))		
	}

	js.Global().Get("window").Call("requestAnimationFrame", js.FuncOf(loop))

	<-c
}

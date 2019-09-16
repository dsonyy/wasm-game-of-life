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
		cell := 4
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
	for id := range games {
		// start := time.Now()

		context := js.Global().Get("document").Call("getElementById", id).Call("getContext", "2d")

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
	if game, exist := games[this.Get("id").String()] ; exist && len(args) > 0 && args[0].Type() == js.TypeString {
		game.backgroundColor = args[0].String()
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

// Fills game board with random noice cells. Recieves number of alive cells probability.
func jsNoise(this js.Value, args []js.Value) interface{} {
	id := this.Get("id").String()

	percentage := 50

	if len(args) >= 1 && args[0].Type() == js.TypeNumber {
		percentage = args[0].Int()
	}

	log.Println(len(games))

	for y := 1; y < games[id].height - 1; y++ {
		for x := 1; x < games[id].width - 1; x++ {
			if (percentage < rand.Intn(100)) {
				birth(id, x, y)
			}
		}
	}

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

func jsStartGameOfLife(this js.Value, args []js.Value) interface{} {
	if len(args) == 0 && args[0].Type() != js.TypeString {
		return js.Value{}
	}

	id := args[0].String()

	canvas := js.Global().Get("document").Call("getElementById", id)

	if canvas.Type() != js.TypeObject {
		return js.Value{}
	}

	canvas.Set("noise", js.FuncOf(jsNoise))

	canvas.Set("getWidthInPx", js.FuncOf(jsGetWidthInPx))
	canvas.Set("getHeightInPx", js.FuncOf(jsGetHeightInPx))
	canvas.Set("getWidthInCells", js.FuncOf(jsGetWidthInCells))
	canvas.Set("getHeightInCells", js.FuncOf(jsGetHeightInCells))
	canvas.Set("getColor", js.FuncOf(jsGetColor))
	canvas.Set("getBackgroundColor", js.FuncOf(jsGetBackgroundColor))
	canvas.Set("getCellSize", js.FuncOf(jsGetCellSize))
	canvas.Set("isStopped", js.FuncOf(jsIsStopped))
	canvas.Set("setColor", js.FuncOf(jsSetColor))
	canvas.Set("setBackgroundColor", js.FuncOf(jsSetBackgroundColor))
	canvas.Set("stop", js.FuncOf(jsStop))
	canvas.Set("resume", js.FuncOf(jsResume))
	canvas.Set("clear", js.FuncOf(jsClear))
	canvas.Set("birth", js.FuncOf(jsBirth))
	canvas.Set("kill", js.FuncOf(jsKill))		
	canvas.Set("get", js.FuncOf(jsGet))		
	canvas.Set("getNeighbours", js.FuncOf(jsGetNeighbours))	

	cell := 10
	width := 2 + canvas.Get("width").Int() / cell
	height := 2 + canvas.Get("height").Int() / cell
	game := newGame(width, height, cell, "#eee", "#555")
	games[id] = &game

	log.Println(games[id].board)

	return canvas
}


func main() {
	c := make(chan struct{}, 0)
	games = make(map[string]*Game)

	js.Global().Set("startGameOfLife", js.FuncOf(jsStartGameOfLife))

	js.Global().Get("window").Call("requestAnimationFrame", js.FuncOf(loop))

	<-c
}

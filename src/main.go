package main

import (
	"log"
	"math/rand"
	"syscall/js"
)

// Struct representing a single instance of Game of Life.
type Game struct {
	board  [][]uint8
	width  int
	height int
	cell   int
	stopped bool
	color string
	backgroundColor string
}

// newGame is a constructor for Game object.
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

// games is a dictionary containing every game data. Keys are ids of canvases. Values are Game objects.
var games map[string]*Game

// birth sets the cell with coordinates (x, y) of the game with given id as active.
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

// kill sets the cell with coordinates (x, y) of the game with given id as inactive.
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

// updateCanvases updates every game.
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

// renderCanvases renders every game to proper canvas.
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

// jsStop is a javascript function but cannot be called directly from javascript.
// Updates and renders every game. Calls itself after that.
func jsLoop(this js.Value, args []js.Value) interface{} {

	updateCanvases()
	renderCanvases()

	js.Global().Get("window").Call("requestAnimationFrame", js.FuncOf(jsLoop))
	return js.Value{}
}

// jsStop is a function called from javascript on an initialized canvas element.
// Stops updating the game.
func jsStop(this js.Value, args []js.Value) interface{} {
	games[this.Get("id").String()].stopped = true
	return js.Value{}
}

// jsResume is a function called from javascript on an initialized canvas element.
// Resumes updating the game.
func jsResume(this js.Value, args []js.Value) interface{} {
	games[this.Get("id").String()].stopped = false
	return js.Value{}
}

// jsSetColor is a function called from javascript on an initialized canvas element.
// Sets the game cells color.
// As arg[0] it recieves a valid CSS color value (javascript string).
func jsSetColor(this js.Value, args []js.Value) interface{} {
	if len(args) > 0 && args[0].Type() == js.TypeString {
		games[this.Get("id").String()].color = args[0].String()
	}
	return js.Value{}
}

// jsSetBackgroundColor is a function called from javascript on an initialized canvas element.
// Sets the game backgorund color.
// As arg[0] it recieves a valid CSS color value (javascript string).
func jsSetBackgroundColor(this js.Value, args []js.Value) interface{} {
	if game, exist := games[this.Get("id").String()] ; exist && len(args) > 0 && args[0].Type() == js.TypeString {
		game.backgroundColor = args[0].String()
	}
	return js.Value{}
}

// jsSetMinInterval is a function called from javascript on an initialized canvas element.
// Sets the minimal time interval between two calls of rendering funcion 
// (note that if the game board was big, the time to calculate and render would be longer that this interval) 
func jsSetMinInterval(this js.Value, args []js.Value) interface{} {
	return js.Value{} 
}

// jsClear is a function called from javascript on an initialized canvas element.
// Fills the game board with inactive cells
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

// jsBirth is a function called from javascript on an initialized canvas element.
// Sets the cell as active.
// As arg[0] it recieves x position in cells of the canvas (javascript number).
// As arg[1] it recieves y position in cells of the canvas (javascript number).
// Returns true if the given cell exist. Otherwise returns javascript undefined.
func jsBirth(this js.Value, args []js.Value) interface{} {
	id := this.Get("id").String()

	if len(args) >= 2 && args[0].Type() == js.TypeNumber && args[1].Type() == js.TypeNumber {
		x := args[0].Int() + 1
		y := args[1].Int() + 1

		if x >= 1 && y >= 1 && x < games[id].width - 1 && y < games[id].height - 1 &&
		games[id].board[y][x] < 100 {
			birth(id, x, y)
			return true
		}
	}

	return js.Value{} 
}

// jsKill is a function called from javascript on an initialized canvas element.
// Sets the cell as inactive.
// As arg[0] it recieves x position in cells of the canvas (javascript number).
// As arg[1] it recieves y position in cells of the canvas (javascript number).
// Returns true if the given cell exist. Otherwise returns javascript undefined.
func jsKill(this js.Value, args []js.Value) interface{} {
	id := this.Get("id").String()

	if len(args) >= 2 && args[0].Type() == js.TypeNumber && args[1].Type() == js.TypeNumber {
		x := args[0].Int() + 1
		y := args[1].Int() + 1

		if x >= 1 && y >= 1 && x < games[id].width - 1 && y < games[id].height - 1 &&
		games[id].board[y][x] >= 100 {
			kill(id, x, y)
			return true
		}
	}

	return js.Value{} 
}

// jsGetNeighbours is a function called from javascript on an initialized canvas element.
// As arg[0] it recieves x position in cells of the canvas (javascript number).
// As arg[1] it recieves y position in cells of the canvas (javascript number).
// Returns number of neighbours of the given cell. If fails, returns undefined.
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

// jsGet is a function called from javascript on an initialized canvas element.
// As arg[0] it recieves x position in cells of the canvas (javascript number).
// As arg[1] it recieves y position in cells of the canvas (javascript number).
// Returns true if the given cell is alive, otherwise returns false. If fails, returns undefined.
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

// jsGetWidthInPx is a function called from javascript on an initialized canvas element. 
// Returns used width of game board in pixels (note that it is not always equals to the width of the canvas).
func jsGetWidthInPx(this js.Value, args []js.Value) interface{} {
	return (games[this.Get("id").String()].width - 2) * games[this.Get("id").String()].cell
}

// jsGetHeightInPx is a function called from javascript on an initialized canvas element. 
// Returns used height of game board in pixels (note that it is not always equals to the height of the canvas).
func jsGetHeightInPx(this js.Value, args []js.Value) interface{} {
	return (games[this.Get("id").String()].height - 2) * games[this.Get("id").String()].cell
}

// jsGetWidthInCells is a function called from javascript on an initialized canvas element. 
// Returns width of the game board in cells.
func jsGetWidthInCells(this js.Value, args []js.Value) interface{} {
	return games[this.Get("id").String()].width - 2
}

// jsGetHeightInCells is a function called from javascript on an initialized canvas element. 
// Returns height of the game board in cells.
func jsGetHeightInCells(this js.Value, args []js.Value) interface{} {
	return games[this.Get("id").String()].height - 2
}

// jsGetColor is a function called from javascript on an initialized canvas element. 
// Returns the game cells color.
func jsGetColor(this js.Value, args []js.Value) interface{} {
	return games[this.Get("id").String()].color
}

// jsGetBackgroundColor is a function called from javascript on an initialized canvas element. 
// Returns the game background color.
func jsGetBackgroundColor(this js.Value, args []js.Value) interface{} {
	return games[this.Get("id").String()].backgroundColor
}

// jsGetCellSize is a function called from javascript on an initialized canvas element. 
// Returns the edge size of a single game cell in pixels.
func jsGetCellSize(this js.Value, args []js.Value) interface{} {
	return games[this.Get("id").String()].cell
}

// jsIsStopped is a function called from javascript on an initialized canvas element. 
// Returns true if the game is stopped. Otherwise returns false.
func jsIsStopped(this js.Value, args []js.Value) interface{} {
	return games[this.Get("id").String()].stopped
}

// jsStartGameOfLife is a function called from javascript. Initializes Game of Life for the canvas with given id.
// As arg[0] it recieves id of the canvas (javascript string).
// Returns javascript undefined if cannot initialize the canvas. Otherwise returns the initialized canvas.
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

	return canvas
}

// jsEndGameOfLife is a function called from javascript. Finishes and cleans up Game of Life for the canvas with given id.
// As arg[0] it recieves id of the canvas (javascript string).
// Returns javascript undefined if cannot finish the canvas. Otherwise returns finished canvas.
func jsEndGameOfLife(this js.Value, args []js.Value) interface{} {
	return js.Value{}
}


// main is called when the wasm code starts.
func main() {
	c := make(chan struct{}, 0)
	games = make(map[string]*Game)

	js.Global().Set("startGameOfLife", js.FuncOf(jsStartGameOfLife))

	js.Global().Get("window").Call("requestAnimationFrame", js.FuncOf(jsLoop))

	<-c
}

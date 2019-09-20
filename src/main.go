package main

import (
	"math/rand"
	"syscall/js"
	"time"
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
	interval int64
	nextRefresh int64
	toRender bool
}

// newGame is a constructor for Game object.
func newGame(width int, height int, cell int, color string, backgroundColor string, interval int64, nextRefresh int64) Game {
	game := Game{[][]uint8{}, width, height, cell, false, color, backgroundColor, interval, nextRefresh, false}
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

// updates the game board with given id.
func update(id string) {
	updated := newGame(games[id].width, games[id].height, games[id].cell, games[id].color, games[id].backgroundColor, games[id].interval, games[id].nextRefresh)

	for y := 1; y < games[id].height - 1; y++ {
		for x := 1; x < games[id].width - 1; x++ { 
			if games[id].board[x][y] == 103 || games[id].board[x][y] == 102 || games[id].board[x][y] == 3 {
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
	
	games[id] = &updated
	games[id].toRender = true
}

// renders the game board with given id to proper canvas.
func render(id string) {
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

	games[id].toRender = false
}



// jsStop is a javascript function but cannot be called directly from javascript.
// Updates and renders every game. Calls itself after that.
func jsLoop(this js.Value, args []js.Value) interface{} {

	for id := range games {

		if games[id].nextRefresh <= time.Now().UnixNano() {
			if !games[id].stopped {
				update(id)
			}
	
			if games[id].toRender {
				render(id)
			}

			games[id].nextRefresh = time.Now().UnixNano() + games[id].interval
		}

	}

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

// jsSetSeed is a function called from javascript on an initialized canvas element.
// Sets random seed for entire Game of Life (if no arguments sets UTF time as the seed).
// As optional arg[0] it recieves the seed (javascript number).
func jsSetSeed(this js.Value, args []js.Value) interface{} {
	seed := time.Now().UTC().UnixNano()

	if len(args) > 0 && args[0].Type() == js.TypeNumber {
		seed = int64(args[0].Int())
	}

	rand.Seed(seed)
	return seed
}

// jsSetColor is a function called from javascript on an initialized canvas element.
// Sets the game cells color.
// As arg[0] it recieves a valid CSS color value (javascript string).
func jsSetColor(this js.Value, args []js.Value) interface{} {
	if len(args) > 0 && args[0].Type() == js.TypeString {
		games[this.Get("id").String()].color = args[0].String()
		games[this.Get("id").String()].toRender = true
	}
	return js.Value{}
}

// jsSetBackgroundColor is a function called from javascript on an initialized canvas element.
// Sets the game backgorund color.
// As arg[0] it recieves a valid CSS color value (javascript string).
func jsSetBackgroundColor(this js.Value, args []js.Value) interface{} {
	if game, exist := games[this.Get("id").String()] ; exist && len(args) > 0 && args[0].Type() == js.TypeString {
		game.backgroundColor = args[0].String()
		game.toRender = true
	}
	return js.Value{}
}

// jsSetMinInterval is a function called from javascript on an initialized canvas element.
// Sets the minimum interval between two calls of rendering funcion in milliseconds.
// (note that if the game board was big, the time to calculate and render would be longer that this interval)
// As arg[0] it recieves minimal interval in ms (javascript number).
func jsSetMinInterval(this js.Value, args []js.Value) interface{} {
	if len(args) == 0 || (len(args) >= 1 && args[0].Type() != js.TypeNumber) {
		return js.Value{}
	}

	id := this.Get("id").String()
	interval := 0
	if args[0].Int() > 0 {
		interval = args[0].Int()
	}
	games[id].interval = int64(interval * 1000000)
	games[id].nextRefresh = int64(0)

	return js.Value{} 
}

// jsClear is a function called from javascript on an initialized canvas element.
// Fills the game board with inactive cells
func jsClear(this js.Value, args []js.Value) interface{} {
	id := this.Get("id").String()

	for y := 0; y < games[id].height; y++ {
		for x := 0; x < games[id].width; x++ {
			games[id].board[x][y] = 0
		}
	}

	games[id].toRender = true
	return js.Value{} 
}

// jsRandomBirth is a function called from javascript on an initialized canvas element.
// Randomly sets active cells with a given probability in % (if no arguments sets 50%).
// As optional arg[0] it recieves probability in % (javascript number).
func jsRandomBirth(this js.Value, args []js.Value) interface{} {
	id := this.Get("id").String()

	percentage := 50

	if len(args) >= 1 && args[0].Type() == js.TypeNumber {
		percentage = args[0].Int()
	}

	for y := 1; y < games[id].height - 1; y++ {
		for x := 1; x < games[id].width - 1; x++ {
			if games[id].board[x][y] < 100 && percentage > rand.Intn(100) {
				birth(id, x, y)
			}
		}
	}

	games[id].toRender = true
	return js.Value{} 
}

// jsRandomKill is a function called from javascript on an initialized canvas element.
// Randomly sets inactive cells with a given probability in % (if no arguments sets 50%).
// As optional arg[0] it recieves probability in % (javascript number).
func jsRandomKill(this js.Value, args []js.Value) interface{} {
	id := this.Get("id").String()

	percentage := 50

	if len(args) >= 1 && args[0].Type() == js.TypeNumber {
		percentage = args[0].Int()
	}

	for y := 1; y < games[id].height - 1; y++ {
		for x := 1; x < games[id].width - 1; x++ {
			if games[id].board[x][y] >= 100 && percentage > rand.Intn(100) {
				kill(id, x, y)
			}
		}
	}

	games[id].toRender = true
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
		games[id].board[x][y] < 100 {
			birth(id, x, y)
			games[id].toRender = true
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
		games[id].board[x][y] >= 100 {
			kill(id, x, y)
			games[id].toRender = true
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
			if games[id].board[x][y] >= 100 {
				return games[id].board[x][y] - 100
			} else {
				return games[id].board[x][y]
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
			if games[id].board[x][y] >= 100 {
				return true
			} else {
				return false
			}
		}
	}

	return js.Value{} 
}

// jsGetMinInterval is a function called from javascript on an initialized canvas element.
// Returns the minimum interval between two calls of rendering funcion in milliseconds.
// (note that if the game board was big, the time to calculate and render would be longer that this interval)
func jsGetMinInterval(this js.Value, args []js.Value) interface{} {
	return games[this.Get("id").String()].interval / 1000000
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

// jsEndGameOfLife is a function called from javascript. Finishes and cleans up Game of Life for the canvas with given id.
// As arg[0] it recieves id of the canvas (javascript string).
// Returns javascript undefined if cannot finish the canvas. Otherwise returns finished canvas.
func jsEndGameOfLife(this js.Value, args []js.Value) interface{} {
	if len(args) == 0 || (len(args) >= 1 && args[0].Type() != js.TypeString){
		return js.Value{}
	}

	id := args[0].String()

	canvas := js.Global().Get("document").Call("getElementById", id)

	if canvas.Type() != js.TypeObject {
		return js.Value{}
	}

	canvas.Set("getWidthInPx", js.ValueOf(js.Value{}))
	canvas.Set("getHeightInPx", js.ValueOf(js.Value{}))
	canvas.Set("getWidthInCells", js.ValueOf(js.Value{}))
	canvas.Set("getHeightInCells", js.ValueOf(js.Value{}))
	canvas.Set("getColor", js.ValueOf(js.Value{}))
	canvas.Set("getBackgroundColor", js.ValueOf(js.Value{}))
	canvas.Set("getCellSize", js.ValueOf(js.Value{}))
	canvas.Set("getMinInterval", js.ValueOf(js.Value{}))
	canvas.Set("isStopped", js.ValueOf(js.Value{}))
	canvas.Set("setColor", js.ValueOf(js.Value{}))
	canvas.Set("setBackgroundColor", js.ValueOf(js.Value{}))
	canvas.Set("setMinInterval", js.ValueOf(js.Value{}))
	canvas.Set("stop", js.ValueOf(js.Value{}))
	canvas.Set("resume", js.ValueOf(js.Value{}))
	canvas.Set("clear", js.ValueOf(js.Value{}))
	canvas.Set("birth", js.ValueOf(js.Value{}))
	canvas.Set("kill", js.ValueOf(js.Value{}))
	canvas.Set("randomKill", js.ValueOf(js.Value{}))
	canvas.Set("randomBirth", js.ValueOf(js.Value{}))
	canvas.Set("get", js.ValueOf(js.Value{}))
	canvas.Set("getNeighbours", js.ValueOf(js.Value{}))

	delete(games, id)

	return js.Value{}
}

// jsStartGameOfLife is a function called from javascript. Initializes Game of Life for the canvas with given id.
// As arg[0] it recieves id of the canvas (javascript string).
// As optional arg[1] it recieves game cell size in pixels (javascript number).
// Returns javascript undefined if cannot initialize the canvas. Otherwise returns the initialized canvas.
func jsStartGameOfLife(this js.Value, args []js.Value) interface{} {
	if len(args) == 0 || (len(args) >= 1 && args[0].Type() != js.TypeString){
		return js.Value{}
	}

	id := args[0].String()

	if _, exist := games[id] ; exist {
		jsEndGameOfLife(this, args)
	}

	canvas := js.Global().Get("document").Call("getElementById", id)

	if canvas.Type() != js.TypeObject {
		return js.Value{}
	}

	canvas.Set("getWidthInPx", js.FuncOf(jsGetWidthInPx))
	canvas.Set("getHeightInPx", js.FuncOf(jsGetHeightInPx))
	canvas.Set("getWidthInCells", js.FuncOf(jsGetWidthInCells))
	canvas.Set("getHeightInCells", js.FuncOf(jsGetHeightInCells))
	canvas.Set("getColor", js.FuncOf(jsGetColor))
	canvas.Set("getBackgroundColor", js.FuncOf(jsGetBackgroundColor))
	canvas.Set("getCellSize", js.FuncOf(jsGetCellSize))
	canvas.Set("getMinInterval", js.FuncOf(jsGetMinInterval))
	canvas.Set("isStopped", js.FuncOf(jsIsStopped))
	canvas.Set("setColor", js.FuncOf(jsSetColor))
	canvas.Set("setBackgroundColor", js.FuncOf(jsSetBackgroundColor))
	canvas.Set("setMinInterval", js.FuncOf(jsSetMinInterval))
	canvas.Set("stop", js.FuncOf(jsStop))
	canvas.Set("resume", js.FuncOf(jsResume))
	canvas.Set("clear", js.FuncOf(jsClear))
	canvas.Set("birth", js.FuncOf(jsBirth))
	canvas.Set("kill", js.FuncOf(jsKill))		
	canvas.Set("randomKill", js.FuncOf(jsRandomKill))
	canvas.Set("randomBirth", js.FuncOf(jsRandomBirth))
	canvas.Set("get", js.FuncOf(jsGet))		
	canvas.Set("getNeighbours", js.FuncOf(jsGetNeighbours))	



	cell := 10
	if len(args) >= 2 && args[1].Type() == js.TypeNumber && args[1].Int() > 0 {
		cell = args[1].Int()
	}
	width := 2 + canvas.Get("width").Int() / cell
	height := 2 + canvas.Get("height").Int() / cell
	interval := 150
	if len(args) >= 3 && args[2].Type() == js.TypeNumber && args[2].Int() >= 0 {
		interval = args[2].Int() 
	}

	game := newGame(width, height, cell, "#eee", "#555", int64(interval * 1000000), 0)
	games[id] = &game

	return canvas
}


// main is called when the wasm code starts.
func main() {
	c := make(chan struct{}, 0)
	games = make(map[string]*Game)

	js.Global().Set("startGameOfLife", js.FuncOf(jsStartGameOfLife))
	js.Global().Set("endGameOfLife", js.FuncOf(jsEndGameOfLife))
	js.Global().Set("setGameOfLifeSeed", js.FuncOf(jsSetSeed))

	js.Global().Get("window").Call("requestAnimationFrame", js.FuncOf(jsLoop))

	<-c
}

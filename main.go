package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/fatih/color"
)

type grid struct {
	width, height int
	board         [][]int
}

var initStates map[string]func(gameBoard [][]int)

// draw renders the grid into an outputable string depiction
func (g *grid) draw() string {
	cmd := exec.Command("clear") //Linux example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()

	white := color.New(color.BgWhite)
	black := color.New(color.BgBlack)

	// draw the grid
	for row := range g.board {
		for column := range g.board[row] {
			if g.board[row][column] == 1 {
				white.Printf("o")
			} else {
				black.Printf(" ")
			}
		}
		fmt.Println()
	}
	return ""
}

// tick progresses time forwards by one tick
func (g *grid) tick() {
	newBoard := make([][]int, g.height) // initialise a slice of height slices
	for i := range newBoard {
		newBoard[i] = make([]int, g.width) // initialise a slice of width strings in each dy slice
	}

	// traverse the board and calculate each cells new state
	for row := 1; row < g.height-1; row++ {

		for column := 1; column < g.width-1; column++ {

			neighbours := calculateNeighbours(g.board, row, column)
			val := 0

			if g.board[row][column] == 1 && (neighbours == 2 || neighbours == 3) { //Each cell with two or three neighbors survives.
				val = 1
			} else if g.board[row][column] == 1 && neighbours < 2 { //Each cell with one or no neighbors dies, as if by solitude.
				val = 0
			} else if g.board[row][column] == 1 && neighbours > 3 { //Each cell with four or more neighbors dies, as if by overpopulation.
				val = 0
			} else if g.board[row][column] == 0 && neighbours == 3 { //Each cell with three neighbors becomes populated.
				val = 1
			} else {
				val = g.board[row][column]
			}

			newBoard[row][column] = val

		}

	}

	g.board = newBoard
}

func calculateNeighbours(board [][]int, target_row int, target_column int) int {
	neighbours := 0
	for row := -1; row <= 1; row++ {
		for column := -1; column <= 1; column++ {
			neighbours += board[row+target_row][column+target_column]
		}
	}
	neighbours -= board[target_row][target_column]
	return neighbours
}

func NewGrid(height int, width int, initFunc func(gameBoard [][]int)) grid {
	gameBoard := make([][]int, height) // initialise a slice of height slices
	for i := range gameBoard {
		gameBoard[i] = make([]int, width) // initialise a slice of width strings in each dy slice
	}

	initFunc(gameBoard)

	return grid{height: height, width: width, board: gameBoard}
}

func initialiseBoardRandom(gameBoard [][]int) {

	// initialise grid with random sequence
	for i := range gameBoard {
		for x := range gameBoard[i] {
			num := rand.Intn(13)
			if num >= 1 {
				gameBoard[i][x] = 0
			} else {
				gameBoard[i][x] = 1
			}
		}
	}

}

func initialiseBoardBar(gameBoard [][]int) {
	// calc half the height
	targHeight := len(gameBoard) / 2
	midPoint := len(gameBoard[0]) / 2

	gameBoard[targHeight][midPoint-1] = 1
	gameBoard[targHeight][midPoint-2] = 1
	gameBoard[targHeight][midPoint-3] = 1
	gameBoard[targHeight][midPoint-4] = 1
	gameBoard[targHeight][midPoint-5] = 1
	gameBoard[targHeight][midPoint] = 1
	gameBoard[targHeight][midPoint+1] = 1
	gameBoard[targHeight][midPoint+2] = 1
	gameBoard[targHeight][midPoint+3] = 1
	gameBoard[targHeight][midPoint+4] = 1
}

func initialiseBoardLightSpaceship(gameBoard [][]int) {
	// calc half the height
	targHeight := len(gameBoard) / 2
	midPoint := len(gameBoard[0]) / 2

	gameBoard[targHeight][midPoint] = 1
	gameBoard[targHeight-2][midPoint] = 1
	gameBoard[targHeight][midPoint+3] = 1
	gameBoard[targHeight-3][midPoint+1] = 1
	gameBoard[targHeight-3][midPoint+2] = 1
	gameBoard[targHeight-3][midPoint+3] = 1
	gameBoard[targHeight-3][midPoint+4] = 1
	gameBoard[targHeight-2][midPoint+4] = 1
	gameBoard[targHeight-1][midPoint+4] = 1
}

func main() {

	initStates = make(map[string]func(gameBoard [][]int))

	initStates["rand"] = initialiseBoardRandom
	initStates["bar"] = initialiseBoardBar
	initStates["spaceship"] = initialiseBoardLightSpaceship

	heightPtr := flag.Int("height", 30, "Grid Height")
	widthPtr := flag.Int("width", 50, "Grid Width")
	initState := flag.String("init", "rand", "rand|bar|spaceship")

	flag.Parse()

	g := NewGrid(*heightPtr, *widthPtr, initStates[*initState])
	g.draw()

	for {
		g.tick()
		time.Sleep(500 * time.Millisecond)
		g.draw()
	}

}

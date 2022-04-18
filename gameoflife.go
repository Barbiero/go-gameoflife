package main

import (
	"math/rand"
	"strings"
	"unsafe"
)

type State uint8

const (
	Dead State = iota
	Alive
)

type GameOfLife struct {
	// size of the GOL board
	width, height int

	// the board itself
	board [][]State

	boardReady chan []byte
}

type point struct {
	x, y int
}

func (g *GameOfLife) outOfBounds(x, y int) bool {
	return x < 0 || x >= g.width || y < 0 || y >= g.height
}

func (g *GameOfLife) numOfNeighbours(p point) int {
	count := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if g.outOfBounds(p.x+i, p.y+j) || (i == 0 && j == 0) {
				continue
			}
			if g.board[p.x+i][p.y+j] == Alive {
				count++
			}
		}
	}
	return count
}

func (g *GameOfLife) RunStep() {
	nextBoard := initBoard(g.width, g.height)
	for i := 0; i < g.width; i++ {
		for j := 0; j < g.height; j++ {
			nh := g.numOfNeighbours(point{i, j})

			elem := g.board[i][j]
			if elem == Dead && nh == 3 {
				nextBoard[i][j] = Alive
			} else if elem == Alive && nh <= 1 {
				nextBoard[i][j] = Dead
			} else if elem == Alive && nh >= 4 {
				nextBoard[i][j] = Dead
			} else {
				nextBoard[i][j] = g.board[i][j]
			}
		}
	}
	copy(g.board, nextBoard)
	g.boardReady <- g.flatBoard()
}

func initBoard(w, h int) [][]State {
	board := make([][]State, w)
	for i := 0; i < w; i++ {
		board[i] = make([]State, h)
		for j := 0; j < h; j++ {
			board[i][j] = Dead
		}
	}
	return board
}

func (g *GameOfLife) Init() {
	g.board = initBoard(g.width, g.height)

	for i := 0; i < g.width; i++ {
		for j := 0; j < g.height; j++ {
			g.board[i][j] = State(rand.Intn(2))
		}
	}
}

func (g *GameOfLife) String() string {
	builder := new(strings.Builder)
	for i := 0; i < g.width; i++ {
		for j := 0; j < g.height; j++ {
			if g.board[i][j] == Alive {
				builder.WriteString("1")
			} else {
				builder.WriteString("0")
			}
		}
		builder.WriteString("\n")
	}
	return builder.String()
}

func (g *GameOfLife) flatBoard() []byte {
	flatBoard := make([]State, g.width*g.height)
	for i := 0; i < g.width; i++ {
		copy(flatBoard[(i*g.width):(i+1)*g.width], g.board[i])
	}

	ptr := unsafe.Pointer(&flatBoard)
	byteArr := (*[]byte)(ptr)
	return *byteArr
}

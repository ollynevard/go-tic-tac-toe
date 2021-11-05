package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type Board [][]Token

type Token string
type GridRef string

type Player struct {
	name  string
	token Token
}

type Move struct {
	x int
	y int
}

const OToken Token = "O"
const XToken Token = "X"
const NullToken Token = " "

func newBoard(x int, y int, defaultToken Token) Board {
	b := make([][]Token, y)

	for i := 0; i < y; i++ {
		b[i] = make([]Token, x)

		for j := 0; j < x; j++ {
			b[i][j] = defaultToken
		}
	}

	return b
}

func clearScreen() {
	fmt.Print("\033[2J")
}

func moveCursorTopLeft() {
	fmt.Print("\033[H")
}

func printTitle() {
	title := `
▀█▀ █ █▀▀ ▄▄ ▀█▀ ▄▀█ █▀▀ ▄▄ ▀█▀ █▀█ █▀▀
░█░ █ █▄▄ ░░ ░█░ █▀█ █▄▄ ░░ ░█░ █▄█ ██▄
  `
	fmt.Print(title)
}

func render(b Board) {
	x := len(b[0]) - 1
	y := len(b) - 1

	clearScreen()
	moveCursorTopLeft()
	printTitle()

	fmt.Print("\n   ┌" + strings.Repeat("───┬", x) + "───┐\n")

	for i, row := range b {
		for j := range row {
			if j == 0 {
				fmt.Printf(" %d ", y-i+1)
			}

			fmt.Printf("│ %s ", b[y-i][j])

			if j == x {
				fmt.Print("│\n")
			}
		}

		if i != y {
			fmt.Print("   ├" + strings.Repeat("───┼", x) + "───┤\n")
		}
	}

	fmt.Print("   └" + strings.Repeat("───┴", x) + "───┘\n  ")

	for i := range b[0] {
		fmt.Printf("   %s", string(intToChar(i)))
	}

	fmt.Print("\n\n")
}

func getMove(b Board, p Player) Move {
	fmt.Printf("%s, what is your move (e.g. 'A1')?\n", p.name)

	var gr string

	for {
		fmt.Scan(&gr)

		if isValidGridRef(gr) {
			m := gridRefToMove(gr)

			if isValidMove(b, m) {
				return m
			}
		}

		fmt.Println("Enter a valid move")
	}
}

func isValidGridRef(s string) bool {
	re := regexp.MustCompile(`^[a-zA-Z][1-9][0-9]?$`)

	return re.FindString(s) != ""
}

func gridRefToMove(gr string) Move {
	r := []rune(gr)
	x := charToInt(r[0])
	y, _ := strconv.Atoi(string(r[1:]))

	return Move{x, y - 1}
}

func makeMove(b Board, m Move, p Player) Board {
	b[m.y][m.x] = p.token
	return b
}

func isValidMove(b Board, m Move) bool {
	return !isMoveOutOfBounds(b, m) && !isMoveTaken(b, m)
}

func isMoveTaken(b Board, m Move) bool {
	return b[m.y][m.x] != NullToken
}

func isMoveOutOfBounds(b Board, m Move) bool {
	return len(b) <= m.y || len(b[m.y]) <= m.x
}

func intToChar(i int) rune {
	return rune('A' + i)
}

func charToInt(r rune) int {
	return int(unicode.ToUpper(r) - 'A')
}

func isWinner(b Board) bool {
	lines := getAllLines(b)

	for _, l := range lines {
		if isWinningLine(b, l) {
			return true
		}
	}

	return false
}

func isWinningLine(b Board, l []Move) bool {
	for i := 1; i < len(l); i++ {
		if b[l[i].x][l[i].y] != b[l[0].x][l[0].y] || b[l[i].x][l[i].y] == NullToken {
			return false
		}
	}
	return true
}

func getAllLines(b Board) [][]Move {
	width := len(b[0])
	height := len(b)
	lines := [][]Move{}

	mainDiagonal := []Move{}
	for x := 0; x < width; x++ {
		column := []Move{}
		for y := 0; y < height; y++ {
			column = append(column, Move{x, y})
			if y+x == width-1 {
				mainDiagonal = append(mainDiagonal, Move{x, y})
			}
		}
		lines = append(lines, column)
	}
	lines = append(lines, mainDiagonal)

	antiDiagonal := []Move{}
	for y := 0; y < height; y++ {
		row := []Move{}
		for x := 0; x < width; x++ {
			row = append(row, Move{x, y})
			if y == x {
				antiDiagonal = append(antiDiagonal, Move{x, y})
			}
		}
		lines = append(lines, row)
	}
	lines = append(lines, antiDiagonal)

	return lines
}

func main() {
	player1 := Player{"Player 1", XToken}
	player2 := Player{"Player 2", OToken}

	b := newBoard(3, 3, NullToken)
	render(b)

	turn := 1

	for {
		p := player1
		if turn%2 == 0 {
			p = player2
		}
		m := getMove(b, p)
		b = makeMove(b, m, p)
		render(b)

		if isWinner(b) {
			fmt.Printf("%s is the winner!\n", p.name)
			break
		}

		turn++
	}
}

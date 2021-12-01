package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

//DEFINE GAME NECESARY VARIABLES
type sprite struct {
	row int
	col int
}

var player sprite
var ghosts []*sprite
var score int
var numDots int
var lives = 1

//LOADING SPRITES FROM JSON
type Config struct {
	Player   string `json:"player"`
	Ghost    string `json:"ghost"`
	Wall     string `json:"wall"`
	Dot      string `json:"dot"`
	Pill     string `json:"pill"`
	Death    string `json:"death"`
	Space    string `json:"space"`
	UseEmoji bool   `json:"use_emoji"`
}

//Colors for interface and borders

type Color int

const reset = "\x1b[0m"

const (
	BLACK Color = iota
	RED
	GREEN
	BROWN
	BLUE
	MAGENTA
	CYAN
	GREY
)

var colors = map[Color]string{
	BLACK:   "\x1b[1;30;40m",
	RED:     "\x1b[1;31;41m",
	GREEN:   "\x1b[1;32;42m",
	BROWN:   "\x1b[1;33;43m",
	BLUE:    "\x1b[1;34;44m",
	MAGENTA: "\x1b[1;35;45m",
	CYAN:    "\x1b[1;36;46m",
	GREY:    "\x1b[1;37;47m",
}

func blueBackground(text string) string {
	return "\x1b[44m" + text + reset
}

func blackBackground(text string, color Color) string {
	if c, ok := colors[color]; ok {
		return c + text + reset
	}
	return blueBackground(text)
}

//Reading  configuration file
var cfg Config

func loadConfig(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return err
	}
	return nil
}

func initSettings() {
	cbTerm := exec.Command("stty", "cbreak", "-echo")
	cbTerm.Stdin = os.Stdin

	err := cbTerm.Run()
	if err != nil {
		log.Fatalln("Unable to activate cbreak mode: ", err)
	}
}
func cleanup() {
	cookedTerm := exec.Command("stty", "-cbreak", "echo")
	cookedTerm.Stdin = os.Stdin

	err := cookedTerm.Run()
	if err != nil {
		log.Fatalln("Unable to restore cooked mode: ", err)
	}
}

func ClearScreen() {
	fmt.Print("\x1b[2J")
	moveCursorEmoji(0, 0)
}
func MoveCursor(row, col int) {
	fmt.Printf("\x1b[%d;%df", row+1, col+1)
}
func moveCursorEmoji(row, col int) {
	if cfg.UseEmoji {
		MoveCursor(row, col*2)
	} else {
		MoveCursor(row, col)
	}
}

func readInput() (string, error) {
	buffer := make([]byte, 100)

	cnt, err := os.Stdin.Read(buffer)
	if err != nil {
		return "", err
	}

	if cnt == 1 && buffer[0] == 0x1b {
		return "ESC", nil
	} else if cnt >= 3 {
		if buffer[0] == 0x1b && buffer[1] == '[' {
			switch buffer[2] {
			case 'A':
				return "UP", nil
			case 'B':
				return "DOWN", nil
			case 'C':
				return "RIGHT", nil
			case 'D':
				return "LEFT", nil
			}

		}
	}

	return "", nil
}

var level []string

func loadLevel(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		level = append(level, line)
	}

	//Finding the start position of the player
	for row, line := range level {
		for col, char := range line {
			switch char {
			case 'P':
				player = sprite{row, col}
			case 'G':
				ghosts = append(ghosts, &sprite{row, col})
			case '.':
				numDots++
			}
		}
	}

	return nil
}

var gui []string

func printScreen() {
	//ClearScreen()
	for _, line := range level {
		for _, chr := range line {
			switch chr {
			case '#':
				gui = append(gui, "#")
				//fmt.Print("#")
			case '.':
				gui = append(gui, ".")
				//fmt.Print(".")

				//fmt.Print("*")

			default:
				gui = append(gui, " ")
				//fmt.Print(" ")
			}
		}
		//fmt.Println()

	}
	fmt.Printf("%d", player.row)
	fmt.Printf("%d", player.col)
	gui[(player.row)*28+player.col] = "P"

	for _, ghost := range ghosts {
		fmt.Printf("%d", ghost.row)
		fmt.Printf("%d", ghost.col)
		gui[(ghost.row)*28+ghost.col] = "G"
	}

	//gui = append(gui,"P")
	//fmt.Print("P")

	moveCursorEmoji(len(level)+1, 0)
	fmt.Print("New row\n")

	for i := 0; i < len(gui); i++ {
		if i%28 == 0 {
			fmt.Print("\n")
		}
		fmt.Printf("%s", gui[i])

	}
	gui = nil
}

func makeMove(oldRow, oldCol int, dir string) (newRow, newCol int) {
	newRow, newCol = oldRow, oldCol

	switch dir {
	case "UP":

		newRow = newRow - 1
		if newRow < 0 {
			newRow = len(level) - 1
		}
	case "DOWN":

		newRow = newRow + 1
		if newRow == len(level) {
			newRow = 0
		}
	case "RIGHT":

		newCol = newCol + 1
		if newCol == len(level[0]) {
			newCol = 0
		}
	case "LEFT":

		newCol = newCol - 1
		if newCol < 0 {
			newCol = len(level[0]) - 1
		}
	}
	if level[newRow][newCol] == '#' {
		newRow = oldRow
		newCol = oldCol
	}

	return
}

func movePlayer(dir string) {
	player.row, player.col = makeMove(player.row, player.col, dir)
	switch level[player.row][player.col] {
	case '.':
		numDots--
		score++
		//Remove the dot from the level
		level[player.row] = level[player.row][0:player.col] + " " + level[player.row][player.col+1:]
	}
}

func directionGhost() string {
	dir := rand.Intn(4)
	move := map[int]string{
		0: "UP",
		1: "DOWN",
		2: "RIGHT",
		3: "LEFT",
	}
	return move[dir]
}

func moveGhosts() {
	for _, g := range ghosts {
		dir := directionGhost()
		g.row, g.col = makeMove(g.row, g.col, dir)
	}
}

func main() {

	initSettings()
	defer cleanup()

	//Loading the level from text file
	err := loadLevel("level1.txt")
	if err != nil {
		log.Println("Failed to load the level", err)
		return
	}

	//Reading config
	err = loadConfig("config.json")
	if err != nil {
		log.Println("Failed to load configuration", err)
		return
	}

	//Reading input from keyboard async
	input := make(chan string)
	go func(ch chan<- string) {
		for {
			input, err := readInput()
			if err != nil {
				log.Println("error reading input:", err)
				ch <- "ESC"
			}
			ch <- input
		}
	}(input)

	//Game loop
	for {

		printScreen()

		//Process player movement
		select { //Select is a switch but for channels
		case inp := <-input:
			if inp == "ESC" {
				lives = 0
			}
			movePlayer(inp)
		default:

		}

		moveGhosts() //Moves the gosts

		for _, g := range ghosts {
			if player == *g {
				lives = 0
			}
		}

		//Check for game over

		if numDots == 0 || lives <= 0 {
			if lives == 0 {
				moveCursorEmoji(player.row, player.col)
				fmt.Print(cfg.Death)
				moveCursorEmoji(len(level)+2, 0)
			}
			break
		}

		time.Sleep(300 * time.Millisecond)

		//fmt.Println("Hello, Pacman!")
		//break
	}
}

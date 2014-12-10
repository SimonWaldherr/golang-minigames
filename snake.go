package main

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

type Field struct {
	cells  [][]int
	width  int
	height int
}

type Snake struct {
	x int
	y int
}

var field *Field
var snake *Snake

var (
	setfps    int
	setwidth  int
	setheight int
	direction int
)

func newSnake() *Snake {
	return &Snake{x: 0, y: 0}
}

func (snake *Snake) set(x, y int) {
	snake.x = x
	snake.y = y
}

func (snake *Snake) move(x, y int) {
	snake.x = snake.x + x
	snake.y = snake.y + y
}

func newField(width, height int) *Field {
	cells := make([][]int, height)
	for cols := range cells {
		cells[cols] = make([]int, width)
	}
	return &Field{cells: cells, width: width, height: height}
}

func (field *Field) setVitality(x, y int, vitality int) {
	field.cells[y][x] = vitality
}

func (field *Field) getVitality(x, y int) int {
	x += field.width
	x %= field.width
	y += field.height
	y %= field.height
	return field.cells[y][x]
}

func generateFirstRound(width, height int) *Field {
	field := newField(width, height)
	field.setVitality(rand.Intn(width), rand.Intn(height), 1)
	return field
}

func (field *Field) nextRound(snake *Snake) *Field {
	new_field := field
	switch direction {
	case 1:
		snake.move(0, -1)
	case 2:
		snake.move(1, 0)
	case 3:
		snake.move(0, 1)
	case 4:
		snake.move(-1, 0)
	}
	new_field.setVitality(snake.x, snake.y, 1)
	return new_field
}

func (field *Field) printField() string {
	var buffer bytes.Buffer
	for y := 0; y < field.height; y++ {
		for x := 0; x < field.width; x++ {
			if field.getVitality(x, y) > 0 {
				buffer.WriteByte(byte('#'))
			} else {
				buffer.WriteByte(byte(' '))
			}
		}
		buffer.WriteByte('\n')
	}
	return buffer.String()
}

func getDirection() {
	defer func() {
		exec.Command("stty", "-f", "/dev/tty", "echo").Run()
	}()
	exec.Command("stty", "-f", "/dev/tty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-f", "/dev/tty", "-echo").Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Printf("Signal: %v\n", sig)
			exec.Command("stty", "-f", "/dev/tty", "echo").Run()
			os.Exit(0)
		}
	}()

	var b []byte = make([]byte, 1)
	for {
		os.Stdin.Read(b)
		kc := fmt.Sprintf("%v", b)
		switch kc {
		case "[65]", "[119]", "[107]", "up":
			direction = 1
		case "[67]", "[100]", "[108]", "right":
			direction = 2
		case "[66]", "[115]", "[106]", "down":
			direction = 3
		case "[68]", "[97]", "[104]", "left":
			direction = 4
		case "[113]", "quit":
			exec.Command("stty", "-f", "/dev/tty", "echo").Run()
			os.Exit(0)
		}
	}
}

func main() {
	defer func() {
		exec.Command("stty", "-f", "/dev/tty", "echo").Run()
	}()
	flag.IntVar(&setwidth, "w", 80, "terminal width")
	flag.IntVar(&setheight, "h", 20, "terminal height")
	flag.IntVar(&setfps, "f", 6, "frames per second")
	flag.Parse()

	field = generateFirstRound(setwidth, setheight)

	snake := newSnake()
	snake.set(rand.Intn(setwidth), rand.Intn(setheight))

	go getDirection()

	for {
		time.Sleep(time.Second / time.Duration(setfps))
		fmt.Print("\033[2J")
		field = field.nextRound(snake)
		str := field.printField()
		fmt.Print(str)
	}
}

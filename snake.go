package main

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"simonwaldherr.de/go/golibs/ansi"
	"simonwaldherr.de/go/golibs/as"
	"simonwaldherr.de/go/golibs/gcurses"
	"time"
)

type Field struct {
	cells  [][]int
	width  int
	height int
}

type Snake struct {
	x      int
	y      int
	length int
}

var field *Field
var snake *Snake
var gover bool

var (
	pts       int
	setfps    int
	setwidth  int
	setheight int
	direction int
)

func newSnake() *Snake {
	return &Snake{x: 0, y: 0, length: 8}
}

func (snake *Snake) set(x, y int) {
	snake.x = x
	snake.y = y
}

func (snake *Snake) get() (int, int) {
	return snake.x, snake.y
}

func (snake *Snake) len(l int) {
	snake.length += l
}

func (snake *Snake) mov(x, y int) {
	if snake.x+x >= setwidth {
		snake.x = 0
	} else if snake.x+x < 0 {
		snake.x = setwidth - 1
	} else {
		snake.x += x
	}
	if snake.y+y >= setheight {
		snake.y = 0
	} else if snake.y+y < 0 {
		snake.y = setheight - 1
	} else {
		snake.y += y
	}
}

func newField(width, height int) *Field {
	cells := make([][]int, height)
	for cols := range cells {
		cells[cols] = make([]int, width)
	}
	return &Field{cells: cells, width: width, height: height}
}

func (field *Field) set(x, y, vitality int) {
	field.cells[y][x] = vitality
}

func (field *Field) get(x, y int) int {
	x += field.width
	x %= field.width
	y += field.height
	y %= field.height
	return field.cells[y][x]
}

func generateFirstRound(width, height int) *Field {
	field := newField(width, height)
	field.nextRound(1)
	return field
}

func (field *Field) nextFrame(snake *Snake) *Field {
	new_field := field
	switch direction {
	case 1:
		snake.mov(0, -1)
	case 2:
		snake.mov(1, 0)
	case 3:
		snake.mov(0, 1)
	case 4:
		snake.mov(-1, 0)
	}
	x, y := snake.get()
	vit := field.get(x, y)
	if vit < 0 {
		pts++
		field.nextRound(pts)
	} else if vit > 0 && direction != 0 {
		end()
		gover = true
		return nil
	}
	new_field.set(snake.x, snake.y, snake.length)
	for y := 0; y < field.height; y++ {
		for x := 0; x < field.width; x++ {
			vit := field.get(x, y)
			if vit > 0 {
				new_field.set(x, y, vit-1)
			}
		}
	}
	return new_field
}

func (field *Field) nextRound(p int) {
	for {
		rand.Seed(time.Now().UnixNano())
		x := rand.Intn(setwidth-10) + 5
		y := rand.Intn(setheight-10) + 5
		if field.get(x, y) == 0 {
			field.set(x, y, p*-1)
			return
		}
	}
}

func (field *Field) print(gameover bool) string {
	var buffer bytes.Buffer
	var ptsstr string = fmt.Sprintf("Points: %v \n", pts)

	buffer.Write([]byte(ptsstr))

	for y := 0; y < field.height; y++ {
		for x := 0; x < field.width; x++ {
			if field.get(x, y) > 0 {
				if gameover {
					buffer.Write(as.Bytes(ansi.Color("█", ansi.Red)))
				} else {
					buffer.Write(as.Bytes(ansi.Color("█", ansi.Green)))
				}
			} else if field.get(x, y) < 0 {
				buffer.Write(as.Bytes(ansi.Color("█", ansi.Blue)))
			} else {
				buffer.WriteByte(byte(' '))
			}
		}
		buffer.WriteByte('\n')
	}
	return buffer.String()
}

func getDirection() {
	defer end()
	exec.Command("stty", "cbreak", "min", "1").Run()
	exec.Command("stty", "-f", "/dev/tty", "-echo").Run()
	exec.Command("stty", "-echo").Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Printf("Signal: %v\n", sig)
			end()
		}
	}()

	var b []byte = make([]byte, 1)
	for {
		os.Stdin.Read(b)
		kc := fmt.Sprintf("%v", b)
		switch kc {
		case "[65]", "[119]", "[107]", "up":
			if direction != 3 {
				direction = 1
			}
		case "[67]", "[100]", "[108]", "right":
			if direction != 4 {
				direction = 2
			}
		case "[66]", "[115]", "[106]", "down":
			if direction != 1 {
				direction = 3
			}
		case "[68]", "[97]", "[104]", "left":
			if direction != 2 {
				direction = 4
			}
		case "[113]", "quit":
			end()
		}
	}
}

func end() {
	fmt.Print(field.print(true))
	time.Sleep(time.Second / time.Duration(setfps) * 10)
	fmt.Printf("GameOver!!!\nYour Score: %v Points\n", pts)
	exec.Command("stty", "-f", "/dev/tty", "echo").Run()
	exec.Command("stty", "echo").Run()
	os.Exit(0)
}

func main() {
	defer end()
	flag.IntVar(&setwidth, "w", 80, "terminal width")
	flag.IntVar(&setheight, "h", 20, "terminal height")
	flag.IntVar(&setfps, "f", 6, "frames per second")
	flag.Parse()

	writer := gcurses.New()
	writer.Start()

	setheight -= 2

	rand.Seed(time.Now().UnixNano())
	field = generateFirstRound(setwidth, setheight)

	snake := newSnake()
	snake.set(rand.Intn(setwidth), rand.Intn(setheight))

	direction = 0
	ite := 0
	pts = 0

	go getDirection()

	for !gover {
		time.Sleep(time.Second / time.Duration(setfps))
		field = field.nextFrame(snake)
		fmt.Fprintf(writer, "%v", field.print(false))
		if ite == 10 {
			snake.len(1)
			ite = 0
		}
		ite++
	}
	writer.Stop()
}

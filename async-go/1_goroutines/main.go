package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

func main() {
	me := &Turtle{name: "YOU"}
	turtles := []*Turtle{{name: "Tom"}, {name: "Anna"}, {name: "Leo"}, me}
	input := keyboard()

	var wg sync.WaitGroup
	for _, t := range turtles[:3] {
		wg.Add(1)
		go race(t, &wg)
	}

	frame := time.NewTicker(50 * time.Millisecond)
	for over := false; !over; {
		select {
		case <-input:
			step(me, 1)
		case <-frame.C:
			over = draw(turtles)
		}
	}
	wg.Wait()

	fmt.Println()
	medals := []string{"GOLD  ", "SILVER", "BRONZE", "4th   "}
	for i, name := range podium {
		fmt.Printf("  %s %s\n", medals[i], name)
	}
}

const trackLength = 30

type Turtle struct {
	name     string
	position int
}

var (
	lock   sync.Mutex
	podium []string
)

func step(t *Turtle, n int) bool {
	lock.Lock()
	defer lock.Unlock()
	if t.position < trackLength {
		t.position += n
		if t.position >= trackLength {
			t.position = trackLength
			podium = append(podium, t.name)
		}
	}
	return t.position >= trackLength
}

func race(t *Turtle, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		time.Sleep(time.Duration(120+rand.Intn(200)) * time.Millisecond)
		if step(t, 1+rand.Intn(2)) {
			return
		}
	}
}

func draw(turtles []*Turtle) bool {
	lock.Lock()
	defer lock.Unlock()

	fmt.Print("\033[H\033[J")
	for _, t := range turtles {
		face := "(o_o)"
		if t.name == "YOU" {
			face = "(^o^)"
		}
		road := strings.Repeat(".", t.position)
		fmt.Printf("%-6s %s%s%s|\n", t.name, road, face, strings.Repeat(" ", trackLength-t.position))
	}
	fmt.Printf("%-6s %s FINISH\n", "", strings.Repeat(" ", trackLength+5))
	return len(podium) == len(turtles)
}

func keyboard() <-chan string {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	ch := make(chan string)
	go func() {
		for sc := bufio.NewScanner(os.Stdin); sc.Scan(); {
			ch <- sc.Text()
		}
	}()
	return ch
}

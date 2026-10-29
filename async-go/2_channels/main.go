package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func main() {
	input := keyboard()

	toYou, toBot := make(chan Ball), make(chan Ball)
	miss, quit := make(chan string), make(chan bool)

	me := Player{"YOU", 0, width, toYou, toBot, human(input)}
	bot := Player{"BOT", width, 0, toBot, toYou, bot}

	go me.play(miss, quit)
	go bot.play(miss, quit)

	serve := toYou
	for score["YOU"] < 3 && score["BOT"] < 3 {
		draw(-1, "Get ready...")
		time.Sleep(time.Second)
		serve <- Ball{}
		loser := <-miss

		winner, next := "YOU", toYou
		if loser == "YOU" {
			winner, next = "BOT", toBot
		}
		score[winner]++
		serve = next
		draw(-1, loser+" missed!  Point for "+winner)
		time.Sleep(1500 * time.Millisecond)
	}
	close(quit)
	draw(-1, "*** "+map[bool]string{true: "YOU WIN!", false: "BOT WINS!"}[score["YOU"] == 3]+" ***")
}

const width = 40

type Ball struct{ hits int }

var score = map[string]int{}

func draw(pos int, msg string) {
	table := []byte(strings.Repeat(".", width+1))
	if pos >= 0 {
		table[pos] = 'O'
	}
	fmt.Printf("\033[H\033[JYOU %d : %d BOT\n\n", score["YOU"], score["BOT"])
	fmt.Printf("  YOU [|]%s[|] BOT\n\n  %s\n", table, msg)
}

func fly(from, to int, msg string) {
	d := 2
	if to < from {
		d = -2
	}
	for p := from; p != to+d; p += d {
		draw(p, msg)
		time.Sleep(40 * time.Millisecond)
	}
}

type brain func(hits int) bool

func bot(hits int) bool {
	time.Sleep(time.Duration(300+rand.Intn(400)) * time.Millisecond)
	return rand.Intn(100) >= 20
}

func human(input <-chan string) brain {
	return func(hits int) bool {
		for drained := false; !drained; {
			select {
			case <-input:
			default:
				drained = true
			}
		}
		draw(0, "The ball is here! Hit it back!")
		select {
		case <-input:
			return true
		case <-time.After(time.Duration(1500-100*min(hits, 9)) * time.Millisecond):
			return false
		}
	}
}

type Player struct {
	name       string
	home, away int
	in         <-chan Ball
	out        chan<- Ball
	think      brain
}

func (p Player) play(miss chan<- string, quit <-chan bool) {
	for {
		select {
		case ball := <-p.in:
			if !p.think(ball.hits) {
				miss <- p.name
				continue
			}
			ball.hits++
			fly(p.home, p.away, fmt.Sprintf("%s hits it!   rally: %d", p.name, ball.hits))
			p.out <- ball
		case <-quit:
			return
		}
	}
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

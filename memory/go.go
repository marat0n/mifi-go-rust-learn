package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

// TODO: РЕАЛИЗОВАТЬ ЭТУ ИГРУ НА ЛЮБОМ ДРУГОМ ЯЗЫКЕ

const (
	blobSize = 1 << 20
	limitMB  = 200
	goalSecs = 60
	leakRate = 2
	bars     = "▁▂▃▄▅▆▇█"
)

type Blob struct{ data []byte }

var (
	alive   int64
	held    []*Blob
	history []int
	msg     = "кэш течёт! используйте команды ниже"
	start   = time.Now()
)

func newBlob() *Blob {
	b := &Blob{data: make([]byte, blobSize)}
	atomic.AddInt64(&alive, 1)

	runtime.SetFinalizer(b, func(*Blob) { atomic.AddInt64(&alive, -1) })
	return b
}

func spawn(n int, keep bool) {
	for i := 0; i < n; i++ {
		b := newBlob()
		if keep {
			held = append(held, b)
		}
	}
}

func sparkline() string {
	var sb strings.Builder
	r := []rune(bars)
	for _, v := range history {
		sb.WriteRune(r[min(max(v*(len(r)-1)/limitMB, 0), len(r)-1)])
	}
	return sb.String()
}

// render перерисовывает экран и возвращает true, когда игра окончена.
func render() bool {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	mb, secs := int(m.HeapAlloc>>20), int(time.Since(start).Seconds())
	fmt.Print("\033[H\033[2J")
	fmt.Println("=== УБОРЩИК МУСОРА (GARBAGE COLLECTOR) ===")
	fmt.Printf("Куча занята      : %4d МБ / лимит %d МБ   (прошло %d/%d с)\n", mb, limitMB, secs, goalSecs)
	fmt.Printf("История кучи     : %s\n", sparkline())
	fmt.Printf("Объектов живо, GC их ещё не освободил : %d\n", atomic.LoadInt64(&alive))
	fmt.Printf("Объектов с ССЫЛКАМИ (held)             : %d\n", len(held))
	fmt.Printf("Циклов GC: %d   последняя пауза: %v\n\n", m.NumGC, time.Duration(m.PauseNs[(m.NumGC+255)%256]))
	fmt.Println("[s] СОЗДАТЬ+ДЕРЖАТЬ  +10 объектов в `held`: достижимы, GC обязан их сохранить")
	fmt.Println("[t] СОЗДАТЬ ВРЕМЕННЫЕ +10 объектов без ссылок: недостижимы, GC может их удалить")
	fmt.Println("[d] СБРОСИТЬ ССЫЛКИ  held = nil: объекты недостижимы, но ещё НЕ освобождены")
	fmt.Println("[g] ЗАПУСТИТЬ GC     runtime.GC(): реально освобождает недостижимые объекты")
	fmt.Println("[q] выход            (экран обновляется 2 раза в секунду; печатайте и жмите Enter)")
	fmt.Printf("\n> %s\n\nкоманда: ", msg)

	if mb > limitMB {
		fmt.Println("\n\nВЫ ПРОИГРАЛИ: куча взорвалась. GC не освобождает то, на что ещё есть ссылки!")
		return true
	}
	if secs >= goalSecs {
		fmt.Println("\n\nВЫ ПОБЕДИЛИ: куча под контролем!")
		return true
	}
	return false
}

func main() {
	input := make(chan string)
	go func() {
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			input <- strings.TrimSpace(sc.Text())
		}
		close(input)
	}()

	tick := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case cmd, ok := <-input:
			if !ok || cmd == "q" || cmd == "й" {
				return
			}
			switch cmd {
			case "s", "ы":
				spawn(10, true)
				msg = "+10 объектов СОХРАНЕНО в `held`: куча растёт, GC не поможет"
			case "t", "е":
				spawn(10, false)
				msg = "+10 ВРЕМЕННЫХ объектов: они стали мусором; следите за следующим циклом GC"
			case "d", "в":
				before := len(held)
				held = nil
				msg = fmt.Sprintf("сброшено ссылок: %d. Объекты недостижимы, но освободит их только цикл GC", before)
			case "g", "п":
				runtime.GC()
				msg = "runtime.GC() выполнен: недостижимые объекты освобождены"
			default:
				msg = "неизвестная команда: " + cmd
			}
			if render() {
				return
			}
		case <-tick.C:
			spawn(leakRate, true)
			spawn(2, false)
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			if history = append(history, int(m.HeapAlloc>>20)); len(history) > 50 {
				history = history[1:]
			}
			if render() {
				return
			}
		}
	}
}

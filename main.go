package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
)

func main() {
	var i uint64

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh)

	lineCh := make(chan string)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			lineCh <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}()

	for {
		select {
		case s := <-signalCh:
			log.Println("Got signal:", s)
			final := atomic.LoadUint64(&i)
			log.Println("Checkpoint:", final)
			os.Exit(0)
		case line := <-lineCh:
			fmt.Println(line)
			atomic.AddUint64(&i, 1)
		}
	}
}

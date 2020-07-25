package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
)

func main() {
	var (
		fileName   string
		skip       uint64
		checkpoint uint64
		canPrint   bool
	)

	flag.StringVar(&fileName, "file", "checkpoint_data", "checkpoint data file")
	flag.StringVar(&fileName, "f", "checkpoint_data", "checkpoint data file")

	flag.Parse()

	content, err := ioutil.ReadFile(fileName)
	if err != nil{
		log.Print(err)
        content = []byte("0")
	}

	skip, err = strconv.ParseUint(string(content), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	if skip == 0 {
		canPrint = true
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	lineCh := make(chan string)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			lineCh <- scanner.Text()
		}
		close(lineCh)
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}()

	for {
		select {
		case s := <-signalCh:
			log.Println("Got signal:", s)
			final := atomic.LoadUint64(&checkpoint)
			log.Println("Checkpoint:", final)
            content := []byte(strconv.FormatUint(final, 10))
			if err := ioutil.WriteFile(fileName, content, 0644); err != nil {
				log.Fatal(err)
			}
			os.Exit(0)
		case line, more := <-lineCh:
			if !more {
				os.Exit(0)
			}
			if canPrint {
				fmt.Println(line)
			}
			current := atomic.AddUint64(&checkpoint, 1)
			if !canPrint && current > skip {
				canPrint = true
			}
		}
	}
}

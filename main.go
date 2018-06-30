package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Game struct {
	falseAlarmTimes []time.Time
}

var game Game
var falseAlarmEvents = make(chan struct{})
var falseAlarmTimesChan = make(chan []time.Time)

func GameLoop() {
	for {
		select {
		case <-falseAlarmEvents:
			log.Println("got a false alarm")
			game.falseAlarmTimes = append(game.falseAlarmTimes, time.Now())
		case falseAlarmTimesChan <- game.falseAlarmTimes:
			log.Println("told world about times")
		}
	}
}

func HandleRequest(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		times := <-falseAlarmTimesChan
		for _, t := range times {
			fmt.Fprintln(w, time.Since(t))
		}
	}
}

func main() {
	fmt.Println("hi")
	go GameLoop()
	go (func() {
		for {
			_, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			falseAlarmEvents <- struct{}{}
		}
	})()
	http.HandleFunc("/", HandleRequest)
	log.Fatal(http.ListenAndServe(":2344", nil))
}

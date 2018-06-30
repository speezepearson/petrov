package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type PlayerName string

type PlayerBoard struct {
	falseAlarmTimes []time.Time
}

type Game struct {
	Boards map[PlayerName]*PlayerBoard
}

var game = Game{
	Boards: map[PlayerName]*PlayerBoard{
		"Alice": &PlayerBoard{[]time.Time{}},
		"Bob":   &PlayerBoard{[]time.Time{}},
	},
}
var falseAlarmEvents = make(chan PlayerName)
var gameStateChan = make(chan Game)

func GameLoop() {
	for {
		select {
		case victimName := <-falseAlarmEvents:
			board, ok := game.Boards[victimName]
			log.Println("got a false alarm for '", victimName, "'")
			if !ok {
				log.Println("...who doesn't exist")
				continue
			}
			board.falseAlarmTimes = append(board.falseAlarmTimes, time.Now())
		case gameStateChan <- game:
			log.Println("told world about times")
		}
	}
}

func HandleRequest(w http.ResponseWriter, req *http.Request) {
	log.Println("got request:", *req)
	if req.Method == "GET" {
		g := <-gameStateChan
		askerName := PlayerName(strings.TrimLeft(req.URL.Path, "/"))
		board, ok := g.Boards[askerName]
		if !ok {
			log.Println("request is for player", req.URL.Path, ", who doesn't exist")
			w.WriteHeader(404)
			return
		}
		for _, t := range board.falseAlarmTimes {
			fmt.Fprintln(w, time.Since(t))
		}
	}
}

func main() {
	fmt.Println("hi")
	go GameLoop()
	go (func() {
		for {
			line, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			player := strings.TrimSpace(line)
			falseAlarmEvents <- PlayerName(player)
		}
	})()
	http.HandleFunc("/", HandleRequest)
	log.Fatal(http.ListenAndServe(":2344", nil))
}

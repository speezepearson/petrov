package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

type PlayerName string

type PlayerBoard struct {
	falseAlarmTimes []time.Time
	launchedTime    *time.Time
}

const (
	MissileFlightTime = 1 * time.Minute
)

type Game struct {
	Boards map[PlayerName]*PlayerBoard
}

var game = Game{
	Boards: map[PlayerName]*PlayerBoard{
		"Alice": &PlayerBoard{},
		"Bob":   &PlayerBoard{},
	},
}
var mutex sync.Mutex

func addFalseAlarm(victimName PlayerName) {
	mutex.Lock()
	defer mutex.Unlock()
	board, ok := game.Boards[victimName]
	log.Println("got a false alarm for '", victimName, "'")
	if !ok {
		log.Println("...who doesn't exist")
		return
	}
	board.falseAlarmTimes = append(board.falseAlarmTimes, time.Now())
}

type FuckGo_lessthan_time_dot_Time_greaterthan []time.Time

func (p FuckGo_lessthan_time_dot_Time_greaterthan) Len() int           { return len(p) }
func (p FuckGo_lessthan_time_dot_Time_greaterthan) Less(i, j int) bool { return p[i].Before(p[j]) }
func (p FuckGo_lessthan_time_dot_Time_greaterthan) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func HandleRequest(w http.ResponseWriter, req *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	log.Println("got request:", *req)
	requesterName := PlayerName(strings.TrimLeft(req.URL.Path, "/"))
	board, ok := game.Boards[requesterName]
	if !ok {
		log.Println("request is for player", req.URL.Path, ", who doesn't exist")
		w.WriteHeader(404)
		return
	}
	if req.Method == "GET" {

		countdownTimes := board.falseAlarmTimes

		dead := false
		for launcherName, launcherBoard := range game.Boards {
			if launcherName == requesterName {
				continue
			}
			if launcherBoard.launchedTime == nil {
				continue
			}
			if time.Since(*launcherBoard.launchedTime) > MissileFlightTime {
				dead = true
				fmt.Fprintln(w, "you have been killed by", launcherName)
			}
			countdownTimes = append(countdownTimes, *launcherBoard.launchedTime)
		}

		if dead {
			return
		}

		sort.Sort(FuckGo_lessthan_time_dot_Time_greaterthan(countdownTimes))
		for _, t := range countdownTimes {
			fmt.Fprintln(w, time.Until(t.Add(MissileFlightTime)))
		}
	} else if req.Method == "POST" {
		if board.launchedTime != nil {
			fmt.Fprintln(w, "you have already launched")
			w.WriteHeader(400)
			return
		}
		now := time.Now()
		board.launchedTime = &now
	}
}

func main() {
	fmt.Println("hi")
	go (func() {
		for {
			line, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			victimName := PlayerName(strings.TrimSpace(line))
			addFalseAlarm(victimName)
		}
	})()
	http.HandleFunc("/", HandleRequest)
	log.Fatal(http.ListenAndServe(":2344", nil))
}

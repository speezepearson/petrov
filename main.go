package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
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
	MissileFlightTime        = 1 * time.Minute
	MeanFalseAlarmsPerSecond = 1 / float64(60)
	RollsPerSecond           = 1
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

var playerList []PlayerName

func init() {
	for p := range game.Boards {
		playerList = append(playerList, p)
	}
}

var mutex sync.Mutex

func addFalseAlarm(victimName PlayerName, at time.Time) {
	mutex.Lock()
	defer mutex.Unlock()
	board, ok := game.Boards[victimName]
	log.Println("got a false alarm for '", victimName, "'")
	if !ok {
		log.Println("...who doesn't exist")
		return
	}
	board.falseAlarmTimes = append(board.falseAlarmTimes, at)
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

		countdownTimes := append([]time.Time{}, board.falseAlarmTimes...)

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

func addFalseAlarmToRandomVictim(at time.Time) {
	i := rand.Intn(len(playerList))
	victimName := playerList[i]

	addFalseAlarm(victimName, at)
}

func main() {
	go (func() {
		for {
			line, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			victimName := PlayerName(strings.TrimSpace(line))
			addFalseAlarm(victimName, time.Now())
		}
	})()

	go func() {
		ticker := time.NewTicker(time.Second / RollsPerSecond)
		for {
			<-ticker.C

			if rand.Float64() < float64(MeanFalseAlarmsPerSecond)/RollsPerSecond {
				i := rand.Intn(len(playerList))
				victimName := playerList[i]

				addFalseAlarm(victimName)
			}
		}
	}()

	http.HandleFunc("/", HandleRequest)
	log.Fatal(http.ListenAndServe(":2344", nil))
}

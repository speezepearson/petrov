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
	GameDuration             = 1 * time.Minute
	MissileFlightTime        = 15 * time.Second
	MeanFalseAlarmsPerSecond = 1 / float64(60)
	RollsPerSecond           = 1
)

type Game struct {
	Started *time.Time // nil -> not started
	Boards  map[PlayerName]*PlayerBoard
}

func missileLanded(now, launched time.Time) bool {
	return launched.Add(MissileFlightTime).After(now)
}

func (g *Game) TimersRemainLive(now time.Time) bool {
	for _, board := range g.Boards {
		if board.launchedTime != nil && !missileLanded(now, *board.launchedTime) {
			return true
		}

		for _, alarm := range board.falseAlarmTimes {
			if !missileLanded(now, alarm) {
				return true
			}
		}
	}

	return false
}

type gamePhase int

const (
	PreStart gamePhase = iota
	Running
	Overtime
	Ended
)

func (g *Game) AnyMissileLanded(now time.Time) bool {
	for _, launcherBoard := range g.Boards {
		if launcherBoard.launchedTime == nil {
			continue
		}

		if now.Sub(*launcherBoard.launchedTime) > MissileFlightTime {
			return true
		}
	}

	return false
}

func (g *Game) Phase(now time.Time) gamePhase {
	if g.Started == nil {
		return PreStart
	}

	if now.After((*g.Started).Add(GameDuration)) {
		if g.TimersRemainLive(now) {
			return Overtime
		} else {
			return Ended
		}
	}

	if g.AnyMissileLanded(now) {
		return Ended
	}

	return Running
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
	if game.Phase(at) == Ended || game.Phase(at) == PreStart {
		log.Println("GAME IS OVER! ignored false alarm for", victimName)
		return
	}

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
	now := time.Now()
	requesterName := PlayerName(strings.TrimLeft(req.URL.Path, "/"))
	board, ok := game.Boards[requesterName]
	if !ok {
		log.Println("request is for player", req.URL.Path, ", who doesn't exist")
		w.WriteHeader(404)
		return
	}
	if req.Method == "GET" {

		if game.Started == nil {
			fmt.Fprintln(w, "game not started")
			return
		}

		countdownTimes := append([]time.Time{}, board.falseAlarmTimes...)

		dead := false
		for launcherName, launcherBoard := range game.Boards {
			if launcherName == requesterName {
				continue
			}
			if launcherBoard.launchedTime == nil {
				continue
			}
			if now.Sub(*launcherBoard.launchedTime) > MissileFlightTime {
				dead = true
				fmt.Fprintln(w, "you have been killed by", launcherName)
			}
			countdownTimes = append(countdownTimes, *launcherBoard.launchedTime)
		}

		if dead {
			return
		}

		// Player is alive, so tell them status.
		if game.Phase(now) == Ended {
			fmt.Fprintln(w, "game is over")
			return
		}

		timeLeft := (*game.Started).Add(GameDuration).Sub(now)
		if timeLeft < 0 {
			fmt.Fprintln(w, "** OVERTIME:", -timeLeft, " **")
		} else {
			fmt.Fprintln(w, timeLeft, "remaining...")
		}

		sort.Sort(FuckGo_lessthan_time_dot_Time_greaterthan(countdownTimes))
		for _, t := range countdownTimes {
			if !missileLanded(now, t) {
				fmt.Fprintf(w, "%.3f\n", (t.Add(MissileFlightTime)).Sub(now).Seconds())
			}
		}
	} else if req.Method == "POST" {
		if board.launchedTime != nil {
			fmt.Fprintln(w, "you have already launched")
			w.WriteHeader(400)
			return
		}
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
				addFalseAlarmToRandomVictim(time.Now())
			}
		}
	}()

	port := "2344"
	log.Println("listening on", port)
	http.HandleFunc("/", HandleRequest)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

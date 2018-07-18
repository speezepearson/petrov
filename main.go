package main

import (
	"encoding/json"
	"bufio"
	"flag"
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

var Debug = flag.Bool("debug", false, "")

type PlayerName string

type PlayerBoard struct {
	falseAlarmTimes []time.Time
	launchedTime    *time.Time
}

func (pb *PlayerBoard) String() string {
	var falseAlarmStrs []string
	for _, alarm := range pb.falseAlarmTimes {
		falseAlarmStrs = append(falseAlarmStrs, alarm.String())
	}

	return fmt.Sprintf("launch:%+v, alarms:{%s}",
		pb.launchedTime,
		strings.Join(falseAlarmStrs, ", "))
}

var GameDuration = flag.Duration("GameDuration", 1*time.Minute, "")

const (
	MissileFlightTime                 = 10 * time.Second
	MeanFalseAlarmsPerSecondPerPlayer = 1 / float64(30)
)

type Game struct {
	Started *time.Time // nil -> not started
	Boards  map[PlayerName]*PlayerBoard
}

func (g *Game) Start(now time.Time) {
	g.Started = &now
}

func (g *Game) String() string {
	now := time.Now()
	var boardStrings []string

	for player, board := range g.Boards {
		boardStrings = append(boardStrings, fmt.Sprintf("%s{%s}", player, board.String()))
	}

	return fmt.Sprintf("Game{Phase:%s, Started: %+v, Boards: [%s]}", g.Phase(now), g.Started, strings.Join(boardStrings, ", "))
}

func missileLanded(now, launched time.Time) bool {
	return now.After(launched.Add(MissileFlightTime))
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

type gamePhase string

const (
	PreStart gamePhase = "PreStart"
	Running            = "Running"
	Overtime           = "Overtime"
	Ended              = "Ended"
)

func (g *Game) AnyMissileLanded(now time.Time) bool {
	for _, launcherBoard := range g.Boards {
		if launcherBoard.launchedTime == nil {
			continue
		}

		if missileLanded(now, *launcherBoard.launchedTime) {
			return true
		}
	}

	return false
}

func (g *Game) PlayerIsAlive(now time.Time, player PlayerName) bool {
	for otherPlayer, launcherBoard := range g.Boards {
		if player == otherPlayer {
			continue
		}

		if launcherBoard.launchedTime == nil {
			continue
		}

		if missileLanded(now, *launcherBoard.launchedTime) {
			return false
		}
	}

	return true
}

func (g *Game) Phase(now time.Time) gamePhase {
	if g.Started == nil {
		return PreStart
	}

	if g.AnyMissileLanded(now) && !g.TimersRemainLive(now) {
		return Ended
	}

	if now.After((*g.Started).Add(*GameDuration)) {
		if g.TimersRemainLive(now) {
			return Overtime
		} else {
			return Ended
		}
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

func addFalseAlarmsForever(victimName PlayerName) {
	for {
		delay := rand.ExpFloat64() / MeanFalseAlarmsPerSecondPerPlayer
		time.Sleep(time.Duration(delay * float64(time.Second)))
		func() {
			mutex.Lock()
			defer mutex.Unlock()
			now := time.Now()
			addFalseAlarm(victimName, now)
		}()
	}
}

func addFalseAlarm(victimName PlayerName, at time.Time) {
	if game.Phase(at) == Ended || game.Phase(at) == PreStart {
		log.Println("GAME NOT RUNNING! ignored false alarm for", victimName)
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

type PlayerView struct {
	TimeRemaining time.Duration
	AlarmTimesRemaining []time.Duration
	KilledBy string
	Phase gamePhase
	TimeToMyImpact *time.Duration
}

type FuckGo_lessthan_time_dot_Time_greaterthan []time.Time

func (p FuckGo_lessthan_time_dot_Time_greaterthan) Len() int           { return len(p) }
func (p FuckGo_lessthan_time_dot_Time_greaterthan) Less(i, j int) bool { return p[i].Before(p[j]) }
func (p FuckGo_lessthan_time_dot_Time_greaterthan) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func mustSucceed(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (game *Game) View(p PlayerName, now time.Time) PlayerView {
	result := PlayerView{
		Phase: game.Phase(now),
		TimeRemaining: (*game.Started).Add(*GameDuration).Sub(now),
	}
	if game.Started == nil {
		return result
	}
	board, _ := game.Boards[p]
	if board.launchedTime != nil {
		timeToMyImpact := board.launchedTime.Add(MissileFlightTime).Sub(now)
		result.TimeToMyImpact = &timeToMyImpact
	}
	alarmTimes := append([]time.Time{}, board.falseAlarmTimes...)

	for launcherName, launcherBoard := range game.Boards {
		if launcherName == p {
			continue
		}
		if launcherBoard.launchedTime == nil {
			continue
		}
		if missileLanded(now, *launcherBoard.launchedTime) {
			result.KilledBy = string(launcherName)
			return result
		}
		alarmTimes = append(alarmTimes, *launcherBoard.launchedTime)
	}

	// Player is alive, so tell them status.
	if game.Phase(now) == Ended {
		return result
	}

	var alarmTimesRemaining []time.Duration

	sort.Sort(FuckGo_lessthan_time_dot_Time_greaterthan(alarmTimes))

	for _, t := range alarmTimes {
		if !missileLanded(now, t) {
			alarmTimesRemaining = append(alarmTimesRemaining, t.Add(MissileFlightTime).Sub(now))
		}
	}

	result.AlarmTimesRemaining = alarmTimesRemaining
	return result
}

func HandleRequest(w http.ResponseWriter, req *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	now := time.Now()
	requesterName := PlayerName(strings.TrimLeft(req.URL.Path, "/"))

	if *Debug {
		log.Println(requesterName, "req=", *req, "state=", game.String())
	}

	board, ok := game.Boards[requesterName]
	if !ok {
		log.Println("request is for player", req.URL.Path, ", who doesn't exist")
		w.WriteHeader(404)
		return
	}
	if req.Method == "GET" {
		pv := game.View(requesterName, now)
		j, err := json.Marshal(pv)
		mustSucceed(err)
		_, err = w.Write(j)
		if err != nil {
			log.Println("err:", err)
		}
		return
	} else if req.Method == "POST" {
		if !game.PlayerIsAlive(now, requesterName) {
			w.WriteHeader(400)
			fmt.Fprintln(w, "can't launch - you are dead!")
			log.Println("dead player tried to launch", requesterName)
			return
		}

		if game.Phase(now) != Running && game.Phase(now) != Overtime {
			w.WriteHeader(400)
			fmt.Fprintln(w, "can't launch - game is not running!")
			log.Println("out-of-bounds launch attempt from", requesterName)
			return
		}

		if board.launchedTime != nil {
			w.WriteHeader(400)
			fmt.Fprintln(w, "you have already launched")
			log.Println("dupe launch from", requesterName)
			return
		}
		log.Println("launch! from", requesterName)
		board.launchedTime = &now
	}
}

func main() {
	flag.Parse()

	go (func() {
		for {
			line, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			line = strings.TrimSpace(line)
			args := strings.Fields(line)

			if len(args) == 0 {
				log.Println("couldn't parse:", line)
				continue
			}

			command := args[0]

			func() {
				mutex.Lock()
				defer mutex.Unlock()

				now := time.Now()

				switch command {
				case "a":
					if len(args) != 2 {
						log.Println("wrong number of args for", command)
						return
					}
					victimName := PlayerName(args[1])
					addFalseAlarm(victimName, time.Now())

				case "s":
					if game.Phase(now) != PreStart {
						log.Println("game not in prestart")
						return
					}

					game.Start(now)
					log.Println("started game at", *game.Started)

				case "d":
					log.Println(game)
					log.Println(game.Phase(now))
				}
			}()

		}
	})()

	mutex.Lock()
	for player, _ := range game.Boards {
		go addFalseAlarmsForever(player)
	}
	mutex.Unlock()

	port := "2344"
	log.Println("listening on", port)
	http.HandleFunc("/", HandleRequest)
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./client"))))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

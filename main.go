package main

import (
	"bufio"
	secure_rand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime/debug"
	"sort"
	"strings"
	"sync"
	"time"
)

var Debug = flag.Bool("debug", false, "")
var port = flag.String("port", "2344", "")
var players = flag.String("players", "Alice,Bob",
	"Comma-delimited list of players. Optionally, specify a password with e.g. Alice:password,Bob:otherpassword")

// var hostname = flag.String("hostname", "globalthermonuclearwar.org", "Hostname to use for generating secret URLs")
var hostname = flag.String("hostname", "localhost", "Hostname to use for generating secret URLs")
var GameDuration = flag.Duration("GameDuration", 1*time.Minute, "")
var MissileFlightTime = flag.Duration("MissileFlightTime", 10*time.Second,
	"Time from launch notification to missile impact. Also controls false alarm duration.")
var MeanFalseAlarmInterval = flag.Duration(
	"MeanFalseAlarmInterval", 30*time.Second,
	"Mean interval between false alarms.")

type Password string
type PlayerName string

type PlayerBoard struct {
	falseAlarmTimes []time.Time
	launchedTime    *time.Time
	showsIfLaunched bool
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

type Game struct {
	Started *time.Time // nil -> not started
	Boards  map[PlayerName]*PlayerBoard

	PasswordToPlayer map[Password]PlayerName
	PlayerToPassword map[PlayerName]Password
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

func missileLandingTime(launched time.Time) time.Time {
	return launched.Add(*MissileFlightTime)
}

func missileLanded(now, launched time.Time) bool {
	return now.After(missileLandingTime(launched))
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
			return Running
		} else {
			return Ended
		}
	}

	return Running
}

// lol global variables
var game Game

var mutex sync.Mutex

func addFalseAlarmsForever(victimName PlayerName, meanFalseAlarmsPerSecond float64) {
	for {
		delay :=
			time.Duration((rand.ExpFloat64() / meanFalseAlarmsPerSecond) * float64(time.Second))
		log.Println(
			"scheduled false alarm for", victimName,
			"at approximately", time.Now().Add(delay).Format("2006-01-02 15:04:05"))
		time.Sleep(delay)
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
	Player              string
	TimeRemaining       time.Duration
	AlarmTimesRemaining []time.Duration
	KilledBy            string
	Phase               gamePhase
	TimeToMyImpact      *time.Duration
	MissileFlightTime   time.Duration
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

func assert(condition bool) {
	if !condition {
		log.Println("OH NOES AN INVARIANT:", string(debug.Stack()))
		if *Debug {
			panic("INVARIANT VIOLATED D:")
		}
	}
}

func (game *Game) View(p PlayerName, now time.Time) PlayerView {
	result := PlayerView{
		Player:            string(p),
		Phase:             game.Phase(now),
		MissileFlightTime: *MissileFlightTime,
	}
	if game.Started == nil {
		assert(game.Phase(now) == PreStart)
		return result
	}

	result.TimeRemaining = (*game.Started).Add(*GameDuration).Sub(now)

	board, _ := game.Boards[p]
	if (game.Phase(now) == Ended || board.showsIfLaunched) && board.launchedTime != nil {
		timeToMyImpact := board.launchedTime.Add(*MissileFlightTime).Sub(now)
		result.TimeToMyImpact = &timeToMyImpact
	}
	alarmTimes := append([]time.Time{}, board.falseAlarmTimes...)

	{
		type killedByRecord struct {
			who PlayerName
			at  time.Time
		}

		var killedBy *killedByRecord
		for launcherName, launcherBoard := range game.Boards {
			if launcherName == p {
				continue
			}
			if launcherBoard.launchedTime == nil {
				continue
			}
			if missileLanded(now, *launcherBoard.launchedTime) {
				possiblyKilledBy := killedByRecord{
					who: launcherName,
					at:  missileLandingTime(*launcherBoard.launchedTime),
				}

				if killedBy == nil {
					killedBy = &possiblyKilledBy
				}

				if possiblyKilledBy.at.Before(killedBy.at) {
					killedBy = &possiblyKilledBy
				}
			}

			alarmTimes = append(alarmTimes, *launcherBoard.launchedTime)
		}

		if killedBy != nil {
			result.KilledBy = string(killedBy.who)
			return result
		}
	}

	// Player is alive, so tell them status.
	if game.Phase(now) == Ended {
		return result
	}

	var alarmTimesRemaining []time.Duration

	sort.Sort(FuckGo_lessthan_time_dot_Time_greaterthan(alarmTimes))

	for _, t := range alarmTimes {
		if !missileLanded(now, t) {
			alarmTimesRemaining = append(alarmTimesRemaining, t.Add(*MissileFlightTime).Sub(now))
		}
	}

	result.AlarmTimesRemaining = alarmTimesRemaining
	return result
}

type Action int

const (
	Action_Status Action = iota
	Action_Launch
	Action_Conceal
)

func (g gameHandler) action(w http.ResponseWriter, req *http.Request, requesterPassword Password, action Action) {
	mutex.Lock()
	defer mutex.Unlock()

	replyErr := func(code int, msg string) {
		w.WriteHeader(code)
		fmt.Fprintln(w, msg)
		log.Println("bad request:", *req, " ==> ", code, msg)
	}

	now := time.Now()

	requesterName, ok := game.PasswordToPlayer[requesterPassword]
	if !ok {
		replyErr(403, "bad password")
		return
	}

	board, ok := game.Boards[requesterName]
	if !ok {
		replyErr(
			404,
			fmt.Sprintf("request is for player %s , who doesn't exist", requesterName))
		return
	}

	switch action {
	case Action_Status:
		pv := game.View(requesterName, now)
		j, err := json.Marshal(pv)
		mustSucceed(err)
		_, err = w.Write(j)
		if err != nil {
			log.Println("err:", err)
		}
		return

	case Action_Launch:
		board.showsIfLaunched = true

		if !game.PlayerIsAlive(now, requesterName) {
			replyErr(400, "can't launch - you are dead!")
			return
		}

		if game.Phase(now) != Running {
			replyErr(400, "can't launch - game is not running!")
			return
		}

		if board.launchedTime != nil {
			replyErr(400, "you have already launched")
			return
		}
		log.Println("launch! from", requesterName)
		board.launchedTime = &now

	case Action_Conceal:
		log.Println("conceal! from", requesterName)
		board.showsIfLaunched = false

	default:
		replyErr(500, fmt.Sprintf("unknown action %d", action))
	}
}

type gameHandler struct {
	files map[string][]byte
}

func NewGameHandler(filePaths map[string]string) (*gameHandler, error) {
	files := make(map[string][]byte)
	for servePath, path := range filePaths {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		files[servePath] = data
	}

	return &gameHandler{files: files}, nil
}

func (g gameHandler) serveFile(w http.ResponseWriter, path string) error {
	data, ok := g.files[path]
	if !ok {
		return errors.New(fmt.Sprintf("unknown path %s", path))
	}

	w.WriteHeader(200)
	if _, err := w.Write(data); err != nil {
		log.Println("failed writing response?!", err)
	}
	return nil
}

func SetAnticachingHeaders(headers http.Header) {
	// taken from https://stackoverflow.com/questions/49547/how-to-control-web-page-caching-across-all-browsers
	headers["Cache-Control:"] = []string{"no-cache", "no-store", "must-revalidate"}
	headers["Pragma:"] = []string{"no-cache"}
	headers["Expires:"] = []string{"0"}
}

func (g gameHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if *Debug {
		log.Println("req=", *req, "state=", game.String())
	}

	SetAnticachingHeaders(w.Header())

	replyErr := func(code int, msg string) {
		w.WriteHeader(code)
		fmt.Fprintln(w, msg)
		log.Println("bad request:", *req, " ==> ", code, msg)
	}

	rawComponents := strings.Split(strings.TrimLeft(req.URL.Path, "/"), "/")
	var components []string
	for _, raw := range rawComponents {
		if raw != "" {
			components = append(components, raw)
		}
	}

	switch req.Method {
	case "GET":
		switch len(components) {
		case 0:
			w.WriteHeader(200)
			w.Write([]byte("<h1>THIS SPACE INTENTIONALLY LEFT BLANK</h1>"))

		case 1:
			if err := g.serveFile(w, "/"); err != nil {
				replyErr(404, err.Error())
			}
			return

		case 2:
			if components[1] == "status" {
				g.action(w, req, Password(components[0]), Action_Status)
			} else {
				if err := g.serveFile(w, components[1]); err != nil {
					replyErr(404, err.Error())
				}
				return
			}

		default:
			replyErr(404, fmt.Sprintf("bad GET path: %s", req.URL.Path))

		}
		return

	case "POST":
		if len(components) != 2 {
			replyErr(404, fmt.Sprintf("bad POST path: %s", req.URL.Path))
			return
		}

		password, actionStr := Password(components[0]), components[1]
		switch actionStr {
		case "launch":
			g.action(w, req, password, Action_Launch)
		case "conceal":
			g.action(w, req, password, Action_Conceal)
		default:
			replyErr(404, fmt.Sprintf("unknown action: %s", actionStr))
		}
		return

	default:
		replyErr(400, fmt.Sprintf("bad request type: %s", req.Method))
		return
	}

}

func makePassword() Password {
	bytes := make([]byte, 32)
	_, err := secure_rand.Read(bytes)
	mustSucceed(err)
	return Password(hex.EncodeToString(bytes))
}

func parsePlayerSpec(specStr string) (PlayerName, Password) {
	spec := strings.Split(specStr, ":")
	switch len(spec) {
	case 1:
		return PlayerName(spec[0]), makePassword()
	case 2:
		return PlayerName(spec[0]), Password(spec[1])
	default:
		log.Fatalf("malformed player spec: %s", specStr)
		panic("wat")
	}
}

func main() {
	flag.Parse()

	logFileName := time.Now().Format("2006-01-02_15:04:05") + ".log"
	logFile, err := os.Create(logFileName)
	mustSucceed(err)
	log.SetOutput(io.MultiWriter(logFile, os.Stdout))

	playerSpecs := strings.Split(*players, ",")
	if len(playerSpecs) != 2 {
		log.Fatal("only 2 players supported ;)")
	}

	// LOL, global variables
	game.Boards = make(map[PlayerName]*PlayerBoard)
	game.PasswordToPlayer = make(map[Password]PlayerName)
	game.PlayerToPassword = make(map[PlayerName]Password)
	for _, spec := range playerSpecs {
		playerName, password := parsePlayerSpec(spec)

		game.Boards[playerName] = &PlayerBoard{}
		game.PasswordToPlayer[password] = playerName
		game.PlayerToPassword[playerName] = password
	}

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

	meanFalseAlarmsPerSecondPerPlayer := 1 / MeanFalseAlarmInterval.Seconds()

	mutex.Lock()
	for player, _ := range game.Boards {
		go addFalseAlarmsForever(player, meanFalseAlarmsPerSecondPerPlayer)
	}
	mutex.Unlock()

	gameHandler, err := NewGameHandler(map[string]string{
		"/":             "./client/dist/index.html",
		"app.bundle.js": "./client/dist/app.bundle.js",
		"map.png":       "./client/dist/src/map.png",
		"klaxon.mp3":    "./client/dist/src/klaxon.mp3",
		"launch.mp3":    "./client/dist/src/launch.mp3",
		"kaboom.mp3":    "./client/dist/src/kaboom.mp3",
	})
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", gameHandler)

	log.Println("listening on", *port)

	for playerName, password := range game.PlayerToPassword {
		url := fmt.Sprintf("http://%s:%s/%s", *hostname, *port, password)

		log.Println(playerName, "@", url)
	}

	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

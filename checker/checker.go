package main

import (
	"encoding/json"
	"flag"
	"github.com/bmizerany/pat"
	"go-concurrency/drunker/client"
	"go-concurrency/drunker/database"
	"log"
	"net/http"
	"strconv"
	"go-concurrency/messages"
	"io"
	"bytes"
)

var (
	orderChan = make(chan message.Order, 30)
	clearChan = make(chan bool, 1)
	redisHost string
	redisPort string
)

type checker struct {
}

func newChecker() *checker {
	d := new(checker)
	return d
}

func main() {
	// flag parsing
	port := flag.String("port", "3002", "rest interface listening port")
	host := flag.String("host", "192.168.1.2", "rest host")
	flag.StringVar(&redisHost, "redisHost", "127.0.0.1", "redis address ip")
	flag.StringVar(&redisPort, "redisPort", "6379", "redis port")
	flag.Parse()

	go printScore()
	initChecker(newChecker(), *host, *port)
}

// init the db connection
func getRedisConnection() client.DbClient {
	r, errR := database.NewRedis(redisHost + ":" + redisPort)
	if errR != nil {
		log.Panicf("failed to connect to redis bd")
	}
	return r
}

// init the checker
func initChecker(d *checker, host, port string) {
	m := pat.New()
	bind(m, d)
	http.Handle("/", m)
	error := http.ListenAndServe(host + ":" + port, nil)
	if error != nil {
		log.Printf("The server stop because of %v", error)
	}
}

func bind(p *pat.PatternServeMux, d *checker) {
	p.Post("/orders", http.HandlerFunc(d.onCheck))
	p.Post("/clear", http.HandlerFunc(clear))
}

func clear(w http.ResponseWriter, r *http.Request) {
	clearChan <- true
}

// what to do when receiving order request to check
func (d *checker) onCheck(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	redis := getRedisConnection()
	defer redis.Close()
	var m message.OrderCheck
	var o message.Order
	unmarshallRestBody(r, &m)
	res, e := redis.Get(strconv.Itoa(int(m.Id)))
	if e != nil {
		log.Panic(e)
	}
	umarshallMess(bytes.NewBuffer(res.([]byte)), &o)
	if o.PlayerId != m.PlayerId {
		return
	}
	redis.Remove(strconv.Itoa(int(m.Id)))
	orderChan <- o

}

func unmarshallRestBody(r *http.Request, m interface{}) {
	umarshallMess(r.Body, m)
}

func umarshallMess(from io.Reader, to interface{}) {
	err := json.NewDecoder(from).Decode(to)
	if err != nil {
		log.Panicf("Error when trying to decode request body : %s", err.Error())
	}
}

func printScore() {
	keys := make([]string, 0)
	score := make(map[string]int)
	for {
		select {
		case <- clearChan:
			keys = make([]string, 0)
			score = make(map[string]int)
		default:
		o := <-orderChan
		prevScore, ok := score[o.PlayerId]
		if !ok {
			keys = append(keys, o.PlayerId)
			score[o.PlayerId] = getScore(&o)
		} else {
			score[o.PlayerId] = prevScore + getScore(&o)
		}
		scoreStr := ""
		for _, k := range keys {
			value, _ := score[k]
			scoreStr = scoreStr + " " + k + " : " + strconv.Itoa(value)
		}
		log.Print(scoreStr)
	}
	}
}

func getScore(order *message.Order) (score int) {
	switch order.Type {
	case message.Beer:
		score = 500
	case message.Cocktail:
		score = 10000
	case message.RedWine:
		score = 1000
	case message.Vodka:
		score = 2500
	case message.Whisky:
		score = 2000
	case message.WhiteWine:
		score = 1500
	default:
		panic("unrecognized tyope")
	}
	return
}

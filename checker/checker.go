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
)

type checker struct {
	redis client.DbClient
}

func newChecker(r client.DbClient) *checker {
	d := new(checker)
	d.redis = r
	return d
}

func main() {
	// flag parsing
	port := flag.String("port", "8080", "rest interface listening port")
	host := flag.String("host", "", "rest host")
	rPort := flag.String("rPort", "6379", "redis db port")
	rhost := flag.String("rHost", "", "redis host")
	flag.Parse()

	r := initRedis(*rhost, *rPort)
	initChecker(newChecker(r), *host, *port)
}

// init the db connection
func initRedis(host, port string) client.DbClient {
	r, errR := database.NewRedis(host + ":" + port)
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
}

// what to do when receiving order request to check
func (d *checker) onCheck(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovery on some error : %s", r)
		}
	}()
	var m message.OrderCheck
	var o message.Order
	unmarshallRestBody(r, &m)
	res, e := d.redis.Get(strconv.Itoa(int(m.Id)))
	if e != nil {
		log.Panic(e)
	}
	umarshallMess(message.GetReader(res), &o)
	if o.PlayerId != m.PlayerId {
		return
	}
	d.redis.Remove(strconv.Itoa(int(m.Id)))
	//todo increase the score.

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

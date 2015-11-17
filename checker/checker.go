package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"flag"
	"github.com/bmizerany/pat"
	"github.com/jeromedoucet/go-concurrency/drunker/client"
	"github.com/jeromedoucet/go-concurrency/drunker/database"
	"github.com/jeromedoucet/go-concurrency/drunker/message"
	"log"
	"net/http"
	"strconv"
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
	error := http.ListenAndServe(host+":"+port, nil)
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
			log.Printf("Recovery on some error : %f", r)
		}
	}()
	var m message.Order
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		log.Panicf("Error when trying to decode request body : %s", err.Error())
	}
	res, e := d.redis.Get(strconv.Itoa(int(m.Id)))
	if e != nil {
		log.Panic(e)
	}

	if &res != nil {
		//ref := message.Order(res)
		ref := umarshallMess(res)

		//Check PlayerId
		if ref.PlayerId != m.PlayerId {
			log.Panicf("Error PlayerIds do not match. Expected %s Actual %s", m.PlayerId, ref.PlayerId)
		}

		//Check beverage
		if ref.Type != m.Type {
			log.Panicf("Error Beverage Types do not match. Expected %s Actual %s", m.Type, ref.Type)
		}

		//Check OK : Try deleting key
		e := d.redis.Remove(strconv.Itoa(int(m.Id)))
		if e != nil {
			log.Panic(e)
		}

		//TODO: store the new score !
	}
}

func umarshallMess(data interface{}) message.Order {
	var m message.Order
	b := &bytes.Buffer{}
	binary.Write(b, binary.BigEndian, data)
	json.Unmarshal(b.Bytes(), &m)
	return m
}

package main

import (
	"flag"
	"github.com/bitly/go-nsq"
	"github.com/bmizerany/pat"
	"go-concurrency/drunker/client"
	"go-concurrency/drunker/database"
	"log"
	"sync"
	"net/http"
	"strconv"
)

const (
	addProducerPath string = "/drunker/producers/add/:nb"
	removeProducerPath string = "/drunker/producers/remove/:nb"
	getProducerNbPath string = "/drunker/producers/nb"
)

var (
	nbProducer int
	nsqHost string
	nsqPort string
	redisHost string
	redisPort string
	host string
	port string
	wg *sync.WaitGroup
	clients    []*client.Client
)

func main() {
	flag.IntVar(&nbProducer, "nbProducer", 1, "number of producer to run")
	flag.StringVar(&nsqHost, "nsqHost", "127.0.0.1", "nsq host")
	flag.StringVar(&nsqPort, "nsqPort", "4150", "nsq port")
	flag.StringVar(&redisHost, "redisHost", "127.0.0.1", "redis host")
	flag.StringVar(&redisPort, "redisPort", "6379", "redis port")
	flag.StringVar(&host, "host", "127.0.0.1", "rest api host")
	flag.StringVar(&port, "port", "8088", "rest api port")
	flag.Parse()
	log.Printf("GO-CONCURRENCY producer module is starting with %d prducer", nbProducer)
	wg = new(sync.WaitGroup)
	clients = make([]*client.Client, nbProducer)
	for i := 0; i < nbProducer; i++ {
		startOneProducer(wg)
	}
	trimArray()
	p := pat.New()
	bind(p)
	http.Handle("/", p)
	error := http.ListenAndServe(host + ":" + port, nil)
	if error != nil {
		log.Printf("The server stop because of %v", error)
		return
	}
	wg.Wait()
}

func bind(p *pat.PatternServeMux) {
	p.Post(addProducerPath, http.HandlerFunc(addProducer))
	p.Del(removeProducerPath, http.HandlerFunc(removeProducer))
	p.Get(getProducerNbPath, http.HandlerFunc(getProducerNbRest))
}

func startOneProducer(wg *sync.WaitGroup) {
	config := nsq.NewConfig()
	w, errN := nsq.NewProducer(nsqHost + ":" + nsqPort, config)
	if errN != nil {
		log.Printf("error during nsq producer creation: %v", errN)
	} else {
		d, errR := database.NewRedis(redisHost + ":" + redisPort)
		if errR != nil {
			log.Printf("error during redis connection: %v", errR)
		} else {
			c, _ := client.StartClient(d, w, "orders#ephemeral", wg)
			clients = append(clients, c)
		}
	}
}

func addProducer(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	nb, err := strconv.Atoi(v.Get(":nb"))
	if err != nil {
		w.WriteHeader(400)
	} else {
		for i := 0; i < nb; i++ {
			startOneProducer(wg)
		}
		w.WriteHeader(200)
	}
	trimArray()
}

func removeProducer(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	nb, err := strconv.Atoi(v.Get(":nb"))
	if err != nil {
		w.WriteHeader(400)
	} else {
		nb = min(getNbProducer(), nb)
		for i := 0; i < nb; i ++ {
			log.Println("Stopping one client")
			clients[i].StopClient()
			clients[i] = nil
		}
		clients = append(clients[:nb-1], clients[nb:]...)
		w.Write([]byte("{nbProducerRemoved:" + strconv.Itoa(nb) + "}"))
		w.WriteHeader(200)
	}
}

func getProducerNbRest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{nbProducer:" + strconv.Itoa(getNbProducer()) + "}"))
	w.WriteHeader(200)
}

func trimArray() {
	for i := 0; i < len(clients); i ++ {
		if clients[i] == nil {
			clients = append(clients[:i], clients[i + 1:]...)
		}
	}
}

func getNbProducer() int {
	trimArray()
	return len(clients)
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}



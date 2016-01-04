package main

import (
	"flag"
	"github.com/bitly/go-nsq"
	"github.com/bmizerany/pat"
	"go-concurrency/drunker/client"
	"go-concurrency/drunker/database"
	"log"
	"net/http"
	"strconv"
)

const (
	addProducerPath    string = "/drunker/producers/add/:nb"
	removeProducerPath string = "/drunker/producers/remove/:nb"
	getProducerNbPath  string = "/drunker/producers/nb"
)

var (
	nbProducer int
	nsqHost    string
	nsqPort    string
	redisHost  string
	redisPort  string
	host       string
	port       string
	frequency  int
	ttl		   int
	clients    []*client.Client
)

func main() {
	flag.IntVar(&nbProducer, "nbProducer", 1, "number of producer to run")
	flag.StringVar(&nsqHost, "nsqHost", "127.0.0.1", "nsq host")
	flag.StringVar(&nsqPort, "nsqPort", "4150", "nsq port")
	flag.StringVar(&redisHost, "redisHost", "127.0.0.1", "redis host")
	flag.StringVar(&redisPort, "redisPort", "6379", "redis port")
	flag.StringVar(&host, "host", "0.0.0.0", "rest api host")
	flag.StringVar(&port, "port", "8088", "rest api port")
	flag.IntVar(&frequency, "frequency", 2000, "interval between each message in millisecond")
	flag.IntVar(&ttl, "ttl", 9000000000, "time in second to keep msg in redis")
	flag.Parse()
	log.Printf("GO-CONCURRENCY producer module is starting with %d prducer", nbProducer)
	log.Printf("nsqHost: %s", nsqHost)
	log.Printf("nsqPort: %s", nsqPort)
	log.Printf("redisHost: %s", redisHost)
	log.Printf("redisPort: %s", redisPort)
	log.Printf("host: %s", host)
	log.Printf("port: %s", port)
	log.Printf("frequency: %d msg / second", frequency)
	log.Printf("ttl: msg are kept %d second in redis", ttl)
	clients = make([]*client.Client, nbProducer)
	for i := 0; i < nbProducer; i++ {
		startOneProducer()
	}
	trimArray()
	p := pat.New()
	bind(p)
	http.Handle("/", p)
	error := http.ListenAndServe(host+":"+port, nil)
	if error != nil {
		log.Printf("The server stop because of %v", error)
		return
	}
}

func bind(p *pat.PatternServeMux) {
	p.Post(addProducerPath, http.HandlerFunc(addProducer))
	p.Del(removeProducerPath, http.HandlerFunc(removeProducer))
	p.Get(getProducerNbPath, http.HandlerFunc(getProducerNbRest))
}

func startOneProducer() {
	config := nsq.NewConfig()
	w, errN := nsq.NewProducer(nsqHost+":"+nsqPort, config)
	if errN != nil {
		log.Printf("error during nsq producer creation: %v", errN)
	} else {
		d, errR := database.NewRedis(redisHost + ":" + redisPort)
		if errR != nil {
			log.Printf("error during redis connection: %v", errR)
		} else {
			c, _ := client.StartClient(d, w, "orders", frequency, ttl)
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
			startOneProducer()
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
		nb = doRemove(nb)
		w.Write([]byte("{nbProducerRemoved:" + strconv.Itoa(nb) + "}"))
		w.WriteHeader(200)
	}
}

func doRemove(nb int) int {
	nb = min(getNbProducer(), nb)
	for i := 0; i < nb; i++ {
		log.Println("Stopping one client")
		clients[i].StopClient()
		clients[i] = nil
	}
	clients = append(clients[:nb-1], clients[nb:]...)
	return nb
}

func getProducerNbRest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{nbProducer:" + strconv.Itoa(getNbProducer()) + "}"))
	w.WriteHeader(200)
}

func trimArray() {
	for i := 0; i < len(clients); i++ {
		if clients[i] == nil {
			clients = append(clients[:i], clients[i+1:]...)
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

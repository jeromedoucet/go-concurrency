package client

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/vil-coyote-acme/go-concurrency/commons"
)

func startNewOrderMaker(url string, rAddr string, reg commons.Registration, unRegChan chan commons.Registration) {
	go func() {
		r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
		c, erredis := redis.Dial("tcp", rAddr)
		if erredis != nil {
			log.Println(erredis.Error())
			return
		}
		defer c.Close()
		decCount := 0;
		for {
			id := int(time.Now().UTC().UnixNano())
			t := r.Intn(6)
			nb := r.Intn(21)
			o := commons.Order{Id: id, Quantity: nb, Type: commons.BeverageType(t), CallBackUrl: url + "/bill/" + reg.PlayerId + "/" + strconv.Itoa(id)}
			bdOrder, _ := json.Marshal(o)
			c.Do("SET", o.Id, string(bdOrder))
			resp, err := http.Post(reg.Ip+"/orders", "application/json", bytes.NewBuffer(bdOrder))

			if err != nil || resp.StatusCode != 200 {
				decCount ++
				if decCount >= 3 {
					unRegChan <- reg
					return
				}
				time.Sleep(5 * time.Second)
			} else {
				decCount = 0
			}
		}
	}()
}

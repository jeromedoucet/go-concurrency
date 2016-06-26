package client

import (
	"github.com/vil-coyote-acme/go-concurrency/commons"
	"time"
	"math/rand"
	"strconv"
	"github.com/garyburd/redigo/redis"
	"encoding/json"
	"net/http"
	"bytes"
)

func startNewOrderMaker(url string, rAddr string, reg commons.Registration, unRegChan chan commons.Registration) {
	go func() {
		r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
		c, _ := redis.Dial("tcp", rAddr)
		defer c.Close()
		for {
			id := int(time.Now().UTC().UnixNano())
			t := r.Intn(6)
			nb := r.Intn(21)
			o := commons.Order{Id:id, Quantity:nb, Type:commons.BeverageType(t), CallBackUrl:url + "/bill/" + reg.PlayerId + "/" + strconv.Itoa(id)}
			bdOrder, _ := json.Marshal(o)
			c.Do("SET", o.Id, string(bdOrder))
			resp, err := http.Post(reg.Ip + "/orders", "application/json", bytes.NewBuffer(bdOrder))
			if err != nil || resp.StatusCode != 200 {
				unRegChan <- reg
				return
			}
		}
	}()
}

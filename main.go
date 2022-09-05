package main

import (
	"github.com/gomodule/redigo/redis"
	"log"
	"net/http"
)

var serverAddress = ":8888"

func main() {
	gRedisConn, err := gRedisConn()
	if err != nil {
		panic(err)
	}
	defer gRedisConn.Close()

	gPubSubConn = &redis.PubSubConn{Conn: gRedisConn}
	defer gPubSubConn.Close()

	go deliverMessages()

	http.HandleFunc("/ws", wsHandler)
	log.Printf("server started at %s\n", serverAddress)

	log.Fatal(http.ListenAndServe(serverAddress, nil))
}

package main

import (
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"log"
)

var (
	gStore      *Store
	gPubSubConn *redis.PubSubConn
	gRedisConn  = func() (redis.Conn, error) {
		return redis.Dial("tcp", ":6379")
	}
)

func init() {
	gStore = &Store{
		Users: make([]*User, 0, 1),
	}
}

func (s *Store) newUser(conn *websocket.Conn) *User {
	u := &User{
		ID:   uuid.NewV4().String(),
		conn: conn,
	}

	if err := gPubSubConn.Subscribe(u.ID); err != nil {
		panic(err)
	}

	s.Lock()
	defer s.Unlock()

	s.Users = append(s.Users, u)
	return u
}

func deliverMessages() {
	for {
		switch v := gPubSubConn.Receive().(type) {
		case redis.Message:
			gStore.findAndDeliver(v.Channel, string(v.Data))

		case redis.Subscription:
			log.Printf("subscription message: %s: %s %d\n", v.Channel, v.Kind, v.Count)

		case error:
			log.Println("error pub/sub, delivery has stopped")
			return
		}
	}
}

func (s *Store) findAndDeliver(userID string, content string) {
	m := Message{
		Content: content,
	}

	for _, u := range s.Users {
		if u.ID == userID {
			if err := u.conn.WriteJSON(m); err != nil {
				log.Printf("error on message delivery e: %s\n", err)
			} else {
				log.Printf("user %s found, message sent \n", userID)
			}
			return
		}
	}

	log.Printf("user %s not found at our store\n", userID)
}

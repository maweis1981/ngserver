package main

import (
	"github.com/gorilla/websocket"
	"sync"
)

type User struct {
	ID   string
	conn *websocket.Conn
}

type Store struct {
	Users []*User
	sync.Mutex
}

type Message struct {
	DeliveryID string `json:"id"`
	Content    string `json:"content"`
}

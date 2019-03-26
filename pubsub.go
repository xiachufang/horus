package main

import (
	"container/list"
	"log"
	"sync"
)

var pubsubLock = &sync.Mutex{}
var subscribers = list.New()

func subscribe() chan interface{} {
	pubsubLock.Lock()
	defer pubsubLock.Unlock()

	ch := make(chan interface{}, 0)
	subscribers.PushBack(ch)

	return ch
}

func unsubscribe(ch chan interface{}) {
	pubsubLock.Lock()
	defer pubsubLock.Unlock()

	for e := subscribers.Front(); e != nil; e = e.Next() {
		if e.Value == ch {
			subscribers.Remove(e)
		}
	}
	close(ch)
}

func publish(data interface{}) bool {
	pubsubLock.Lock()
	defer pubsubLock.Unlock()

	published := 0
	for e := subscribers.Front(); e != nil; e = e.Next() {
		ch := e.Value.(chan interface{})
		ch <- data
		published++
	}

	log.Printf("Published to %d subscribers", published)
	return published > 0
}

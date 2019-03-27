package main

import (
	"container/list"

	log "github.com/sirupsen/logrus"
)

var subscribers = list.New()
var subscribeCh = make(chan chan interface{})
var unsubscribeCh = make(chan chan interface{})
var publishCh = make(chan interface{})
var donePubSub = make(chan struct{})

func removeFromList(val interface{}, l *list.List) {
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value == val {
			l.Remove(e)
		}
	}
}

func runPubSub() {
	go func() {
		for {
			select {
			case ch := <-subscribeCh:
				log.Debugf("subscribe")
				subscribers.PushBack(ch)
			case ch := <-unsubscribeCh:
				log.Debugf("unsubscribe\n")
				removeFromList(ch, subscribers)
				close(ch)
			case data := <-publishCh:
				for e := subscribers.Front(); e != nil; e = e.Next() {
					subCh := e.Value.(chan interface{})
				FORLOOP:
					for {
						log.Debugf("pub in for\n")
						select {
						case subCh <- data:
							log.Debugf("send success\n")
							break FORLOOP
						case ch := <-unsubscribeCh:
							log.Debugf("unsub in publish\n")
							if ch == subCh {
								log.Debugf("remove current publish\n")
								subscribers.Remove(e)
								close(ch)
								break FORLOOP
							} else {
								log.Debugf("remove other sub in publish\n")
								removeFromList(ch, subscribers)
								close(ch)
							}
						case ch := <-subscribeCh:
							log.Debugf("sub in publish\n")
							subscribers.PushBack(ch)
							break FORLOOP
						}
					}
				}
			case <-donePubSub:
				return
			}
		}
	}()
}

func stopPubSub() {
	donePubSub <- struct{}{}
}

func subscribe() chan interface{} {
	log.Debugf("Begin subscribe\n")
	ch := make(chan interface{}, 0)
	subscribeCh <- ch
	log.Debugf("Finish subscribe\n")
	return ch
}

func unsubscribe(ch chan interface{}) {
	log.Debugf("Begin unsubscribe\n")
	unsubscribeCh <- ch
	log.Debugf("Finish unsubscribe\n")
}

func publish(data interface{}) {
	log.Debugf("Begin publish\n")
	publishCh <- data
	log.Debugf("Finish publish\n")
}

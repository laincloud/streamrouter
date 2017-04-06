package watcher

import "time"

const retryTime = time.Second * 3

type Watcher interface {
	watch(notify chan interface{})
}

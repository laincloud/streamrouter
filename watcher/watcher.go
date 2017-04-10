package watcher

import (
	"fmt"
	"time"

	"github.com/laincloud/lainlet/client"
	"github.com/laincloud/streamrouter/utils"
)

const retryTime = time.Second * 3

var LainletClient = client.New(fmt.Sprintf("lainlet.lain:%s", utils.GetEnvWithDefault("LAINLET_PORT", "9001")))

type Watcher interface {
	Watch(notify chan interface{})
}

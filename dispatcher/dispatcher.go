package dispatcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/laincloud/streamrouter/backend"
	"github.com/laincloud/streamrouter/model"
	"github.com/laincloud/streamrouter/utils"
	"github.com/laincloud/streamrouter/watcher"
	"github.com/mijia/sweb/log"
)

const checkInterval = time.Minute

var (
	graphiteEndpoint = net.JoinHostPort("graphite.lain", utils.GetEnvWithDefault("GRAPHITE_PORT", "2003"))
	lainDomain       = utils.GetEnvWithDefault("LAIN_DOMAIN", "lain.local")
	instanceNo       = utils.GetEnvWithDefault("INSTANCE_NO", "1")
)

type Dispatcher struct {
	notify  chan interface{}
	backend backend.Backend
}

func Run() {
	dispatcher := Dispatcher{
		notify:  make(chan interface{}),
		backend: backend.HaproxyBackend{},
	}
	reportTicker := time.Tick(checkInterval)
	watcherList := []watcher.Watcher{
		watcher.StreamWatcher{},
	}
	for _, w := range watcherList {
		go w.Watch(dispatcher.notify)
	}
	for {
		select {
		case <-reportTicker:
			dispatcher.report()
		case data := <-dispatcher.notify:
			dispatcher.reload(data)
		}
	}
}

func (disp *Dispatcher) report() {

	valid := 1
	if disp.backend.Check() != nil {
		valid = 0
	}

	graphiteConf := make(map[string]string)
	if data, err := watcher.LainletClient.Get("/v2/configwatcher?target=features/graphite", 2*time.Second); err != nil {
		log.Errorf("Get graphite feature failed: %s", err.Error())
		return
	} else if err := json.Unmarshal(data, &graphiteConf); err != nil {
		log.Errorf("Unmarshal graphite feature failed: %s", err.Error())
		return
	} else if needReport, _ := strconv.ParseBool(graphiteConf["features/graphite"]); !needReport {
		return
	}

	conn, err := net.DialTimeout("tcp", graphiteEndpoint, time.Second*2)
	if err != nil {
		log.Errorf("Dial failed: %s", err.Error())
		return
	}
	defer conn.Close()
	timeStamp := time.Now().Unix()
	sendData := fmt.Sprintf("%s.streamrouter.syntax_valid.%s %s %d\n", lainDomain, instanceNo, valid, timeStamp)
	if _, err = conn.Write([]byte(sendData)); err != nil {
		log.Errorf("Send report data failed: %s", err.Error())
	}
}

func (disp *Dispatcher) reload(data interface{}) {
	var err error
	switch data.(type) {
	case model.StreamAppList:
		err = disp.backend.RenderStreamFiles(data.([]model.StreamApp))
	default:
		err = errors.New("Not supported datatype")
	}
	if err == nil {
		disp.backend.Reload()
	}
}

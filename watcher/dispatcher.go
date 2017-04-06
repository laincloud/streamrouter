package watcher

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"net"

	"github.com/laincloud/lainlet/client"
	"github.com/mijia/sweb/log"
)

const (
	checkInterval  = time.Minute
	execTimeout    = time.Second * 3
	haproxyConfDir = "/etc/haproxy"
	haproxyConfExt = ".cfg"
	commandName    = "haproxy"
)

var (
	lainletClient     = client.New(fmt.Sprintf("lainlet.lain:%s", getEnvWithDefault("LAINLET_PORT", "9001")))
	graphiteEndpoint  = net.JoinHostPort("graphite.lain", getEnvWithDefault("GRAPHITE_PORT", "2003"))
	lainDomain        = getEnvWithDefault("LAIN_DOMAIN", "lain.local")
	currentInstanceNo = getEnvWithDefault("INSTANCE_NO", 1)
	pidfile           = getEnvWithDefault("HAPROXY_PID_FILE", "/var/run/haproxy.pid")
	cmdParams         = []string{
		"-f",
		filepath.Join(haproxyConfDir, "haproxy"+haproxyConfExt),
		"-f",
		streamConfDir,
	}
	runParams   = append(cmdParams, "-D")
	checkParams = append(cmdParams, "-c")
)

type Dispatcher struct {
	notify chan interface{}
}

func Run() error {
	dispatcher := Dispatcher{
		notify: make(chan interface{}),
	}
	reportTicker := time.Tick(checkInterval)
	watcherList := []Watcher{
		StreamWatcher{},
	}
	for _, watcher := range watcherList {
		go watcher.watch(dispatcher.notify)
	}
	for {
		select {
		case <-reportTicker:
			dispatcher.report()
		case <-dispatcher.notify:
			dispatcher.reload()
		}
	}
	return nil
}

func (disp *Dispatcher) check() error {
	ctx, cancel := context.WithTimeout(context.Background(), execTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, commandName, checkParams...)
	return cmd.Run()
}

func (disp *Dispatcher) reload() {
	log.Info("Haproxy reload start")
	if err := disp.check(); err != nil {
		log.Errorf("Haproxy syntax checking failed", err.Error())
		return
	}
	disp.runHaproxy()
}

func (disp *Dispatcher) runHaproxy() {
	ctx, cancel := context.WithTimeout(context.Background(), execTimeout)
	defer cancel()
	var realRunParams []string
	if haproxyPID := getPidFromPidfile(pidfile); haproxyPID != 0 {
		realRunParams = append(runParams, "-sf", strconv.Itoa(haproxyPID))
	}
	cmd := exec.CommandContext(ctx, commandName, realRunParams...)
	if err := cmd.Run(); err != nil {
		log.Error(err.Error())
	} else {
		log.Info("Haproxy reload successfully")
	}
}

func (disp *Dispatcher) report() {
	haproxyPID := getPidFromPidfile(pidfile)
	if haproxyPID != 0 {
		if err := checkProcessAlive(haproxyPID); err != nil {
			log.Errorf("Haproxy process check error: %s. Now attempt to start...", err.Error())
			disp.runHaproxy()
		}
	}
	valid := 1
	if err := disp.check(); err != nil {
		log.Errorf("Check config files failed: %s", err.Error())
		valid = 0
	}

	graphiteConf := make(map[string]string)
	if data, err := lainletClient.Get("/v2/configwatcher?target=features/graphite", 2*time.Second); err != nil {
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
	sendData := fmt.Sprintf("%s.streamrouter.syntax_valid.%s %d %d\n", lainDomain, currentInstanceNo, valid, timeStamp)
	if _, err = conn.Write([]byte(sendData)); err != nil {
		log.Errorf("Send report data failed: %s", err.Error())
	}

}

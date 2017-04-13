package watcher

import (
	"context"
	"encoding/json"
	"time"

	"github.com/laincloud/streamrouter/model"
	"github.com/mijia/sweb/log"
)

type StreamWatcher struct {
}

func (s *StreamWatcher) Watch(notify chan interface{}) {
	for {
		if ch, err := LainletClient.Watch("/v2/streamrouter/streamprocs", context.Background()); err != nil {
			log.Errorf("StreamWatcher connect to lainlet failed. Retry in 3 seconds")
		} else {
			for resp := range ch {
				newData := make(model.StreamAppList)
				if resp.Event == "init" || resp.Event == "update" || resp.Event == "delete" {
					if err := json.Unmarshal(resp.Data, &newData); err != nil {
						log.Errorf("StreamWatcher unmarshall data failed: %s", err.Error())
					} else {
						appList := make([]model.StreamApp, 0, len(newData))
						for name, app := range newData {
							appList = append(appList, model.StreamApp{
								Name:        name,
								StreamProcs: app,
							})
						}
						notify <- appList
					}
				}
				time.Sleep(retryTime)
			}
			log.Infof("StreamWatcher channel closed")
		}
		time.Sleep(retryTime)
	}
}

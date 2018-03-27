package watcher

import (
	"context"
	"encoding/json"
	"reflect"
	"sort"
	"time"

	"github.com/laincloud/streamrouter/model"
	"github.com/mijia/sweb/log"
)

type StreamWatcher struct {
	oldAppList []model.StreamApp
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
						sort.Slice(appList, func(i, j int) bool {
							return appList[i].Name < appList[j].Name
						})

						if reflect.DeepEqual(appList, s.oldAppList) {
							log.Warnf("StreamWatcher receives old data: %+v.", appList)
						} else {
							s.oldAppList = appList
							notify <- appList
						}
					}
				}
				time.Sleep(retryTime)
			}
			log.Infof("StreamWatcher channel closed")
		}
		time.Sleep(retryTime)
	}
}

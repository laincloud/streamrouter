package watcher

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/mijia/sweb/log"
)

const streamConfTempl = `
{{- $appName := .Name }}
{{- range $proc := .StreamProcs }}
{{- range $service := $proc.Services }}
listen {{ $appName }}_{{ $proc.Name }}_{{ $service.ListenPort }}
  mode tcp
  bind :{{- $service.ListenPort}}
  {{- if $service.EnableHealthCheck }}
  option tcp-check
  tcp-check send {{ $service.Send }}
  tcp-check expect string {{ $service.Expect }}
  {{- end }}
{{- range $upstream := $proc.Upstreams }}
  server {{ $proc.Name }}_{{ $service.ListenPort }}_{{ $upstream.InstanceNo }} {{ $upstream.Host }}:{{ $service.UpstreamPort }} check
{{- end }}
{{ end -}}
{{ end -}}
`

var (
	streamConfTemplate = template.Must(template.New("stream").Parse(streamConfTempl))
	streamConfDir      = filepath.Join(haproxyConfDir, "stream.d")
)

type StreamWatcher struct {
}

type StreamApp struct {
	Name        string
	StreamProcs []StreamProc
}

type StreamProc struct {
	Name      string
	Upstreams []StreamUpstream
	Services  []StreamService
}

type StreamUpstream struct {
	Host       string
	InstanceNo int
}

type StreamService struct {
	UpstreamPort      int
	ListenPort        int
	Send              string
	Expect            string
	EnableHealthCheck bool
}

func (s StreamWatcher) watch(notify chan interface{}) {
	var newData []StreamApp
	for {
		if ch, err := lainletClient.Watch("/v2/streamrouter/streamwatcher", context.Background()); err != nil {
			log.Errorf("StreamWatcher connect to lainlet failed. Retry in 3 seconds")
		} else {
			for resp := range ch {
				if resp.Event == "init" || resp.Event == "update" || resp.Event == "delete" {
					if err := json.Unmarshal(resp.Data, &newData); err != nil {
						log.Errorf("StreamWatcher unmarshall data failed: %s", err.Error())
					} else {
						s.renderFiles(notify, newData)
					}
				}
				time.Sleep(retryTime)
			}
		}
		time.Sleep(retryTime)
	}
}

func (s *StreamWatcher) renderFiles(notify chan interface{}, data []StreamApp) {
	if len(data) > 0 {
		if err := removeContents(streamConfDir); err != nil {
			log.Errorf("Remove stream conf dir contents failed: %s", err.Error())
			return
		}
		for _, app := range data {
			if confData, err := s.parseStreamApp(app); err != nil {
				log.Errorf("Rendering template of app[%s] failed: %s", app.Name, err.Error())
			} else if err := ioutil.WriteFile(filepath.Join(streamConfDir, app.Name+haproxyConfExt), confData, 0644); err != nil {
				log.Errorf("Writing conf file of app[%s] failed: %s", err.Error())
			} else {
				log.Infof("Writing conf file of app[%s] successfully", err.Error())
			}
		}
		notify <- 1
	}
}

func (s *StreamWatcher) parseStreamApp(data StreamApp) ([]byte, error) {
	var byteArr []byte
	buf := bytes.NewBuffer(byteArr)

	if err := streamConfTemplate.Execute(buf, data); err != nil {
		return byteArr, err
	}
	return buf.Bytes(), nil
}

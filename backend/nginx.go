package backend

import (
	"context"
	"io/ioutil"
	"os/exec"
	"text/template"

	"time"

	"path/filepath"

	"bytes"

	"github.com/laincloud/streamrouter/model"
	"github.com/laincloud/streamrouter/utils"
	"github.com/mijia/sweb/log"
)

type NginxBackend struct {
}

const (
	nginxCommand       = "nginx"
	nginxExecTimeout   = time.Second * 3
	nginxStreamConfDir = "/etc/nginx/stream.d"
	nginxConfExt       = ".conf"
)

const nginxStreamConfTeml = `
{{- $appName := .Name }}
{{- range $proc := .StreamProcs }}
{{- range $service := $proc.Services }}
server {
    listen {{ $service.ListenPort }};
    proxy_pass {{ $appName }}_{{ $proc.Name }}_{{ $service.ListenPort}}_{{ $service.UpstreamPort }};
    access_log /var/log/nginx/{{ $appName }}.log stream;
}

upstream {{ $appName }}_{{ $proc.Name }}_{{ $service.ListenPort}}_{{ $service.UpstreamPort }} {
{{- range $upstream := $proc.Upstreams }}
    server {{ $upstream.Host }}:{{ $service.UpstreamPort }};
{{- end }}
    check interval=3000 rise=2 fall=5 timeout=1000 type=http;
}
{{ end -}}
{{ end -}}
`

var (
	nginxCheckParam      = []string{"-t"}
	nginxReloadParam     = []string{"-s", "reload"}
	ngStreamConfTemplate = template.Must(template.New("hp_stream").Parse(nginxStreamConfTeml))
)

func (nb *NginxBackend) Check() error {
	ctx, cancel := context.WithTimeout(context.Background(), nginxExecTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, nginxCommand, nginxCheckParam...)
	return cmd.Run()
}

func (nb *NginxBackend) Reload() {
	log.Info("Nginx reload start")
	ctx, cancel := context.WithTimeout(context.Background(), nginxExecTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, nginxCommand, nginxReloadParam...)
	if err := cmd.Run(); err != nil {
		log.Errorf("Nginx reload failed: %s", err.Error())
	}
	log.Info("Nginx reload successfully")
}

func (nb *NginxBackend) RenderStreamFiles(data []model.StreamApp) error {
	if len(data) > 0 {
		if err := utils.RemoveContents(nginxStreamConfDir); err != nil {
			log.Errorf("Remove stream conf dir contents failed: %s", err.Error())
			return err
		}
		for _, app := range data {
			if confData, err := nb.parseStreamApp(app); err != nil {
				log.Errorf("Rendering template of app[%s] failed: %s", app.Name, err.Error())
			} else if err := ioutil.WriteFile(filepath.Join(nginxStreamConfDir, app.Name+nginxConfExt), confData, 0644); err != nil {
				log.Errorf("Writing conf file of app[%s] failed: %s", app.Name, err.Error())
			} else {
				log.Infof("Writing conf file of app[%s] successfully", app.Name)
			}
		}
	}
	return nil
}

func (nb *NginxBackend) parseStreamApp(data model.StreamApp) ([]byte, error) {
	var byteArr []byte
	buf := bytes.NewBuffer(byteArr)

	if err := ngStreamConfTemplate.Execute(buf, data); err != nil {
		return byteArr, err
	}
	return buf.Bytes(), nil
}

package backend

import (
	"context"
	"html/template"
	"io/ioutil"
	"os/exec"
	"time"

	"path/filepath"

	"bytes"

	"github.com/laincloud/streamrouter/model"
	"github.com/laincloud/streamrouter/utils"
	"github.com/mijia/sweb/log"
)

const (
	hpExecTimeout = time.Second * 3
	hpConfDir     = "/etc/haproxy"
	hpConfExt     = ".cfg"
	hpCmdName     = "haproxy"
)

const hpStreamConfTempl = `
{{- $appName := .Name }}
{{- range $proc := .StreamProcs }}
{{- range $service := $proc.Services }}
listen {{ $appName }}_{{ $proc.Name }}_{{ $service.ListenPort }}
  mode tcp
  bind :{{- $service.ListenPort}}
  {{- if and $service.Send $service.Expect }}
  option tcp-check
  tcp-check send-binary {{ $service.Send }}
  tcp-check expect binary {{ $service.Expect }}
  {{- end }}
{{- range $upstream := $proc.Upstreams }}
  server {{ $proc.Name }}_{{ $service.ListenPort }}_{{ $upstream.InstanceNo }} {{ $upstream.Host }}:{{ $service.UpstreamPort }} check
{{- end }}
{{ end -}}
{{ end -}}
`

var (
	hpStreamConfDir      = filepath.Join(hpConfDir, "stream.d")
	hpStreamConfTemplate = template.Must(template.New("hp_stream").Parse(hpStreamConfTempl))
	hpPIDfile            = utils.GetEnvWithDefault("HAPROXY_PID_FILE", "/var/run/haproxy.pid")
	hpCmdParams          = []string{
		"-f",
		filepath.Join(hpConfDir, "haproxy"+hpConfExt),
		"-f",
		hpStreamConfDir,
	}
	hpRunParams   = append(hpCmdParams, "-D")
	hpCheckParams = append(hpCmdParams, "-c")
)

// Implements BackendInterface
type HaproxyBackend struct {
}

func (hb *HaproxyBackend) Check() error {
	ctx, cancel := context.WithTimeout(context.Background(), hpExecTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, hpCmdName, hpCheckParams...)
	return cmd.Run()
}

func (hb *HaproxyBackend) Reload() {
	log.Info("Haproxy reload start")
	if err := hb.Check(); err != nil {
		log.Errorf("Haproxy syntax checking failed", err.Error())
		return
	}
	hb.runHaproxy()
}

func (hb *HaproxyBackend) RenderStreamFiles(data []model.StreamApp) error {
	if len(data) > 0 {
		if err := utils.RemoveContents(hpStreamConfDir); err != nil {
			log.Errorf("Remove stream conf dir contents failed: %s", err.Error())
			return err
		}
		for _, app := range data {
			if confData, err := hb.parseStreamApp(app); err != nil {
				log.Errorf("Rendering template of app[%s] failed: %s", app.Name, err.Error())
			} else if err := ioutil.WriteFile(filepath.Join(hpStreamConfDir, app.Name+hpConfExt), confData, 0644); err != nil {
				log.Errorf("Writing conf file of app[%s] failed: %s", app.Name, err.Error())
			} else {
				log.Infof("Writing conf file of app[%s] successfully", app.Name)
			}
		}
	}
	return nil
}

func (hb *HaproxyBackend) runHaproxy() {
	ctx, cancel := context.WithTimeout(context.Background(), hpExecTimeout)
	defer cancel()
	realRunParams := hpRunParams
	if data, err := ioutil.ReadFile(hpPIDfile); err == nil {
		log.Infof("Haproxy is already running(PID file exists). Soft stop instead.")
		realRunParams = append(hpRunParams, "-sf", string(data))
	}
	cmd := exec.CommandContext(ctx, hpCmdName, realRunParams...)
	if err := cmd.Run(); err != nil {
		log.Error(err.Error())
	} else {
		log.Info("Haproxy reload successfully")
	}
}

func (hb *HaproxyBackend) parseStreamApp(data model.StreamApp) ([]byte, error) {
	var byteArr []byte
	buf := bytes.NewBuffer(byteArr)

	if err := hpStreamConfTemplate.Execute(buf, data); err != nil {
		return byteArr, err
	}
	return buf.Bytes(), nil
}
